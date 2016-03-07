package bootstrap

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/satori/go.uuid"
)

const Path = "/etc/consul/consul.json"

type SecInfo struct {
	RootToken string
	GossipKey string
	Cert      string
}


// GetSecurityInfo()
// Return the Consul security information needed to add a node to an
// existing Consul cluster. Return a boolean flag to indicate whether
// the current node is already bootstrapped
//
func (c *Config) GetSecurityInfo() (err error) {
	cconfig := ReadConsulConfig()
	if cconfig == nil {
		c.isBootstrapped = false
		c.secInfo, err = c.GenerateSecInfo()

		return err
	}
	c.isBootstrapped = true

	cert, err := ioutil.ReadFile(cconfig["cacert"].(string))
	if err != nil {
		return err
	}

	c.secInfo = &SecInfo{
		RootToken: cconfig["acl_master_token"].(string),
		GossipKey: cconfig["encrypt"].(string),
		Cert:      string(cert),
	}

	return nil
}

func ReadConsulConfig() map[string]interface{} {
	if _, err := os.Stat(Path); os.IsNotExist(err) {
		return nil
	}

	return nil
}

func (c *Config) GenerateSecInfo() (*SecInfo, error) {
	// Generate gossip key
	key := make([]byte, 16)
	n, err := rand.Reader.Read(key)
	if err != nil {
		return nil, err
	}

	if n != 16 {
		return nil, fmt.Errorf("Couldn't read enough entropy.")
	}

	cert, err := GenerateCaCert(c.Cert)
	if err != nil {
		return nil, err
	}

	si := &SecInfo{
		RootToken: uuid.NewV4().String(),
		GossipKey: base64.StdEncoding.EncodeToString(key),
		Cert: cert,
	}

	return si, nil
}

