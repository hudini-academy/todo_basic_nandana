package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
)

// The serverError helper writes an error message and stack trace to the errorLog
// then sends a generic 500 Internal Server Error response to the user
func (a *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s", err.Error())
	a.errorLog.Println(trace)
	http.Error(w, trace, http.StatusInternalServerError)
}

func infoError() (*log.Logger, *log.Logger) {
	f, err := os.OpenFile("./info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666) //Logging the page info into a file
	if err != nil {
		log.Fatal(err)
	}
	// defer f.Close() //file closed

	f1, err := os.OpenFile("./error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666) //Logging the error into a file f1
	if err != nil {
		log.Fatal(err)
	}
	// defer f1.Close()

	infoLog := log.New(f, "INFO\t", log.Ldate|log.Ltime|log.LUTC)          //storing the formatted info message to file f
	errorLog := log.New(f1, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile) //storing the formatted error message to file f

	return errorLog, infoLog
}

// The clientError helper sends status code and corresponding desc to user
// send responses like 400 "Bad Request" when there's a problem with the request that the user sent
func (a *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// also implement a notFound helper
// it is convenience wrapper around clientError which sends a 404 Not Found response to the user.
func (a *application) notFound(w http.ResponseWriter) {
	a.clientError(w, http.StatusNotFound)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
