package structs

type Bootstrap struct {
	GossipKey     string   `json:"gossip_key"`
	RootToken     string   `json:"root_token"`
	RetryJoin     []string `json:"retry_join"`
	AdvertiseAddr string   `json:"advertise_addr"`
	Cacert        string   `json:"cacert"`
	IsServer      bool     `json:"is_server"`
}
