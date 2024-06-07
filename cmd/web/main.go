package main

import (
	"Todo_Application/pkg/models/mysql"
	"flag"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	Todo     *mysql.TodoModel
	session  *sessions.Session
	users    *mysql.UserModel
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "root:root@/Todo?parseTime=true", "MySQL database")
	flag.Parse()
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret Key")
	flag.Parse()
	errorLog, infoLog := infoError()
	//Creating a database connection pool
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()
	session := sessions.New([]byte(*secret)) //Initializing a new sessio manager
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	a := &application{ //an interface of structure application storing error and info
		errorLog: errorLog,
		infoLog:  infoLog,
		Todo:     &mysql.TodoModel{DB: db},
		session:  session,
		users:    &mysql.UserModel{DB: db},
	}
	// tlsConfig := &tls.Config{
	// 	PreferServerCipherSuites: true,
	// 	CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	// }
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  a.routes(),
		//TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *addr)

	errr := srv.ListenAndServe() //Listen to port on addr and pass the mux handler and display
	errorLog.Fatal(errr)

}
