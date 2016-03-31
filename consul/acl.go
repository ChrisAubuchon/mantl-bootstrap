// +build
//

package consul

import (
	"github.com/hashicorp/consul/api"
)

type AclRule struct {
	Key map[string]*RulePath `json:"key,omitempty"`
	Service map[string]*RulePath `json:"key,omitempty"`
}

type RulePath struct {
	Policy string `json:"policy"`
}


