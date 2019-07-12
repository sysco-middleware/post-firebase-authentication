package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/labstack/echo"
	"google.golang.org/api/option"
)

var client *auth.Client

func main() {
	var err error
	var app *firebase.App
	e := echo.New()
	e.HideBanner = true

	// initialize firebase app
	opt := option.WithCredentialsFile("fb.json")
	if app, err = firebase.NewApp(context.Background(), nil, opt); err != nil {
		e.Logger.Fatalf("error initializing app: %v", err)
	}

	// create a client to communicate with firebase project
	if client, err = app.Auth(context.Background()); err != nil {
		e.Logger.Fatalf("error creating firebase client: %v", err)
	}

	e.Static("/", "public")

	grp := e.Group("/api/")
	grp.GET("serviceA", serviceA)
	grp.GET("serviceB", serviceB)

	e.Logger.Fatal(e.StartServer(server()))
}

func serviceA(ctx echo.Context) error {
	idToken := ctx.Request().Header.Get("id-token")
	if idToken == "" {
		return ctx.JSON(200, map[string]string{
			"error":   "no token",
			"details": "id-token missing in header",
		})
	}

	// use the client to verify the id token sent by client
	t, err := client.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		return ctx.JSON(200, map[string]string{
			"error":   "invalid token",
			"message": " invalid id-token provided",
		})
	}

	fmt.Println("id-token verified")

	resp := make(map[string]string)
	resp["name"] = "Service A"
	resp["time"] = time.Now().Format(time.UnixDate)
	resp["uid"] = t.UID
	resp["project-url"] = t.Issuer
	resp["project-id"] = t.Audience

	return ctx.JSON(200, resp)
}
func serviceB(ctx echo.Context) error {
	idToken := ctx.Request().Header.Get("id-token")
	if idToken == "" {
		return ctx.JSON(200, map[string]string{
			"error": "no token",
		})
	}
	t, err := client.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		return ctx.JSON(200, map[string]string{
			"error":   "invalid token",
			"message": " invalid id-token provided",
		})
	}
	fmt.Println("id-token verified")

	resp := make(map[string]string)
	resp["name"] = "Service B"
	resp["time"] = time.Now().Format(time.UnixDate)
	resp["uid"] = t.UID
	resp["project-url"] = t.Issuer
	resp["project-id"] = t.Audience

	return ctx.JSON(200, resp)
}

func server() *http.Server {
	return &http.Server{
		Addr:         ":1313",          // port
		ReadTimeout:  20 * time.Minute, // read timeout
		WriteTimeout: 20 * time.Minute, // write timeout
	}
}
