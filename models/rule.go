package models

type Rule struct {
	Host           string     `yaml:"host"`
	DefaultBackend string     `yaml:"defaultBackend,omitempty"`
	Paths          []RulePath `yaml:"paths,omitempty"`
}

type RulePath struct {
	Path    string `yaml:"path"`
	Backend string `yaml:"backend"`
}
