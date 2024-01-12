package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"web3/pkg/helpers"

	"github.com/justinas/nosurf"
)

func LogRequestInfo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		fmt.Printf("%d/%d/%d : %d:%d ", now.Month(), now.Day(), now.Year(), now.Hour(), now.Minute())
		fmt.Println(r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func SetupSession(next http.Handler) http.Handler {
	return sessionManager.LoadAndSave(next)
}

func NoSurf(next http.Handler) http.Handler {
	noSurfHandler := nosurf.New(next)
	noSurfHandler.SetBaseCookie(http.Cookie{
		Name:     "mycsrfcookie",
		Path:     "/",
		Domain:   "",
		Secure:   false,
		HttpOnly: true,
		MaxAge:   3600,
		SameSite: http.SameSiteLaxMode,
	})
	return noSurfHandler
}

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsAuthenticated(r) {
			app.Session.Put(r.Context(), "error", "You aren't logged in")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			log.Fatal("Error logging in")
			return
		}
		next.ServeHTTP(w, r)
	})
}
