package common

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

func NewConsulConfig() *ConsulStatic {
	return &ConsulStatic{
		Encrypt: "",
	}
}
