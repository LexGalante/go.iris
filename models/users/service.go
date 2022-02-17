package users

import (
	"errors"

	"github.com/lexgalante/go.iris/utils"
	"gopkg.in/dealancer/validate.v2"
)

//CryptPassword -> crypt password
func (u *User) CryptPassword() error {
	password, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	u.Password = password

	return nil
}

//GrantRole -> give a role to user
func (u *User) GrantRole(role string) {
	if !utils.ContainsInSlice(role, u.Roles) {
		u.Roles = append(u.Roles, role)
	}
}

//InRole -> check user has a role
func (u User) InRole(role string) bool {
	return utils.ContainsInSlice(role, u.Roles)
}

//Validate -> custom validations for github.com/dealancer/validate
func (u User) Validate() error {
	if len(u.Roles) == 0 {
		return errors.New("must at least once role 'admin' or 'guest'")
	}

	return nil
}

//ExecuteValidations -> execute a pipeline of validations
func (u *User) ExecuteValidations() error {
	return validate.Validate(&u)
}
