package vehicles

import (
	"errors"

	"github.com/lexgalante/go.iris/utils"
)

//Validate -> validate struct
func (v Vehicle) Validate() error {
	if utils.IsValidLicensePlate(v.License) {
		return errors.New("invalid license plate")
	}

	return nil
}
