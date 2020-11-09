package tools

import (
	"bytes"
	"encoding/json"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

// Info tool is used to get general information
// of the repository
func Info() InformationStruct {
	// git config --get remote.origin.url
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

	return InformationStruct{URL: string(out)}
}

// InformationStruct defines the structure of git
// information and also the keys when generated to
// json
type InformationStruct struct {
	URL string `json:"git_url"`
}

// CommitLog tool is used to get the last n commits
// of the current git repository branch
// n: -1 -> All commits
func CommitLog(depth int) []CommitLogStruct {
	// git config --get remote.origin.url
	format := `--pretty=format:{%n  "commit": "%H",%n  "author": "%aN <%aE>",%n  "date": "%ad",%n  "message": "%f"%n},`
	logDepth := "-n "

	var cmd *exec.Cmd

	if depth > 0 {
		d := strconv.Itoa(depth)
		logDepth += d

		cmd = exec.Command("git", "log", logDepth, format)
	} else {
		cmd = exec.Command("git", "log", format)
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

	/*[OUT Sample]
	{
		"commit": "5f8e958de03f5468ece45cb622f7a41460962df3",
		"author": "schettn <nicoschett@icloud.com>",
		"date": "Tue Jul 28 20:30:17 2020 +0200",
		"message": "Add-go.sum"
	},
	*/

	// The last character of 'out' needs to be removed in order to
	// be valid JSON
	out = out[:len(out)-1]

	/*[OUT Sample]
	{
		"commit": "5f8e958de03f5468ece45cb622f7a41460962df3",
		"author": "schettn <nicoschett@icloud.com>",
		"date": "Tue Jul 28 20:30:17 2020 +0200",
		"message": "Add-go.sum"
	}
	*/

	// Append list identifiers before and after 'out'
	prefix := []byte("[")
	suffix := []byte("]")

	out = append(out, suffix...)
	out = append(prefix, out...)

	/*[OUT Sample]
	[{
		"commit": "5f8e958de03f5468ece45cb622f7a41460962df3",
		"author": "schettn <nicoschett@icloud.com>",
		"date": "Tue Jul 28 20:30:17 2020 +0200",
		"message": "Add-go.sum"
	}]
	*/

	// Escape backslashes
	out = bytes.ReplaceAll(out, []byte(`\`), []byte(`\\`))

	// Return list
	var logData []CommitLogStruct

	err = json.Unmarshal([]byte(out), &logData)

	if err != nil {
		panic(err)
	}

	return logData
}

// CommitLogStruct defines the structure of a list
// which items contain information about commits
type CommitLogStruct struct {
	Commit  string       `json:"commit"`
	Author  string       `json:"author"`
	Date    string       `json:"date"`
	Message string       `json:"message"`
	Files   []FileStruct `json:"files"`
}

// CommitLogFiles tool is used to get the last n commits
// of the current git repository branch, containing a for
// each commit
// n: -1 -> All commits
func CommitLogFiles(depth int) []CommitFilesStruct {
	format := `--pretty=format:'%H'`
	logDepth := "-n "

	var cmd *exec.Cmd

	if depth > 0 {
		d := strconv.Itoa(depth)
		logDepth += d

		cmd = exec.Command("git", "log", logDepth, "--numstat", format)
	} else {
		cmd = exec.Command("git", "log", "--numstat", format)
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

	/*[OUT Sample]
	'5f8e958de03f5468ece45cb622f7a41460962df3'
	10	0	go.sum

	'5f8e958de03f5468ece45cb622f7a41460962df3'
	10	0	go.sum

	*/

	// Split by empty line
	commitArray := strings.Split(string(out), "\n\n")

	var commitLogData []CommitFilesStruct

	for _, commit := range commitArray {
		commitInfo := strings.Split(string(commit), "\n")

		// Remove first and last character from commit hash
		/*[commitHash Sample]
		'5f8e958de03f5468ece45cb622f7a41460962df3'

		--> 5f8e958de03f5468ece45cb622f7a41460962df3
		*/
		commitHash := commitInfo[0][:len(commitInfo[0])-1][1:]

		commitFiles := commitInfo[1:]

		commitLogEntryData := CommitFilesStruct{Commit: commitHash}

		for _, commitFile := range commitFiles {
			commitFilesInfo := strings.Split(string(commitFile), "\t")

			if len(commitFilesInfo) > 1 {
				insertions, deletions, path := commitFilesInfo[0], commitFilesInfo[1], commitFilesInfo[2]
				commitLogEntryData.Files = append(commitLogEntryData.Files, FileStruct{Insertions: insertions, Deletions: deletions, Path: path})
			}
		}

		commitLogData = append(commitLogData, commitLogEntryData)
	}
	return commitLogData
}

// CommitFilesStruct defines the structure of a commit
// containing a changelog of all files
type CommitFilesStruct struct {
	Commit string       `json:"commitHash"`
	Files  []FileStruct `json:"files"`
}

// FileStruct defines the structure of a specific changelog
// entry
type FileStruct struct {
	Insertions string `json:"insertions"`
	Deletions  string `json:"deletions"`
	Path       string `json:"path"`
	RawChanges string `json:"raw_changes"`
}

// CommitLogChanges tool is used to get the specific
// content which was changed by a given commit and file
func CommitLogChanges(version, file string) string {
	cmd := exec.Command("git", "show", "--color", version, "--", file)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

	lines := strings.Split(string(out), "\n")

	changes := ""

	for _, line := range lines {
		if strings.HasPrefix(line, "\033[32m+") || strings.HasPrefix(line, "\033[31m-") {
			changes += line + "\n"
		}
	}
	return changes
}
