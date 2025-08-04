package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/benj3coll/bookings/pkg/config"
	"github.com/benj3coll/bookings/pkg/handlers"
	"github.com/benj3coll/bookings/pkg/render"
)

const portNr = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {

	// change this to true when in production
	app.InProduction = false

	createSession()

	app.UseCache = true
	render.InitCache(&app)

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	fmt.Printf("Starting application on port %s\n", portNr)

	server := &http.Server{
		Addr:    portNr,
		Handler: routes(&app),
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

func createSession() {
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session
}
