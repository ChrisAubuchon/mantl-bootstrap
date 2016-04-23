package bootstrap

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/asteris-llc/mantl-bootstrap/common"
	"github.com/asteris-llc/mantl-bootstrap/consul"
	pb "github.com/asteris-llc/mantl-bootstrap/proto"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
)

const (
	tag_isBootstrapped = "is-bootstrapped"
	tag_caFileBytes = "ca_file_bytes"
	tag_caKeyBytes = "ca_key_bytes"
)

type Bootstrap struct {
	Config *viper.Viper
	Consul *common.ConsulStatic

	caFile []byte
	caKey []byte

	client pb.BootstrapRPCClient
}

func Init(root *cobra.Command) {
	b := Bootstrap{
		Config: viper.New(),
		Consul: common.NewConsulConfig(),
	}

	bCmd := &cobra.Command{
		Use:   "bootstrap",
		Short: "Bootstrap Mantl nodes",
		Long:  "Bootstrap Mantl nodes",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if !b.Config.IsSet("servers") {
				return fmt.Errorf("Must supply list of Consul server IPs")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return b.Bootstrap(args)
		},
	}

	bCmd.Flags().String("servers", "", "Comma separated list of Consul server IPs")
	bCmd.Flags().String("clients", "", "Comma separated list of Consul client IPs")
	bCmd.Flags().String("domain", "consul", "Consul DNS domain")
	bCmd.Flags().String("dc", "dc1", "Consul datacenter")

	bCmd.Flags().String("cert-country", "US", "Certificate country")
	bCmd.Flags().String("cert-state", "New York", "Certificate state/province")
	bCmd.Flags().String("cert-locality", "Anytown", "Certificate locality/city")
	bCmd.Flags().String("cert-organization", "Example Company Inc", "Certificate organization")
	bCmd.Flags().String("cert-unit", "Operations", "Certificate organizational unit inside of organization")
	bCmd.Flags().String("cert-common", "mantl", "Certificate common name")

	b.Config.BindPFlags(bCmd.Flags())

	root.AddCommand(bCmd)
}

func (b *Bootstrap) Bootstrap(args []string) error {
	if err := b.GenConsulConfig(); err != nil {
		return err
	}

	// Save my ip address for the end
	servers := strings.Split(b.Config.GetString("servers"), ",")
	clients := strings.Split(b.Config.GetString("clients"), ",")
	myip := servers[0]

	if b.Config.GetBool(tag_isBootstrapped) {
		fmt.Println("Getting consul Ips")
		consulIps, err := consul.GetIps()
		if err != nil {
			return err
		}

		b.Config.Set("existing-ips", consulIps)
	}

	fmt.Printf("Bootstrapping hosts: %s\n", b.Config.GetString("servers"))

	b.Consul.RetryJoin = servers

	// Server nodes first
	for i, ip := range servers {
		if i == 0 {
			continue
		}

		if err := b.BootstrapNode(ip, true); err != nil {
			return err
		}
	}

	fmt.Println("Bootstrapping clients")
	if len(clients) > 0 {
		for _, ip := range clients {
			if err := b.BootstrapNode(ip, false); err != nil {
				return err
			}
		}
	}

	fmt.Println("Bootstrapping self")
	return b.BootstrapNode(myip, true)
}

func (b *Bootstrap) BootstrapNode(ip string, isServer bool) error {
	if ip == "" {
		return nil
	}

	fmt.Printf("Bootstrapping ip %s\n", ip)

	b.Consul.AdvertiseAddr = ip
	b.Consul.Server = isServer
	b.Consul.Domain = b.Config.GetString("domain")
	b.Consul.Datacenter = b.Config.GetString("dc")
	b.Consul.AclDatacenter = b.Config.GetString("dc")

	if !b.Config.GetBool(tag_isBootstrapped) {
		if err := b.Configure(ip); err != nil {
			return err
		}
	} else {
		found := false
		for _, cip := range b.Config.GetStringSlice("existing-ips") {
			if cip == ip {
				fmt.Printf("ip: %s found. Skipping\n", ip)
				found = true
				break
			}
		}
		if !found {
			if err := b.Configure(ip); err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *Bootstrap) Configure(ip string) error {
	nc := common.DefaultConfig()

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", ip, nc.Port), grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	b.client = pb.NewBootstrapRPCClient(conn)

	consulJson, err := json.Marshal(b.Consul)
	if err != nil {
		return err
	}

	// Write certificate and key
	fmt.Printf("Sending '%s' file\n", b.Consul.CaFile)
	if err := b.writeFile(b.caFile, b.Consul.CaFile, 0644); err != nil {
		return err
	}

	fmt.Printf("Sending '%s' file\n", common.CaCertAnchorPath)
	if err := b.writeFile(b.caFile, common.CaCertAnchorPath, 0600); err != nil {
		return err
	}

	fmt.Printf("Sending '%s' file\n", common.CaKeyPath)
	if err := b.writeFile(b.caKey, common.CaKeyPath, 0644); err != nil {
		return err
	}

	// Write Consul Configuration
	if err := b.configureConsul(consulJson); err != nil {
		return err
	}

	if err := b.shutdown(true); err != nil {
		return err
	}

	return nil
}

// Convenience functions for RPC commands
//

func (b *Bootstrap) writeFile(data []byte, path string, mode uint32) error {
	r, err := b.client.WriteFile(context.Background(),
		&pb.FileData{
			Data: data,
			Path: path,
			Mode: mode,
		})
	if err != nil {
		return err
	}

	if r.Code != pb.Response_Success {
		return fmt.Errorf("%s", r.Mesg)	
	}

	return nil
}

func (b *Bootstrap) configureConsul(cdata []byte) error {
	r, err := b.client.ConfigureConsul(context.Background(), &pb.ConsulConfig{Data:cdata})
	if err != nil {
		return err
	}

	if r.Code != pb.Response_Success {
		return fmt.Errorf("%s", r.Mesg)	
	}

	return nil
}

func (b *Bootstrap) shutdown(success bool) error {
	msg := pb.ShutdownMsg{ Code: pb.ShutdownMsg_Failure }
	if success {
		msg.Code = pb.ShutdownMsg_Success
	}
	r, err := b.client.Shutdown(context.Background(), &msg)
	if err != nil {
		return err
	}

	if r.Code != pb.Response_Success {
		return fmt.Errorf("%s", r.Mesg)	
	}

	return nil
}
