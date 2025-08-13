package render

import (
	"bytes"
	"errors"
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

func InitCache(a *config.AppConfig) error {
	app = a
	// create a template cache
	tcache, err := createTemplateCache(app.TemplatePath)
	if err != nil {
		return err
	}
	templateCache = tcache
	return nil
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	return td
}

// RenderTemplate renders a template
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	if !app.UseCache {
		templateCache, _ = createTemplateCache(app.TemplatePath)
	}
	// get requested template from cache
	t, ok := templateCache[tmpl]
	if !ok {
		log.Println("Could not get template from template cache")
		return errors.New("could not get template from template cache")
	}

	td = AddDefaultData(td, r)

	buf := new(bytes.Buffer)
	err := t.Execute(buf, td)
	if err != nil {
		log.Println(err)
		return err
	}

	// render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println("error writing template to browser", err)
		return err
	}
	return nil
}

func createTemplateCache(templatePath string) (map[string]*template.Template, error) {
	tcache := map[string]*template.Template{}

	// get all of the files named *.page.html from templatePath
	pages, err := filepath.Glob(templatePath + "/*.page.html")
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
		matches, err := filepath.Glob(templatePath + "/*.layout.html")
		if err != nil {
			return tcache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(templatePath + "/*.layout.html")
			if err != nil {
				return tcache, err
			}
		}

		tcache[name] = ts
	}

	return tcache, nil
}
