package listen

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"

	"github.com/asteris-llc/mantl-bootstrap/cert"
	"github.com/asteris-llc/mantl-bootstrap/common"
	"github.com/asteris-llc/mantl-bootstrap/structs"

	"github.com/spf13/cobra"
)

func Init(root *cobra.Command) {
	listenCmd := &cobra.Command{
		Use: "listen",
		Short: "Listen for bootstrap information",
		Long: "Listen for bootstrap information",
		RunE: func(cmd *cobra.Command, args []string) error {
			return Listen(args)
		},
	}

	root.AddCommand(listenCmd)
}

func Listen(args []string) error {
	c := common.DefaultConfig()

	l, err := net.Listen(c.Type, fmt.Sprintf(":%d", c.Port))
	if err != nil {
		return err
	}

	var bs structs.Bootstrap

	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		defer conn.Close()

		d := json.NewDecoder(conn)

		if err := d.Decode(&bs); err != nil {
			fmt.Println(err)
			conn.Close()
			continue
		}

		if err := common.SendResponse(common.MSG_SUCCESS, "Success", conn); err != nil {
			fmt.Println(err)
		}

		break
	}

	fmt.Printf("%+v\n", bs)

	return ConfigureConsul(&bs)
}

const path = "/etc/consul/consul.json"

func ConfigureConsul(bs *structs.Bootstrap) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("%s does not exist. Creating\n", path)
	}

	cfg := make(map[string]interface{})
	cfg["server"] = bs.IsServer
	cfg["rejoin_after_leave"] = true
	cfg["advertise_addr"] = bs.AdvertiseAddr

	// XXX - These are for demonstration purposes only. Consul will
	// have some sort of proxy in front of it and these two values
	// must be deleted
	cfg["bind_addr"] = bs.AdvertiseAddr
	cfg["client_addr"] = bs.AdvertiseAddr

	if bs.IsServer {
		be := len(bs.RetryJoin)
		if be > 3 {
			be = 3
		}

		cfg["bootstrap_expect"] = be
	}
	cfg["retry_join"] = bs.RetryJoin

	cfg["encrypt"] = bs.GossipKey
	cfg["acl_master_token"] = bs.RootToken

	ioutil.WriteFile("/etc/pki/CA/cacert.pem", []byte(bs.Cacert), 0644)
	ioutil.WriteFile("/etc/pki/CA/cacert.key", []byte(bs.Cakey), 0644)
	ioutil.WriteFile("/etc/pki/ca-trust/source/anchors/cacert.pem", []byte(bs.Cacert), 0600)
	c, err := cert.GenerateCert(bs.Cacert, bs.Cakey, bs.AdvertiseAddr, bs.Domain)
	if err != nil {
		return err
	}

	ioutil.WriteFile("/etc/consul/ssl/consul.cert", c.Cert, 0660)
	ioutil.WriteFile("/etc/consul/ssl/consul.key", c.Key, 0660)
	cmd := exec.Command("chown", "consul:consul", "/etc/consul/ssl/consul.cert", "/etc/consul/ssl/consul.key")
	if err := cmd.Run(); err != nil {
		return err
	}

	cfg["ca_file"] = "/etc/pki/CA/cacert.pem"
	cfg["cert_file"] = "/etc/consul/ssl/consul.cert"
	cfg["key_file"] = "/etc/consul/ssl/consul.key"
	cfg["verify_incoming"] = true
	cfg["verify_outgoing"] = true

	cfg["domain"] = bs.Domain

	bytes, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	ioutil.WriteFile("/etc/consul/consul.json", bytes, 0660)

	return nil
}
