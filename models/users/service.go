package users

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/lexgalante/go.iris/utils"
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

//VerifyPassword -> check password
func (u *User) VerifyPassword(password string) bool {
	return utils.VerifyPassword(u.Password, password)
}

//CreateAccessToken -> return jwk user token
func (u *User) CreateAccessToken() (*UserAccessToken, error) {
	var userAccessToken UserAccessToken

	expireAt := time.Now().Add(time.Hour * 1).Unix()
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expireAt,
		},
		ID:     u.ID.String(),
		Email:  u.Email,
		Roles:  u.Roles,
		Active: u.Active,
	})

	token, err := claims.SignedString([]byte(os.Getenv("APP_SECRET")))
	if err != nil {
		return &userAccessToken, err
	}

	userAccessToken.TokenType = "Bearer"
	userAccessToken.AccessToken = token
	userAccessToken.ExpiresAt = expireAt

	return &userAccessToken, nil
}

//ValidateAcessToken -> validate access token see https://pkg.go.dev/github.com/golang-jwt/jwt#example-Parse-Hmac
func ValidateAcessToken(accessToken string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(accessToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("APP_SECRET")), nil
	})
}

//GrantRole -> give a role to user
func (u *UserClaims) GrantRole(role string) {
	if !utils.ContainsStringInSlice(role, u.Roles) {
		u.Roles = append(u.Roles, role)
	}
}

//InRole -> check user has a role
func (u UserClaims) InRole(role string) bool {
	return utils.ContainsStringInSlice(role, u.Roles)
}
