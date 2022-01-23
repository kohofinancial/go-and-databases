package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"expvar"
	_ "expvar"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kohofinancial/go-and-databases/services"
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
	var setLimits bool
	flag.StringVar(&dsn, "dsn", "", "PostgreSQL DSN")
	flag.StringVar(&port, "port", "8080", "Service Port")
	flag.BoolVar(&setLimits, "limits", false, "Sets DB limits")
	flag.Parse()

	db, err := openDB(dsn, setLimits)
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

func openDB(dsn string, setLimits bool) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if setLimits {
		fmt.Println("setting limits")
		db.SetMaxOpenConns(5)
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
	w.WriteHeader(http.StatusBadRequest)
	fmt.Println(err)
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
