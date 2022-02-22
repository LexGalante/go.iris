package middlewares

import (
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/lexgalante/go.iris/controllers"
	"github.com/lexgalante/go.iris/models/users"
)

//JwtMiddleware -> middleware for jwk validation
func JwtMiddleware() iris.Handler {
	return func(ctx iris.Context) {

		authorizationHeader := ctx.GetHeader("Authorization")
		accessToken := strings.ReplaceAll(authorizationHeader, "Bearer ", "")

		validateToken, err := users.ValidateAcessToken(accessToken)
		if err != nil {
			controllers.Unauthorized(ctx, controllers.MakeValidationError(controllers.ErrorInvalidCredentials, "access token invalid"))
			return
		}

		if claims, ok := validateToken.Claims.(*users.UserClaims); ok && validateToken.Valid {
			//write user data in current context
			ctx.SetUser(claims)

			ctx.Next()
		} else {
			controllers.Unauthorized(ctx, controllers.MakeValidationError(controllers.ErrorInvalidCredentials, "access token invalid"))
		}
	}
}
