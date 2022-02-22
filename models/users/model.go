package users

import (
	"github.com/golang-jwt/jwt"
	"github.com/kamva/mgm/v3"
)

//User -> represent a user
type User struct {
	mgm.DefaultModel `bson:",inline"`
	Email            string   `json:"email" bson:"email" validate:"required,email,lte=250,lowercase"`
	Password         string   `json:"password,omitempty" bson:"password" validate:"required,gte=6,lte=20"`
	Roles            []string `json:"roles" bson:"roles"`
	Active           bool     `json:"active" bson:"active"`
}

//UserClaims -> custom claims for golang-jwt
type UserClaims struct {
	*jwt.StandardClaims
	ID     string   `json:"id"`
	Email  string   `json:"email"`
	Roles  []string `json:"roles"`
	Active bool     `json:"active" bson:"active"`
}

//UserAccessToken -> represent a jwk user's personal token
type UserAccessToken struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
}
