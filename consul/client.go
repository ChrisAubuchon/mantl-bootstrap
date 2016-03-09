package consul

import (
	"fmt"
	"time"

	api "github.com/hashicorp/consul/api"
)

func GetIps(ip string) ([]string, error) {
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:8500", ip) 

	// Get the ip list from the local hosts Consul instance
	consul, err := api.NewClient(config)
	if err != nil {
		return []string{}, err
	}

	catalog := consul.Catalog()
	nodes, _, err := catalog.Nodes(&api.QueryOptions{})
	if err != nil {
		return []string{}, err
	}

	rval := make([]string, len(nodes))
	for i, node := range(nodes) {
		rval[i] = node.Address
	}

	return rval, nil
}

func Write(ip string, bytes []byte) error {
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:8500", ip)

	consul, err := api.NewClient(config)
	if err != nil {
		return err
	}

	agent := consul.Agent()
	for {
		_, err := agent.Self()
		if err == nil {
			break
		}

		time.Sleep(5)
	}

	kv := consul.KV()
	kvp := api.KVPair{
		Key: "bootstrap/consul",
		Value: bytes,
	}

	if _, err := kv.Put(&kvp, &api.WriteOptions{}); err != nil {
		return err
	}

	return nil
}
