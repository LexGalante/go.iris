package main

import (
	"log"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/kamva/mgm/v3"
	"github.com/kataras/iris/v12"
	"github.com/lexgalante/go.iris/controllers"
	"github.com/lexgalante/go.iris/middlewares"
	"github.com/lexgalante/go.iris/utils"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("unable to load .env file")
	}

	err = mgm.SetDefaultConfig(&mgm.Config{CtxTimeout: time.Second * 30}, os.Getenv("DB_NAME"), options.Client().ApplyURI(os.Getenv("DB_URI")))

	validate := validator.New()
	validate.RegisterValidation("licenseplate", utils.ValidatorLicensePlate)

	app := iris.Default()
	app.Validator = validate

	authAPI := app.Party("/auth")
	{
		authAPI.Post("/register", controllers.Register)
		authAPI.Post("/login", controllers.Login)
	}

	usersAPI := app.Party("/v1/users")
	{
		usersAPI.Use(middlewares.JwtMiddleware())
		usersAPI.Use(middlewares.LogErrorMiddleware())
		usersAPI.Use(iris.Compression)
		usersAPI.Get("/", controllers.GetUser)
		usersAPI.Get("/{id:string}", controllers.GetUser)
		usersAPI.Get("/{id:string}/vehicles", controllers.GetUserVehicles)
		usersAPI.Patch("/users/{id:string}/active", controllers.PatchUser)
		usersAPI.Patch("/users/{id:string}/inactive", controllers.PatchUser)
	}

	vehiclesAPI := app.Party("/v1/vehicles")
	{
		vehiclesAPI.Use(middlewares.JwtMiddleware())
		vehiclesAPI.Use(middlewares.LogErrorMiddleware())
		vehiclesAPI.Use(iris.Compression)
		vehiclesAPI.Get("/", controllers.GetVehicle)
		vehiclesAPI.Get("/{id:string}", controllers.GetVehicleByID)
		vehiclesAPI.Post("/", controllers.PostVehicle)
		vehiclesAPI.Put("/{id:string}", controllers.PutVehicle)
		vehiclesAPI.Delete("/{id:string}", controllers.DeleteVehicle)
		vehiclesAPI.Patch("/{id:string}/active", controllers.PatchVehicle)
		vehiclesAPI.Patch("/{id:string}/inactive", controllers.PatchVehicle)
	}

	log.Fatal(app.Listen(os.Getenv("APP_PORT")))
}
