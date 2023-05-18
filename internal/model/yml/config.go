package ymlmodel

type Config struct {
	Equipment map[int]string `yaml:"equipment"`
	Mapping   map[int]string `yaml:"mapping"`
}
