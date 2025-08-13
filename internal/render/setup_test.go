package render

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/benj3coll/bookings/internal/config"
	"github.com/benj3coll/bookings/internal/models"
)

var testApp config.AppConfig
var session *scs.SessionManager

func TestMain(m *testing.M) {
	// What am I going to put in the session
	gob.Register(models.Reservation{})

	// change this to true when in production
	testApp.InProduction = false
	testApp.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	testApp.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	createSession()
	testApp.UseCache = false
	testApp.TemplatePath = "./../../templates"

	app = &testApp

	os.Exit(m.Run())
}

func createSession() {
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = testApp.InProduction

	testApp.Session = session
}

type myWriter struct{}

func (tw *myWriter) Header() http.Header {
	return http.Header{}
}

func (tw *myWriter) WriteHeader(i int) {

}

func (tw *myWriter) Write(b []byte) (int, error) {
	return len(b), nil
}
