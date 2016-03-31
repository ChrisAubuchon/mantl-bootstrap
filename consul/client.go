package consul

import (
	api "github.com/hashicorp/consul/api"
)

func GetIps() ([]string, error) {
	// Get the ip list from the local hosts Consul instance
	consul, err := api.NewClient(api.DefaultConfig())
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

func NewClient() *api.Client {
	client, _ := api.NewClient(api.DefaultConfig())

	return client
}

func QueryOptions(token string) *api.QueryOptions {
	return &api.QueryOptions{
		Token: token,
	}
}

func WriteOptions(token string) *api.WriteOptions {
	return &api.WriteOptions{
		Token: token,
	}
}
