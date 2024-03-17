package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var AppConfig ApplicationConfig

type ApplicationConfig struct {
	ServerPort            string   `json:"server_port"`
	ZkServers             []string `json:"zk_servers"`
	MongoConnectionString string   `json:"mongo_connection_string"`
	MongoDBName           string   `json:"mongo_db"`
	MongoCollectionName   string   `json:"mongo_collection"`
}

func Initialize() {
	wd, _ := os.Getwd()
	configFile := fmt.Sprintf("%s/src/resources/config.json", wd)
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Panicf("Failed to load config file %s", configFile)
	}
	json.Unmarshal(data, &AppConfig)
}
