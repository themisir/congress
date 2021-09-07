package models

type Config struct {
	Congress CongressConfig `yaml:"congress"`
	Rules    []Rule         `yaml:"rules"`
}

type CongressConfig struct {
	Ip    string `yaml:"ip"`
	Proxy struct {
		Enabled bool `yaml:"enabled"`
		Port    uint `yaml:"port"`
	} `yaml:"proxy"`
	Dns struct {
		Enabled  bool   `yaml:"enabled"`
		Port     uint   `yaml:"port"`
		Fallback string `yaml:"fallback"`
	} `yaml:"dns"`
}
