package main

import (
	"testing"

	"github.com/benj3coll/bookings/internal/config"
	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch mux.(type) {
	case *chi.Mux:
		// do nothing
	default:
		t.Error("type is not *chi.Mux")
	}

}
