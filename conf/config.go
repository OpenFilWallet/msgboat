package conf

import (
	"github.com/spf13/viper"
)

var config *viper.Viper

func LocalConfig() error {
	config = viper.New()
	config.SetConfigName("config")
	config.SetConfigType("yaml")
	config.AddConfigPath(".")
	config.AddConfigPath("./conf")
	config.AddConfigPath("../conf")
	err := config.ReadInConfig()
	return err
}

func GetNodes() map[string]string {
	nodesConf := config.GetStringMap("nodes")

	var nodes = make(map[string]string)
	for name, rpcAddr := range nodesConf {
		nodes[name] = rpcAddr.(string)
	}

	return nodes
}
