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

//ADD ONS :

//To get an environment variable use GetEnv. Also use SetEnv to set a new environment
// addr := os.Getenv("SNIPPETBOX_ADDR")
// os.Setenv("FOO", "1")
// log.Println("Foo:", os.Getenv("FOO"))

// The below code will display the whole environment variables of the OS
//for _, e := range os.Environ() {
// 	pair := strings.SplitN(e, "=", 2)
// 	fmt.Println(pair[0])
// }

//METHOD1
//changed from "err := http.ListenAndServe(*addr, mux)" to :(for passing the serverstruct)
//changed from "log.Fatal(err)" to : (passed the errorlog of serverstruct as err)

//METHOD 2
//logging the information and error of current file. 3 parameters are writer,prefix string and the data passed,"|" passed as appending all items together
//infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
//errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

//:we create a *server struct* to create a new error logger instead of inbuilt logger
// srv := &http.Server{
// 	Addr:     *addr,
// 	ErrorLog: errorLog,
// 	Handler:  mux,
// }
//changed from "log.Printf("Starting server on :%s", *addr)" to:
//infoLog.Printf("Starting server on %s", *addr)

//err := srv.ListenAndServe()
//errorLog.Fatal(err)

//sql.Open("mysql", "web:pass@/Todo?parseTime=true")	//mysql is driver name,2nd parameter describes how to connect to a db
