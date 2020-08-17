package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/snek-at/tools"
	"gopkg.in/yaml.v2"

	"github.com/snek-at/internal/client"
	"github.com/snek-at/internal/git"
)

// YamlConfig is exported.
type YamlConfig struct {
	Connection struct {
		Address string `yaml:"address"`
	} `yaml:"Connection"`
	Delicacies struct {
		Fetchmode string `yaml:"fetchmode"`
		Blocksize int    `yaml:"blocksize"`
	} `yaml:"Delicacies"`
}

func getEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	log.Fatal(errors.New(key + " not set"))
	return ""
}

func main() {
	// Get general git information
	info := git.GetInformation()
	// Full log status
	completelog := false

	// Parse yaml config
	configPath := "./sitp.yaml"

	var config YamlConfig
	var isErr = false
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		isErr = true
	} else {
		err = yaml.Unmarshal(yamlFile, &config)
		if err != nil {
			isErr = true
		}
	}

	if isErr == true {
		config.Connection.Address = "http://localhost:8000/graphql"
		config.Delicacies.Blocksize = 20
		config.Delicacies.Fetchmode = "LCO" // LCO -> Lastest commit only
	}

	fmt.Printf("Result: %v\n", config)
	// Check mode env variable

	switch mode := config.Delicacies.Fetchmode; mode {
	case "COMPLETE":
		completelog = true
	case "LCO":
		completelog = false
	}

	// Initialize log reader channel
	reader, err := git.GetLog(completelog)

	if err != nil {
		fmt.Println("Error in commit item ocurred")
	}

	var buffer CommitLog

	for item := range reader {

		buffer = append(buffer, item)

		if len(buffer) > config.Delicacies.Blocksize {
			send(DataStruct{Git: info, Log: buffer}, config)
			buffer = nil
		}
	}

	if len(buffer) > 0 {
		send(DataStruct{Git: info, Log: buffer}, config)
		buffer = nil
	}
}

// Tranforms data to json and sends it to OPS
func send(data DataStruct, config YamlConfig) {
	// Convert struct to json
	bufx, _ := json.Marshal(data)

	// Get token from env
	token := getEnv("pipeline_token")

	// Send json to OPS
	client.SendToOPS(config.Connection.Address, token, string(bufx))
}

// CommitLog defines a list structure of commit items
type CommitLog = []tools.CommitLogStruct

// DataStruct defines the combined structure of CommitLog and Git
type DataStruct struct {
	Git git.BasicInformation
	Log git.CommitLog
}
