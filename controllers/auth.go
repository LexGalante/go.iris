package controllers

import (
	"errors"
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/lexgalante/go.iris/models/users"
	"go.mongodb.org/mongo-driver/mongo"
)

//Register -> POST /auth/register
func Register(ctx iris.Context) {
	var user users.User
	if err := ctx.ReadBody(&user); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			BadRequest(ctx, MakeValidationErrors(errs))
			return
		}
		BadRequest(ctx, MakeValidationError(ErrorInvalidPayload, "body cannot be desserialize"))
		return
	}

	userExists, _ := users.FindByEmail(user.Email)
	if userExists != nil {
		Conflict(ctx, fmt.Sprintf("%s already exists", user.Email))
		return
	}

	if err := user.CryptPassword(); err != nil {
		InternalServerError(ctx, ErrorUnexpectedError, err)
		return
	}

	user.Roles = []string{"guest"}
	user.Active = true

	id, err := user.Save()
	if err != nil {
		InternalServerError(ctx, ErrorUnexpectedError, err)
		return
	}

	user.Password = ""

	log.Println("new user registered:", id)

	Created(ctx, user)
}

//Login -> POST /auth/login
func Login(ctx iris.Context) {
	var user users.User
	if err := ctx.ReadBody(&user); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			Unauthorized(ctx, MakeValidationErrors(errs))
			return
		}
		BadRequest(ctx, MakeValidationError(ErrorInvalidPayload, "body cannot be desserialize"))
		return
	}

	userInDb, err := users.FindByEmail(user.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			Unauthorized(ctx, MakeValidationError(ErrorInvalidCredentials, "invalid user or password"))
			return
		}
	}

	if !userInDb.VerifyPassword(user.Password) {
		Unauthorized(ctx, MakeValidationError(ErrorInvalidCredentials, "invalid user or password"))
		return
	}

	accessToken, err := userInDb.CreateAccessToken()
	if err != nil {
		InternalServerError(ctx, ErrorUnexpectedError, err)
		return
	}

	Ok(ctx, accessToken)
}
