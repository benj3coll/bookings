package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/benj3coll/bookings/internal/config"
	"github.com/benj3coll/bookings/internal/models"
	"github.com/justinas/nosurf"
)

var app *config.AppConfig
var templateCache map[string]*template.Template

func InitCache(a *config.AppConfig) {
	app = a
	// create a template cache
	tcache, err := createTemplateCache()
	if err != nil {
		log.Fatal(err)
	}
	templateCache = tcache
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)
	return td
}

// RenderTemplate renders a template
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	if !app.UseCache {
		templateCache, _ = createTemplateCache()
	}
	// get requested template from cache
	t, ok := templateCache[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	td = AddDefaultData(td, r)

	buf := new(bytes.Buffer)
	err := t.Execute(buf, td)
	if err != nil {
		log.Println(err)
	}

	// render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

func createTemplateCache() (map[string]*template.Template, error) {
	log.Println("Create template cache")
	tcache := map[string]*template.Template{}

	// get all of the files named *.page.html from ./templates
	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		return tcache, err
	}

	// range through all files ending with *.page.html
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return tcache, err
		}

		// get all of the layout files
		matches, err := filepath.Glob("./templates/*.layout.html")
		if err != nil {
			return tcache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.html")
			if err != nil {
				return tcache, err
			}
		}

		tcache[name] = ts
	}

	return tcache, nil
}
