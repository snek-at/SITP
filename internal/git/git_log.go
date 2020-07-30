package git

import (
	"github.com/snek-at/tools"
)

// GetLog returns the git log history of the
// current checked out repository.
//
// completeLog=false -> Analyze the last commit
// completeLog=true -> Analyze the whole branch
func GetLog(completeLog bool) CommitLog {

	var commitLog CommitLog
	var commitLogFiles CommitFiles

	if completeLog {
		commitLog = tools.CommitLog(-1)
		commitLogFiles = tools.CommitLogFiles(-1)
	} else {
		commitLog = tools.CommitLog(1)
		commitLogFiles = tools.CommitLogFiles(1)
	}

	for commitIndex, commit := range commitLog {
		// fmt.Println(commit)
		for _, logFiles := range commitLogFiles {
			if logFiles.Commit == commit.Commit {
				commitLog[commitIndex].Files = logFiles.Files

				for logFileIndex, logFile := range logFiles.Files {
					commitLog[commitIndex].Files[logFileIndex].RawChanges = tools.CommitLogChanges(commit.Commit, logFile.Path)
				}
			}
		}
	}
	return commitLog
}

type CommitLog = []tools.CommitLogStruct
type CommitFiles = []tools.CommitFilesStruct
