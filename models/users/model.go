package users

import "github.com/kamva/mgm/v3"

//User -> represent a user
type User struct {
	mgm.DefaultModel `bson:",inline"`
	Email            string   `json:"email" bson:"email" validate:"gte=10 & lte=250 & format=email"`
	Password         string   `json:"password" bson:"password" validate:"empty=false"`
	Roles            []string `json:"roles" bson:"roles" validate:"empty=false"`
}
