package ymlmodel

type Equipment struct {
	Equipment map[int32]string `yaml:"equipment"`
	Mapping   map[int32]string `yaml:"mapping"`
}
