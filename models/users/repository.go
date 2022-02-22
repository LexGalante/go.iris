package users

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Save -> create or update model
func (u *User) Save() (primitive.ObjectID, error) {
	collection := mgm.Coll(u)

	if u.ID == primitive.NilObjectID {
		err := collection.Create(u)
		if err != nil {
			return u.ID, err
		}
	}

	err := collection.Update(u)
	if err != nil {
		return u.ID, err
	}

	return u.ID, nil
}

//Remove -> delete model
func (u *User) Remove() error {
	collection := mgm.Coll(u)

	return collection.Delete(u)
}

//FindByID -> find by id
func FindByID(id string) (*User, error) {
	user := &User{}

	collection := mgm.Coll(user)

	err := collection.FindByID(id, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

//FindByEmail -> find by email
func FindByEmail(email string) (*User, error) {
	user := &User{}

	collection := mgm.Coll(user)

	err := collection.First(bson.M{"email": email}, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
