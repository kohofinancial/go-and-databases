package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"expvar"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kohofinancial/go-and-databases/services"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

type app struct {
	UserService services.UserService
}

func main() {
	var dsn string
	var port string
	var connLimits bool
	var idleLimits bool
	flag.StringVar(&dsn, "dsn", "", "PostgreSQL DSN")
	flag.StringVar(&port, "port", "8080", "Service Port")
	flag.BoolVar(&connLimits, "conn-limits", false, "Sets DB limits")
	flag.BoolVar(&idleLimits, "idle-limits", false, "Sets DB limits")
	flag.Parse()

	db, err := openDB(dsn, connLimits, idleLimits)
	if err != nil {
		log.Fatalln(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(db)
	expvar.Publish("db-stats", expvar.Func(func() interface{} {
		return db.Stats()
	}))

	application := app{UserService: services.NewPostgresUserService(db)}

	r := mux.NewRouter()
	userRouter := r.PathPrefix("/users").Subrouter()
	userRouter.HandleFunc("/{id}", application.GetUser).Methods(http.MethodGet)
	userRouter.HandleFunc("/all", application.DeleteAll).Methods(http.MethodDelete)
	userRouter.HandleFunc("/{id}", application.DeleteUser).Methods(http.MethodDelete)
	userRouter.HandleFunc("/{id}", application.UpdateUser).Methods(http.MethodPut)
	userRouter.HandleFunc("", application.AddUser).Methods(http.MethodPost)

	r.Handle("/debug/vars", expvar.Handler())

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func openDB(dsn string, connLimits bool, idleLimits bool) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if connLimits {
		db.SetMaxOpenConns(5)
	}

	if idleLimits {
		db.SetMaxIdleConns(5)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func ReturnError(w http.ResponseWriter, err error) {
	//Ideally the errors are handled lower than this but for a quick project this gave me what I needed.
	if errors.Is(err, context.DeadlineExceeded) {
		w.WriteHeader(http.StatusRequestTimeout)
	} else if errors.Is(err, pq.Error{}) {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	response := ErrorResponse{Error: err.Error()}
	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = w.Write(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

type ErrorResponse struct {
	Error string `json:"error,omitempty"`
}
