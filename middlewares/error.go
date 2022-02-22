package middlewares

import (
	"log"

	"github.com/kataras/iris/v12"
)

//LogErrorMiddleware -> middleware for log unhandle unexpected errors
func LogErrorMiddleware() iris.Handler {
	return func(ctx iris.Context) {
		ctx.Next()

		statusCode := ctx.GetStatusCode()

		if statusCode >= 300 {
			err := ctx.GetErr()
			if err != nil {
				log.Println("[LogErrorMiddleware] path:", ctx.Path(), "error: ", err.Error())
			}
		}
	}
}
