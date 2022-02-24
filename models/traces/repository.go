package traces

import (
	"context"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Save -> create or update model
func (t *Trace) Save() (primitive.ObjectID, error) {
	collection := mgm.Coll(t)

	if t.ID == primitive.NilObjectID {
		err := collection.Create(t)
		if err != nil {
			return t.ID, err
		}
	}

	err := collection.Update(t)
	if err != nil {
		return t.ID, err
	}

	return t.ID, nil
}

//Remove -> delete model
func (t *Trace) Remove(license string) (int64, error) {
	collection := mgm.Coll(t)

	result, err := collection.DeleteMany(context.TODO(), bson.M{"license": license})
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}

//FindByVehicleLicense -> find by vehicle
func FindByVehicleLicense(license string, offset int64, limit int64) (*[]Trace, error) {
	collection := mgm.Coll(&Trace{})

	traces := []Trace{}

	findOptions := options.Find()
	findOptions.SetLimit(limit)
	findOptions.SetSkip(offset)
	findOptions.SetProjection(bson.D{{"license", 0}})

	err := collection.SimpleFind(&traces, bson.M{"license": license}, findOptions)
	if err != nil {
		return nil, err
	}

	return &traces, nil
}
