package main

import (
	"encoding/json"
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

	// Get git log of checked out branch
	log := git.GetLog(completelog)
	// Convert struct to json
	data := DataStruct{Git: info, Log: log}

	bufx, _ := json.Marshal(data)

	// Send json to OPS
	client.SendToOPS(string(bufx))
}

type a = tools.InformationStruct

type DataStruct struct {
	Git git.BasicInformation
	Log git.CommitLog
}
