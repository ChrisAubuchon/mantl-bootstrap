package bootstrap

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/asteris-llc/mantl-bootstrap/cert"

	"github.com/satori/go.uuid"
)

const Path = "/etc/consul/consul.json"

type SecInfo struct {
	RootToken string
	GossipKey string
	Cert      string
	Key       string
}


// GetSecurityInfo()
// Return the Consul security information needed to add a node to an
// existing Consul cluster. Return a boolean flag to indicate whether
// the current node is already bootstrapped
//
func (c *Config) GetSecurityInfo() (err error) {
	cconfig, err := ReadConsulConfig()
	if err != nil {
		return err
	}
	if cconfig == nil {
		c.isBootstrapped = false
		c.secInfo, err = c.GenerateSecInfo()

		return err
	}
	c.isBootstrapped = true

	fmt.Println("isBootstrapped == true")
	fmt.Printf("%+v\n", cconfig)

	cert, err := ioutil.ReadFile(cconfig["ca_file"].(string))
	if err != nil {
		return err
	}

	key, err := ioutil.ReadFile("/etc/pki/CA/cacert.key")
	if err != nil {
		return err
	}

	c.secInfo = &SecInfo{
		RootToken: cconfig["acl_master_token"].(string),
		GossipKey: cconfig["encrypt"].(string),
		Cert:      string(cert),
		Key:       string(key),
	}

	return nil
}

func ReadConsulConfig() (map[string]interface{}, error) {
	if _, err := os.Stat(Path); os.IsNotExist(err) {
		return nil, nil
	}

	consulJson, err := ioutil.ReadFile(Path)
	if err != nil {
		return nil, err
	}

	rval := make(map[string]interface{})

	if err := json.Unmarshal(consulJson, &rval); err != nil {
		return nil, err
	}

	if _, ok := rval["encrypt"]; !ok {
		return nil, nil
	}

	return rval, nil
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

	cert, err := cert.GenerateCaCert(c.Cert)
	if err != nil {
		return nil, err
	}

	si := &SecInfo{
		RootToken: uuid.NewV4().String(),
		GossipKey: base64.StdEncoding.EncodeToString(key),
		Cert: string(cert.Cert),
		Key: string(cert.Key),
	}

	return si, nil
}

