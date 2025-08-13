package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/benj3coll/bookings/internal/config"
	"github.com/benj3coll/bookings/internal/handlers"
	"github.com/benj3coll/bookings/internal/helpers"
	"github.com/benj3coll/bookings/internal/models"
	"github.com/benj3coll/bookings/internal/render"
)

const portNr = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {

	err := run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Starting application on port %s\n", portNr)

	server := &http.Server{
		Addr:    portNr,
		Handler: routes(&app),
	}
	err = server.ListenAndServe()
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

func run() error {
	// What am I going to put in the session
	gob.Register(models.Reservation{})

	// change this to true when in production
	app.InProduction = false

	app.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	createSession()

	app.UseCache = false
	app.TemplatePath = "./templates"
	err := render.InitCache(&app)
	if err != nil {
		return err
	}

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	helpers.NewHelpers(&app)

	return nil
}
