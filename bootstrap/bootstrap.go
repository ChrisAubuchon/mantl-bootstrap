package bootstrap

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/asteris-llc/mantl-bootstrap/common"
	"github.com/asteris-llc/mantl-bootstrap/consul"
	"github.com/asteris-llc/mantl-bootstrap/structs"

	"github.com/spf13/cobra"
)

type Config struct {
	Servers string
	serversSplit []string
	Clients string
	clientsSplit []string
	Domain string
	Cert *CertData

	isBootstrapped bool
	secInfo *SecInfo
	consulIps []string
}

func Init(root *cobra.Command) {
	c := Config{
		Cert: &CertData{},
	}

	bCmd := &cobra.Command{
		Use: "bootstrap",
		Short: "Bootstrap Mantl nodes",
		Long: "Bootstrap Mantl nodes",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if c.Servers == "" {
				return fmt.Errorf("Must supply list of Consul server IPs")
			}
			c.serversSplit = strings.Split(c.Servers, ",")
			c.clientsSplit = strings.Split(c.Clients, ",")

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Bootstrap(args)
		},
	}

	bCmd.Flags().StringVar(&c.Servers, "servers", "", "Comma separated list of Consul server IPs")
	bCmd.Flags().StringVar(&c.Clients, "clients", "", "Comma separated list of Consul client IPs")
        bCmd.Flags().StringVar(&c.Domain, "domain", "consul", "Consul DNS domain")

        bCmd.Flags().StringVar(&c.Cert.Country, "cert-country", "US", "Certificate country")
        bCmd.Flags().StringVar(&c.Cert.State, "cert-state", "New York", "Certificate state/province")
        bCmd.Flags().StringVar(&c.Cert.Locality, "cert-locality", "Anytown", "Certificate locality/city")
        bCmd.Flags().StringVar(&c.Cert.Org, "cert-organization", "Example Company Inc", "Certificate organization")
        bCmd.Flags().StringVar(&c.Cert.Unit, "cert-unit", "Operations", "Certificate organizational unit inside of organization")
        bCmd.Flags().StringVar(&c.Cert.Common, "cert-common", "mantl", "Certificate common name")


	root.AddCommand(bCmd)
}

func (c *Config) Bootstrap(args []string) error {
	if err := c.GetSecurityInfo(); err != nil {
		return err
	}

	if c.isBootstrapped {
		var err error
		if c.consulIps, err = consul.GetIps(); err != nil {
			return err
		}
	}

	// Server nodes first
	for _, ip := range(c.serversSplit) {
		if err := c.BootstrapNode(ip, true); err != nil {
			return err
		}
	}
//	for _, ip := range(c.ipsSplit) {
//		bs := &structs.Bootstrap{
//			GossipKey: secInfo.GossipKey,
//			RootToken: secInfo.RootToken,
//			RetryJoin: c.ipsSplit,
//			AdvertiseAddr: ip,
//			Cacert: secInfo.Cert,
//		}
//
//		if !isBootstrapped {
//			if err := WriteToConsul(bs, ip); err != nil {
//				return err
//			}
//		} else {
//			found := false
//			for _, cip := range consulIps {
//				if cip == ip {
//					found = true
//					break
//				}
//			}
//
//			if !found {
//				if err := WriteToConsul(bs, ip); err != nil {
//					return err
//				}
//			}
//		}
//	}

	return nil
}

func (c *Config) BootstrapNode(ip string, isServer bool) error {
	bs := &structs.Bootstrap{
		GossipKey: c.secInfo.GossipKey,
		RootToken: c.secInfo.RootToken,
		RetryJoin: c.serversSplit,
		AdvertiseAddr: ip,
		Cacert: c.secInfo.Cert,
		IsServer: isServer,
	}

	if !c.isBootstrapped {
		if err := c.SendConsulData(ip, bs); err != nil {	
			return err
		}
	}

	return nil
}

func (c *Config) SendConsulData(ip string, bs *structs.Bootstrap) (err error) {
	nc := common.DefaultConfig()

	var conn net.Conn
	for {
		if conn, err = net.Dial(nc.Type, fmt.Sprintf("%s:%d", ip, nc.Port)); err != nil {
			fmt.Println(err)
			fmt.Println("Trying again in 10 seconds")
			time.Sleep(10 * time.Second)
			continue
		}

		break
	}

	defer conn.Close()

	// Connected. Encode bs and send
	e := json.NewEncoder(conn)
	if err := e.Encode(bs); err != nil {
		return err
	}

	result, err := common.RecvResponse(conn)
	if err != nil {
		return err
	}

	if result.Value != common.MSG_SUCCESS {
		return fmt.Errorf(result.Message)
	}

	return nil
	
}
