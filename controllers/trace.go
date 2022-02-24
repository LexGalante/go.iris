package controllers

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/lexgalante/go.iris/models/traces"
	"github.com/lexgalante/go.iris/models/vehicles"
	"go.mongodb.org/mongo-driver/mongo"
)

//PaginateTraces -> GET /vehicles/{id}/traces?page=1&page_size=50
func PaginateTraces(ctx iris.Context) {
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

	//check user is admin or owner of vehicle
	if !currentUser.InRole("admin") && vehicle.UserID != currentUser.Email {
		Forbidden(ctx)
		return
	}

	page := ctx.URLParamInt64Default("page", 1)
	pageSize := ctx.URLParamInt64Default("page_size", 10)

	traces, err := traces.FindByVehicleLicense(vehicle.License, (page-1)*pageSize, pageSize)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			NoContent(ctx)
			return
		}
		InternalServerError(ctx, ErrorUnexpectedError, err)
		return
	}

	Ok(ctx, traces)
}

//PostTrace -> /vehicles/{id}/traces
func PostTrace(ctx iris.Context) {
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

	//check user is admin or owner of vehicle
	if !currentUser.InRole("admin") && vehicle.UserID != currentUser.Email {
		Forbidden(ctx)
		return
	}

	var trace traces.Trace
	if trace.License == "" {
		trace.License = vehicle.License
	}
	if trace.DateTime.IsZero() {
		trace.DateTime = time.Now()
	}

	if err := ctx.ReadBody(&trace); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			BadRequest(ctx, MakeValidationErrors(errs))
			return
		}
		BadRequest(ctx, MakeValidationError(ErrorInvalidPayload, "body cannot be desserialize"))
		return
	}

	idTrace, err := trace.Save()
	if err != nil {
		InternalServerError(ctx, ErrorUnexpectedError, err)
		return
	}

	log.Println("new trace", idTrace, " by vehicle:", vehicle.License, "created by:", currentUser.Email)

	Created(ctx, trace)
}

//PostTraceFromCSV -> /vehicles/{id}/traces/upload
func PostTraceFromCSV(ctx iris.Context) {
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

	//check user is admin or owner of vehicle
	if !currentUser.InRole("admin") && vehicle.UserID != currentUser.Email {
		Forbidden(ctx)
		return
	}

	_, fileHeader, err := ctx.FormFile("file")
	if err != nil {
		InternalServerError(ctx, ErrorUnexpectedError, err)
		return
	}

	now := time.Now()
	dir := fmt.Sprintf("./uploads/%d-%02d-%02dT%02d:%02d:%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	dest := filepath.Join(dir, vehicle.License)
	ctx.SaveFormFile(fileHeader, dest)

	NoContent(ctx)
}
