package controllers

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/lexgalante/go.iris/models/vehicles"
	"go.mongodb.org/mongo-driver/mongo"
)

//GetVehicle -> GET /v1/vehicles
func GetVehicle(ctx iris.Context) {
	currentUser := GetCurrentUserFromContext(ctx)

	vehicles, err := vehicles.FindByUserID(currentUser.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			NoContent(ctx)
			return
		}
		InternalServerError(ctx, ErrorUnexpectedError, err)
		return
	}

	Ok(ctx, vehicles)
}

//GetVehicleByID -> GET /v1/vehicles/{id}
func GetVehicleByID(ctx iris.Context) {
	currentUser := GetCurrentUserFromContext(ctx)

	id := ctx.Params().Get("id")
	if id == "" {
		BadRequest(ctx, MakeValidationError(ErrorInvalidParameter, "the parameter [id] cannot be null or empty"))
		return
	}

	if currentUser.InRole("admin") {
		vehicle, err := vehicles.FindByID(id)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				NotFound(ctx, "vehicle not found")
				return
			}
			InternalServerError(ctx, ErrorUnexpectedError, err)
			return
		}

		Ok(ctx, vehicle)
	}

	vehicle, err := vehicles.FindByOwner(id, currentUser.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			NotFound(ctx, "vehicle not found")
			return
		}
		InternalServerError(ctx, ErrorUnexpectedError, err)
		return
	}

	Ok(ctx, vehicle)
}

//PostVehicle -> POST /v1/vehicles
func PostVehicle(ctx iris.Context) {
	currentUser := GetCurrentUserFromContext(ctx)

	var vehicle vehicles.Vehicle
	if err := ctx.ReadBody(&vehicle); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			BadRequest(ctx, MakeValidationErrors(errs))
			return
		}
		BadRequest(ctx, MakeValidationError(ErrorInvalidPayload, "body cannot be desserialize"))
		return
	}

	vehicleExists, _ := vehicles.FindByLicense(vehicle.License)
	if vehicleExists != nil {
		Conflict(ctx, fmt.Sprintf("%s already exists by another user", vehicle.License))
		return
	}

	if !currentUser.InRole("admin") || vehicle.UserID == "" {
		vehicle.UserID = currentUser.Email
	}

	vehicle.Active = true

	id, err := vehicle.Save()
	if err != nil {
		InternalServerError(ctx, ErrorUnexpectedError, err)
		return
	}

	log.Println("new vehicle", id, "created by:", currentUser.Email)

	Created(ctx, vehicle)
}

//PutVehicle -> PUT /v1/vehicles/{id}
func PutVehicle(ctx iris.Context) {
	currentUser := GetCurrentUserFromContext(ctx)

	id := ctx.Params().Get("id")
	if id == "" {
		BadRequest(ctx, MakeValidationError(ErrorInvalidParameter, "the parameter [id] cannot be null or empty"))
		return
	}

	var vehicle vehicles.Vehicle
	if err := ctx.ReadBody(&vehicle); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			BadRequest(ctx, MakeValidationErrors(errs))
			return
		}
		BadRequest(ctx, MakeValidationError(ErrorInvalidPayload, "body cannot be desserialize"))
		return
	}

	vehicleInDb, err := vehicles.FindByID(id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			NotFound(ctx, "vehicle not found")
			return
		}
		InternalServerError(ctx, ErrorUnexpectedError, err)
		return
	}

	vehicleInDb.Name = vehicle.Name
	vehicleInDb.Model = vehicle.Model
	vehicleInDb.YearModel = vehicle.YearModel
	vehicleInDb.YearManufactory = vehicle.YearManufactory
	vehicleInDb.Color = vehicle.Color
	vehicleInDb.Active = vehicle.Active

	_, err = vehicleInDb.Save()
	if err != nil {
		InternalServerError(ctx, ErrorUnexpectedError, err)
		return
	}

	log.Println("user:", currentUser.Email, "updated vehicle:", id)

	Accepted(ctx, vehicleInDb)
}

//PatchVehicle -> PATCH /v1/vehicles/{id}/active | /v1/vehicles/{id}/active
func PatchVehicle(ctx iris.Context) {
	currentUser := GetCurrentUserFromContext(ctx)

	id := ctx.Params().Get("id")
	if id == "" {
		BadRequest(ctx, MakeValidationError(ErrorInvalidParameter, "the parameter [id] cannot be null or empty"))
		return
	}

	vehicle, err := vehicles.FindByID(id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			NotFound(ctx, "vehicle not found")
			return
		}
		InternalServerError(ctx, ErrorUnexpectedError, err)
		return
	}

	if currentUser.InRole("admin") && vehicle.UserID != currentUser.Email {
		Forbidden(ctx)
		return
	}

	if strings.Contains(ctx.Request().RequestURI, "inactive") {
		vehicle.Active = false
	} else {
		vehicle.Active = true
	}

	_, err = vehicle.Save()
	if err != nil {
		InternalServerError(ctx, ErrorUnexpectedError, err)
		return
	}

	log.Println("user:", currentUser.Email, "change user active to:", vehicle.Active)

	Accepted(ctx, vehicle)
}

//DeleteVehicle -> DELETE /v1/vehicles/{id}
func DeleteVehicle(ctx iris.Context) {
	currentUser := GetCurrentUserFromContext(ctx)

	id := ctx.Params().Get("id")
	if id == "" {
		BadRequest(ctx, MakeValidationError(ErrorInvalidParameter, "the parameter [id] cannot be null or empty"))
		return
	}

	vehicle, err := vehicles.FindByID(id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			NotFound(ctx, "vehicle not found")
			return
		}
		InternalServerError(ctx, ErrorUnexpectedError, err)
		return
	}

	if currentUser.InRole("admin") && vehicle.UserID != currentUser.Email {
		Forbidden(ctx)
		return
	}

	err = vehicle.Remove()
	if err != nil {
		InternalServerError(ctx, ErrorUnexpectedError, err)
		return
	}

	log.Println("user:", currentUser.Email, "delete vehicle:", id)

	NoContent(ctx)
}
