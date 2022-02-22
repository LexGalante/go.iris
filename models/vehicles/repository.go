package vehicles

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Save -> create or update model
func (v *Vehicle) Save() (primitive.ObjectID, error) {
	collection := mgm.Coll(v)

	if v.ID == primitive.NilObjectID {
		err := collection.Create(v)
		if err != nil {
			return v.ID, err
		}
	}

	err := collection.Update(v)
	if err != nil {
		return v.ID, err
	}

	return v.ID, nil
}

//Remove -> delete model
func (v *Vehicle) Remove() error {
	collection := mgm.Coll(v)

	return collection.Delete(v)
}

//FindByID -> find by id
func FindByID(id string) (*Vehicle, error) {
	vehicle := &Vehicle{}

	collection := mgm.Coll(vehicle)

	err := collection.FindByID(id, vehicle)
	if err != nil {
		return nil, err
	}

	return vehicle, nil
}

//FindByOwner -> find by owner of the vehicle
func FindByOwner(id string, email string) (*Vehicle, error) {
	vehicle := &Vehicle{}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	collection := mgm.Coll(vehicle)

	err = collection.SimpleFind(vehicle, bson.M{"_id": objectID, "user_id": email})
	if err != nil {
		return nil, err
	}

	return vehicle, nil
}

//FindByUserID -> find by id
func FindByUserID(email string) (*[]Vehicle, error) {
	vehicles := []Vehicle{}

	collection := mgm.Coll(&Vehicle{})

	err := collection.SimpleFind(&vehicles, bson.M{"user_id": email})
	if err != nil {
		return nil, err
	}

	return &vehicles, nil
}

//FindByLicense -> find by id
func FindByLicense(license string) (*Vehicle, error) {
	vehicle := &Vehicle{}

	collection := mgm.Coll(vehicle)

	err := collection.First(bson.M{"license": license}, vehicle)
	if err != nil {
		return nil, err
	}

	return vehicle, nil
}
