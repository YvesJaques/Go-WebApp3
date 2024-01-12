package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"time"
	"web3/models"
	"web3/pkg/config"
	"web3/pkg/dbdriver"
	handlers "web3/pkg/handlers"

	"github.com/alexedwards/scs/v2"
)

var sessionManager *scs.SessionManager
var app config.AppConfig

func main() {

	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() (*dbdriver.DB, error) {
	gob.Register(models.Article{})

	// 26 Add table models in session
	gob.Register(models.User{})
	gob.Register(models.Post{})

	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.Secure = false
	// for testing
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	app.Session = sessionManager

	db, err := dbdriver.ConnectSQL("host=localhost port=5432 dbname=blog_db user=postgres password=postgres")
	if err != nil {
		log.Fatal("Can't connect to database")
	}

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	return db, nil
}
