package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/indexcoder/bookings/internal/config"
	"github.com/indexcoder/bookings/internal/handlers"
	"github.com/indexcoder/bookings/internal/render"
	"log"
	"net/http"
	"time"
)

const portNumber = "9090"

var app config.AppConfig
var session *scs.SessionManager

func main() {

	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.TemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandler(repo)
	render.NewTemplate(&app)
	fmt.Println("Starting application on port: ", portNumber)
	srv := &http.Server{
		Addr:    ":" + portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
