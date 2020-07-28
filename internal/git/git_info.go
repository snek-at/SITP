package git

import (
	"encoding/json"
	"fmt"

	"github.com/snek-at/tools"
)

// GetInformation returns basic information of the
// current checked out git repository
func GetInformation() InformationStruct {
	var infoData InformationStruct

	infoStr := tools.ExecuteScript("scripts/git_info.sh", "")

	fmt.Println(infoStr)

	err := json.Unmarshal([]byte(infoStr), &infoData)

	if err != nil {
		panic(err)
	}

	return infoData
}

// InformationStruct defines the structure of git
// information and also the keys when generated to
// json
type InformationStruct struct {
	URL string `json:"git_url"`
}
