package main

import (
	"encoding/json"
	"os"

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

	// Get full git log of checked out branch
	log := git.GetLog(completelog)

	// Convert struct to json
	data := dataStruct{Git: info, Log: log}
	bufx, _ := json.Marshal(data)

	// Send json to OPS
	client.SendToOPS(string(bufx))
}

type dataStruct struct {
	Git git.InformationStruct
	Log git.CommitLogStruct
}
