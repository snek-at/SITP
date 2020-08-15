package git

import (
	"github.com/snek-at/tools"
)

// GetLog returns the git log history of the
// current checked out repository.
//
// completeLog=false -> Analyze the last commit
// completeLog=true -> Analyze the whole branch
func GetLog(completeLog bool) (<-chan tools.CommitLogStruct, error) {

	var commitLog CommitLog
	var commitLogFiles CommitFiles

	if completeLog {
		commitLog = tools.CommitLog(-1)
		commitLogFiles = tools.CommitLogFiles(-1)
	} else {
		commitLog = tools.CommitLog(1)
		commitLogFiles = tools.CommitLogFiles(1)
	}

	chnl := make(chan tools.CommitLogStruct)
	go func() {
		for commitIndex, commit := range commitLog {
			for _, logFiles := range commitLogFiles {
				if logFiles.Commit == commit.Commit {
					commitLog[commitIndex].Files = logFiles.Files

					for logFileIndex, logFile := range logFiles.Files {
						commitLog[commitIndex].Files[logFileIndex].RawChanges = tools.CommitLogChanges(commit.Commit, logFile.Path)
					}
				}
			}
			chnl <- commitLog[commitIndex]
		}
		// Ensure that at the end of the loop we close the channel!
		close(chnl)
	}()
	return chnl, nil
}

type CommitLog = []tools.CommitLogStruct
type CommitFiles = []tools.CommitFilesStruct
