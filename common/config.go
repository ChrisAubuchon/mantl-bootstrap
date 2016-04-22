package common

type Config struct {
	Port int
	Type string
}

func DefaultConfig() *Config {
	return &Config{
		Port: 3242,
		Type: "tcp",
	}
}

const (
	ConsulStaticPath = "/etc/consul/static.json"

	CaCertPath = "/etc/pki/CA/ca.cert"
	CaKeyPath = "/etc/pki/CA/ca.key"
	CaCertAnchorPath = "/etc/pki/ca-trust/source/anchors/ca.cert"
	HostCertPath = "/etc/pki/tls/certs/host.cert"
	HostKeyPath = "/etc/pki/tls/private/host.key"
)
