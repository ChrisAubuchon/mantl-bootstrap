package agentacl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/asteris-llc/mantl-bootstrap/common"
	"github.com/asteris-llc/mantl-bootstrap/consul"

	"github.com/spf13/cobra"
	"github.com/hashicorp/consul/api"
)

var agentAcl = `{"key":{"":{"Policy":"write"}},"service":{"":{"Policy":"write"}}}`

func Init(root *cobra.Command) {
	aclCmd := &cobra.Command{
		Use:   "agent-acl",
		Short: "Bootstrap the agent ACL token",
		Long: "Bootstrap the agent ACL token",
		RunE: func(cmd *cobra.Command, args []string) error {
			return AgentAcl(args)
		},
	}

	root.AddCommand(aclCmd)
}

func AgentAcl(args []string) error {
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

	agentToken, ok := rval["acl_token"].(string)
	if !ok {
		return fmt.Errorf("acl_token not in %s", common.ConsulStaticPath)
	}

	client := consul.NewClient().ACL()
	opts := consul.QueryOptions(rootToken)

	acl, _, err := client.Info(agentToken, opts)
	if err != nil {
		return err
	}

	if acl != nil {
		// ACL already created. Return
		return nil
	}

	entry := &api.ACLEntry{
		ID: agentToken,
		Name: "Agent policy",
		Type: api.ACLClientType,
		Rules: agentAcl,
	}

	_, err = client.Update(entry, consul.WriteOptions(rootToken))
	if err != nil {
		return err
	}

	return nil
}
