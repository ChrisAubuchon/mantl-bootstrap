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
