package vault

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/asteris-llc/mantl-bootstrap/common"

	"github.com/spf13/cobra"
)

const (
	VaultConfigPath = "/etc/vault/vault.json"
)

func Init(root *cobra.Command) {
	vaultCmd := &cobra.Command{
		Use:   "vault",
		Short: "Bootstrap vault",
		Long: "Bootstrap vault",
		RunE: func(cmd *cobra.Command, args []string) error {
			return VaultBootstrap(args)
		},
	}

	root.AddCommand(vaultCmd)
}

type VaultConfig struct {
	Backend map[string]Backend `json:"backend"`
	Listener map[string]Listener `json:"listener"`
}

type Backend struct {
	Address string `json:"address"`
	Path string `json:"path"`
	Scheme string `json:"scheme"`
	Token string `json:"token"`
	AdvertiseAddr string `json:"advertise_addr"`
}

type Listener struct {
	Address string `json:"address"`
	CertFile string `json:"tls_cert_file"`
	KeyFile string `json:"tls_key_file"`
}

func VaultBootstrap(args []string) error {
	if _, err := os.Stat(VaultConfigPath); err == nil {
		// Vault configuration exists. Exit.
		return nil
	}

	consulJson, err := ioutil.ReadFile(common.ConsulStaticPath)
	if err != nil {
		return err
	}

	rval := make(map[string]interface{})
	if err := json.Unmarshal(consulJson, &rval); err != nil {
		return err
	}

	rootToken, ok := rval["acl_master_token"].(string)
	if !ok {
		return fmt.Errorf("acl_master_token not in %s", common.ConsulStaticPath)
	}

	advertiseAddr, ok := rval["advertise_addr"].(string)
	if !ok {
		return fmt.Errorf("advertiseAddr not in %s", common.ConsulStaticPath)
	}

	vaultConfig := VaultConfig{
		Backend: map[string]Backend{
			"consul": {
				Address: "127.0.0.1.8500",
				Path: "vault",
				Scheme: "http",
				Token: rootToken,
				AdvertiseAddr: fmt.Sprintf("https://%s", advertiseAddr),
			},
		},
		Listener: map[string]Listener{
			"tcp": {
				Address: "0.0.0.0:8200",
				CertFile: "/etc/pki/tls/certs/host.cert",
				KeyFile: "/etc/pki/tls/private/host.key",
			},
		},
	}

	bytes, err := json.MarshalIndent(vaultConfig, "", "  ")
	if err != nil {
		return err
	}
	ioutil.WriteFile(VaultConfigPath, bytes, 0660)

	return nil
}
