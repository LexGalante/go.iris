package controllers

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/kataras/iris/v12"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/lexgalante/go.iris/models/users"
	"github.com/lexgalante/go.iris/models/vehicles"
)

//GetCurrentUserFromContext -> retrive user data from context
func GetCurrentUserFromContext(ctx iris.Context) *users.UserClaims {
	currentClaims, err := ctx.User().GetRaw()
	if err != nil {
		panic("cannot recover user data from context")
	}

	usersClaims, _ := currentClaims.(*users.UserClaims)

	return usersClaims
}

//GetUser -> GET /v1/users
func GetUser(ctx iris.Context) {
	currentUser := GetCurrentUserFromContext(ctx)

	if !currentUser.InRole("admin") {
		Forbidden(ctx)
		return
	}

	if ctx.URLParamExists("email") {
		email := ctx.URLParam("email")

		user, err := users.FindByEmail(email)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				NotFound(ctx, fmt.Sprintf("user %s not found", email))
				return
			}
			InternalServerError(ctx, ErrorUnexpectedError, err)
			return
		}

		Ok(ctx, user)
	} else if ctx.Params().Exists("id") {
		id := ctx.Params().Get("id")

		user, err := users.FindByID(id)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				NotFound(ctx, fmt.Sprintf("user %s not found", id))
				return
			}
			InternalServerError(ctx, ErrorUnexpectedError, err)
			return
		}

		Ok(ctx, user)
	} else {
		BadRequest(ctx, MakeValidationError(ErrorInvalidParameter, "use path parameter [id] or query string parameter [email]"))
		return
	}
}

//GetUserVehicles -> GET /v1/users/{id}/vehicles
func GetUserVehicles(ctx iris.Context) {
	id := ctx.Params().Get("id")
	if id == "" {
		BadRequest(ctx, MakeValidationError(ErrorInvalidParameter, "the parameter [id] cannot be null or empty"))
		return
	}

	currentUser := GetCurrentUserFromContext(ctx)

	if !currentUser.InRole("admin") && id != currentUser.Id {
		Forbidden(ctx)
		return
	}

	vehicles, err := vehicles.FindByUserID(currentUser.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			NoContent(ctx)
			return
		}
	}

	Ok(ctx, vehicles)
}

//PatchUser -> GET /v1/users/{id}/active | /v1/users/{id}/inactive
func PatchUser(ctx iris.Context) {
	currentUser := GetCurrentUserFromContext(ctx)

	if !currentUser.InRole("admin") {
		Forbidden(ctx)
		return
	}

	id := ctx.Params().Get("id")
	if id == "" {
		BadRequest(ctx, MakeValidationError(ErrorInvalidParameter, "the parameter [id] cannot be null or empty"))
		return
	}

	user, err := users.FindByID(id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			NotFound(ctx, fmt.Sprintf("user %s not found", id))
			return
		}
	}

	if strings.Contains(ctx.Request().RequestURI, "inactive") {
		user.Active = false
	} else {
		user.Active = true
	}

	_, err = user.Save()
	if err != nil {
		InternalServerError(ctx, ErrorUnexpectedError, err)
		return
	}

	user.Password = ""

	log.Println("admin:", currentUser.Email, " changes user:", user.Email, " active:", user.Active)

	Accepted(ctx, user)
}
