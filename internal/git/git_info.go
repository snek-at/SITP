package git

import (
	"github.com/snek-at/tools"
)

// GetInformation returns basic information of the
// current checked out git repository
func GetInformation() BasicInformation {
	return tools.Info()
}

type BasicInformation = tools.InformationStruct
