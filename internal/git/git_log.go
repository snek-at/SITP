package git

import (
	"encoding/json"

	"github.com/snek-at/tools"
)

// GetLog returns the git log history of the
// current checked out repository.
//
// completelog=false -> Analyze the last commit
// completelog=true -> Analyze the whole branch
func GetLog(completelog bool) CommitLogStruct {
	var logData CommitLogStruct
	var filelogData map[string]FilesStruct
	var log, filelog string

	if completelog {
		log = tools.ExecuteScript("scripts/git_commit_log.sh", "")
		filelog = tools.ExecuteScript("scripts/git_commit_log_files_stat.sh", "")
	} else {
		log = tools.ExecuteScript("scripts/git_commit_log.sh", "1")
		filelog = tools.ExecuteScript("scripts/git_commit_log_files_stat.sh", "1")
	}

	err := json.Unmarshal([]byte(log), &logData)
	err2 := json.Unmarshal([]byte(filelog), &filelogData)

	if err != nil && err2 != nil {
		panic(err)
	}

	// for {key}, {value} := range {list}

	for key, entry := range logData {
		entry.Files = filelogData[entry.Commit]

		for i := 0; i < len(entry.Files); i++ {

			fileContent := tools.ExecuteScript("scripts/git_commit_changes.sh", entry.Commit, entry.Files[i].Path)
			entry.Files[i].RawChanges = (fileContent)
		}

		logData[key] = entry
	}

	return logData
}

// CommitLogStruct defines the structure of a list
// which items contain information about commits
type CommitLogStruct []struct {
	Commit  string      `json:"commit"`
	Author  string      `json:"author"`
	Date    string      `json:"date"`
	Message string      `json:"message"`
	Files   FilesStruct `json:"files"`
}

// FilesStruct defines the structure of a list
// which items contain information about a specific
// file
type FilesStruct []struct {
	Insertions string `json:"insertions"`
	Deletions  string `json:"deletions"`
	Path       string `json:"path"`
	RawChanges string `json:"raw_changes"`
}
