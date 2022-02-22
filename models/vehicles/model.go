package vehicles

import "github.com/kamva/mgm/v3"

//Vehicle -> represents a vehicle
type Vehicle struct {
	mgm.DefaultModel `bson:",inline"`
	License          string `json:"license" bson:"license" validate:"required,licenseplate"`
	UserID           string `json:"user_id" bson:"user_id"`
	Name             string `json:"name" bson:"name" validate:"required,gte=2,lte=50"`
	Model            string `json:"model" bson:"model" validate:"required,gte=2,lte=20"`
	YearModel        int    `json:"year_model" bson:"year_model" validate:"required,min=1900,max=2023"`
	YearManufactory  int    `json:"year_manufactory" bson:"year_manufactory" validate:"required,ltefield=YearModel"`
	Color            string `json:"color" bson:"color" validate:"required,gte=4,lte=25"`
	Active           bool   `json:"active" bson:"active"`
}
