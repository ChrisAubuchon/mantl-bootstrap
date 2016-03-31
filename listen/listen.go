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
		Use:   "listen",
		Short: "Listen for bootstrap information",
		Long:  "Listen for bootstrap information",
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

type ConsulStatic struct {
	AclDatacenter    string   `json:"acl_datacenter"`
	AclDefaultPolicy string   `json:"acl_default_policy"`
	AclDownPolicy    string   `json:"acl_down_policy"`
	AclMasterToken   string   `json:"acl_master_token"`
	AclToken         string   `json:"acl_token"`
	AdvertiseAddr    string   `json:"advertise_addr"`
	BootstrapExpect  int      `json:"bootstrap_expect,omitempty"`
	CaFile           string   `json:"ca_file"`
	CertFile         string   `json:"cert_file"`
	Datacenter       string   `json:"datacenter"`
	Domain           string   `json:"domain'"`
	Encrypt          string   `json:"encrypt"`
	KeyFile          string   `json:"key_file"`
	RejoinAfterLeave bool     `json:"rejoin_after_leave"`
	RetryJoin        []string `json:"retry_join"`
	Server           bool     `json:"server"`
	VerifyIncoming   bool     `json:"verify_incoming"`
	VerifyOutgoing   bool     `json:"verify_outgoing"`
}

func ConfigureConsul(bs *structs.Bootstrap) error {
	if _, err := os.Stat(common.ConsulStaticPath); os.IsNotExist(err) {
		fmt.Printf("%s does not exist. Creating\n", common.ConsulStaticPath)
	}

	cfg := ConsulStatic{
		AclDatacenter:    bs.Datacenter,
		AclDefaultPolicy: "deny",
		AclDownPolicy:    "allow",
		AclMasterToken:   bs.RootToken,
		AclToken:         bs.AgentToken,
		AdvertiseAddr:    bs.AdvertiseAddr,
		CaFile:           "/etc/pki/CA/cacert.pem",
		CertFile:         "/etc/pki/tls/certs/host.cert",
		Datacenter:       bs.Datacenter,
		Domain:           bs.Domain,
		Encrypt:          bs.GossipKey,
		KeyFile:          "/etc/pki/tls/private/host.key",
		RejoinAfterLeave: true,
		RetryJoin:        bs.RetryJoin,
		Server:           bs.IsServer,
		VerifyIncoming:   true,
		VerifyOutgoing:   true,
	}

	if bs.IsServer {
		be := len(bs.RetryJoin)
		if be > 3 {
			be = 3
		}

		cfg.BootstrapExpect = be
	}

	ioutil.WriteFile("/etc/pki/CA/cacert.pem", []byte(bs.Cacert), 0644)
	ioutil.WriteFile("/etc/pki/CA/cacert.key", []byte(bs.Cakey), 0644)
	ioutil.WriteFile("/etc/pki/ca-trust/source/anchors/cacert.pem", []byte(bs.Cacert), 0600)
	c, err := cert.GenerateCert(bs.Cacert, bs.Cakey, bs.AdvertiseAddr, bs.Domain)
	if err != nil {
		return err
	}

	ioutil.WriteFile("/etc/pki/tls/certs/host.cert", c.Cert, 0644)
	ioutil.WriteFile("/etc/pki/tls/private/host.key", c.Key, 0644)

	bytes, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	ioutil.WriteFile(common.ConsulStaticPath, bytes, 0660)
	cmd := exec.Command("chown", "consul:consul", common.ConsulStaticPath)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
