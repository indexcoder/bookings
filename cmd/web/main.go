package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/indexcoder/bookings/internal/config"
	"github.com/indexcoder/bookings/internal/driver"
	"github.com/indexcoder/bookings/internal/handlers"
	"github.com/indexcoder/bookings/internal/helpers"
	"github.com/indexcoder/bookings/internal/models"
	"github.com/indexcoder/bookings/internal/render"
	"log"
	"net/http"
	"os"
	"time"
)

const portNumber = "9090"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {

	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	defer close(app.MailChan)
	fmt.Println("Starting mail listener.")
	listenForMail()

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

func run() (*driver.DB, error) {

	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(models.Reservation{})
	gob.Register(models.RoomRestriction{})
	gob.Register(map[string]int{})

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connect to database
	log.Println("connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=12345 sslmode=disable")
	if err != nil {
		log.Fatal("Cannot connect to database! dying...", err.Error())
	}
	log.Println("connected to database")

	tc, err := render.TemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandler(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
