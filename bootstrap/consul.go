package bootstrap

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/asteris-llc/mantl-bootstrap/cert"
	"github.com/asteris-llc/mantl-bootstrap/common"

	"github.com/satori/go.uuid"
)

// GenConsulConfig()
// Return the Consul security information needed to add a node to an
// existing Consul cluster.
//
func (b *Bootstrap) GenConsulConfig() error {
	if err := b.ReadConsulConfig(); err != nil {
		return err
	}

	if b.Consul.Encrypt == "" {
		b.Config.Set(tag_isBootstrapped, false)

		return b.genConsulConfig()
	}
	b.Config.Set(tag_isBootstrapped, true)

	if b.Consul.CaFile == "" {
		return fmt.Errorf("No ca_file in %s", common.ConsulStaticPath)
	}

	cafile, err := ioutil.ReadFile(b.Consul.CaFile)
	if err != nil {
		return err
	}
	b.caFile = cafile

	cakey, err := ioutil.ReadFile(common.CaKeyPath)
	if err != nil {
		return err
	}
	b.caKey = cakey

	return nil
}

func (b *Bootstrap) ReadConsulConfig() error {
	if _, err := os.Stat(common.ConsulStaticPath); os.IsNotExist(err) {
		return nil
	}

	consulJson, err := ioutil.ReadFile(common.ConsulStaticPath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(consulJson, &b.Consul); err != nil {
		return err
	}

	return nil
}

func (b *Bootstrap) genConsulConfig() error {
	// Generate gossip key
	key := make([]byte, 16)
	n, err := rand.Reader.Read(key)
	if err != nil {
		return err
	}

	if n != 16 {
		return fmt.Errorf("Couldn't read enough entropy.")
	}

	cert, err := cert.GenerateCaCert(b.Config)
	if err != nil {
		return err
	}

	b.Consul.CaFile = common.CaCertPath
	b.Consul.AclToken = uuid.NewV4().String()
	b.Consul.AclMasterToken = uuid.NewV4().String()
	b.Consul.Encrypt = base64.StdEncoding.EncodeToString(key)
	b.Consul.AclDefaultPolicy = "deny"
	b.Consul.AclDownPolicy = "allow"
	b.Consul.RejoinAfterLeave = true
	b.Consul.VerifyIncoming = true
	b.Consul.VerifyOutgoing = true
	b.caFile = cert.Cert
	b.caKey = cert.Key

	return nil
}
