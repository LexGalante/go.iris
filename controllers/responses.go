package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
)

//Ok -> 200
func Ok(ctx iris.Context, response interface{}) {
	ctx.JSON(response)
}

//Created -> 201
func Created(ctx iris.Context, response interface{}) {
	ctx.StatusCode(iris.StatusCreated)
	ctx.JSON(response)
}

//Accepted -> 202
func Accepted(ctx iris.Context, response interface{}) {
	ctx.StatusCode(iris.StatusAccepted)
	ctx.JSON(response)
}

//NoContent -> 204
func NoContent(ctx iris.Context) {
	ctx.StatusCode(iris.StatusNoContent)
}

//BadRequest -> stop with 400
func BadRequest(ctx iris.Context, validationsErros []ValidationError) {
	ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
		Title(ErrorInvalidRequest).
		Detail("invalid request").
		Key("erros", validationsErros))
}

//Unauthorized -> stop with 401
func Unauthorized(ctx iris.Context, validationsErros []ValidationError) {
	ctx.StopWithProblem(iris.StatusUnauthorized, iris.NewProblem().
		Title(ErrorInvalidRequest).
		Detail("unathorized user").
		Key("erros", validationsErros))
}

//Forbidden -> stop with 403
func Forbidden(ctx iris.Context) {
	ctx.StopWithProblem(iris.StatusForbidden, iris.NewProblem().
		Title(ErrorInvalidRequest).
		Detail("forbidden action").
		Key("erros", []ValidationError{{ErrorInsufficientPermission, "you cannot perform this action"}}))
}

//NotFound -> stop with 404
func NotFound(ctx iris.Context, message string) {
	ctx.StopWithProblem(iris.StatusNotFound, iris.NewProblem().
		Title(ErrorInvalidRequest).
		Detail("resource not found").
		Key("erros", []ValidationError{{ErrorNotFound, message}}))
}

//Conflict -> stop with 409
func Conflict(ctx iris.Context, message string) {
	ctx.StopWithProblem(iris.StatusNotFound, iris.NewProblem().
		Title(ErrorInvalidRequest).
		Detail("resource not found").
		Key("erros", []ValidationError{{ErrorNotFound, message}}))
}

//UnprocessableEntity -> stop with 422
func UnprocessableEntity(ctx iris.Context, errs validator.ValidationErrors) {
	ctx.StopWithProblem(iris.StatusUnprocessableEntity, iris.NewProblem().
		Title(ErrorInvalidRequest).
		Detail("invalid entity request").
		Key("erros", MakeValidationErrors(errs)))
}

//InternalServerError -> stop with 500
func InternalServerError(ctx iris.Context, code string, err error) {
	ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
		Title(ErrorInvalidRequest).
		Detail("unexpected error ocurred").
		Key("erros", []ValidationError{{code, err.Error()}}))
}
