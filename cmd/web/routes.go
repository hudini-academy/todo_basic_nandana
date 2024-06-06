package main

import (
	"net/http"

	"github.com/bmizerany/pat"
)

func (a *application) routes() http.Handler { //http.ServeMux changed to http.handler because of middleware
	// a mux handler is created
	//mux := http.NewServeMux()
	mux := pat.New()
	mux.Get("/", a.session.Enable(http.HandlerFunc(a.App)))
	mux.Post("/Add", a.session.Enable(http.HandlerFunc(a.Add)))
	mux.Post("/delete", a.session.Enable(http.HandlerFunc(a.Delete)))
	mux.Post("/Update", a.session.Enable(http.HandlerFunc(a.Update)))
	mux.Get("/user/signup", a.session.Enable(http.HandlerFunc(a.signupUserForm)))
	mux.Post("/user/signup", a.session.Enable(http.HandlerFunc(a.signupUser)))
	mux.Get("/user/login", a.session.Enable(http.HandlerFunc(a.loginUserForm)))
	mux.Post("/user/login", a.session.Enable(http.HandlerFunc(a.loginUser)))
	mux.Post("/user/logout", a.session.Enable(http.HandlerFunc(a.logoutUser)))
	// mux.HandleFunc("/", a.App) //Three functions passed as handler function with mux
	// mux.HandleFunc("/Add", a.Add)
	// mux.HandleFunc("/delete", a.Delete)
	// mux.HandleFunc("/Update", a.Update)

	fileServer := http.FileServer(http.Dir("./ui/static"))       //serves files out of the "./ui/static" directory
	mux.Get("/static/", http.StripPrefix("/static", fileServer)) //strip the prefix /static from the url and passes to fileserver
	return a.recoverPanic(a.logRequest((secureHeaders(mux))))
}
