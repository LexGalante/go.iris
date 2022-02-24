package traces

import (
	"time"

	"github.com/kamva/mgm/v3"
)

//Trace -> represent a vehicle trace
type Trace struct {
	mgm.DefaultModel `bson:",inline"`
	License          string    `json:"license,omitempty" bson:"license" validate:"licenseplate"`
	Latitude         float64   `json:"latitude" bson:"latitude" validate:"required,latitude"`
	Longitude        float64   `json:"longitude" bson:"longitude" validate:"required,longitude"`
	DateTime         time.Time `json:"date_time" bson:"date_time"`
}
