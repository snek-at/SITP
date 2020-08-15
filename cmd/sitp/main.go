package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/snek-at/tools"

	"github.com/snek-at/internal/client"
	"github.com/snek-at/internal/git"
)

func main() {
	// Get general git information
	info := git.GetInformation()
	// Full log status
	completelog := false
	// Check mode env variable

	switch mode := os.Getenv("MODE"); mode {
	case "COMPLETE":
		completelog = true
	}

	// Initialize log reader channel
	reader, err := git.GetLog(completelog)

	if err != nil {
		fmt.Println("Error in commit item ocurred")
	}

	var buffer CommitLog

	for item := range reader {

		buffer = append(buffer, item)

		if len(buffer) > 15 {
			send(DataStruct{Git: info, Log: buffer})
			buffer = nil
		}
	}

	if len(buffer) > 0 {
		send(DataStruct{Git: info, Log: buffer})
		buffer = nil
	}
}

// Tranforms data to json and sends it to OPS
func send(data DataStruct) {
	// Convert struct to json
	bufx, _ := json.Marshal(data)

	// Send json to OPS
	client.SendToOPS(string(bufx))
}

// CommitLog defines a list structure of commit items
type CommitLog = []tools.CommitLogStruct

// DataStruct defines the combined structure of CommitLog and Git
type DataStruct struct {
	Git git.BasicInformation
	Log git.CommitLog
}
