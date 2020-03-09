package filestream

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

type ConfigFilePath struct {
	Verbose    bool          `json:"verbose"`
	Server     string        `json:"server"`
	ListenPort string        `json:"listen_port"`
	Cipher     string        `json:"cipher"`
	Password   string        `json:"password"`
	Plugin     string        `json:"plugin"`
	PluginOpts string        `json:"plugin_opts"`
	UDPTimeout time.Duration `json:"udp_timeout"`
}

func (c *ConfigFilePath) ParseConfigFile(filePath string) error {
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(fileData, c)
	return err
}
