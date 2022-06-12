package main

import (
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"net/http"
	"techpark_db/internal/handler"
	"techpark_db/internal/infra/psql"
)

func main() {
	db, err := psql.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Info("Successful connect to database.")

	psqlStorage := psql.NewStorage(db)

	handler := handler.NewHandler(psqlStorage)

	router := mux.NewRouter()
	routerAPI := router.PathPrefix("/api").Subrouter()

	/*====================== FORUM ======================*/
	routerAPI.HandleFunc("/forum/create", handler.ForumCreate).Methods("POST")
	routerAPI.HandleFunc("/forum/{slug:[A-Za-z0-9._-]+}/details", handler.ForumDetails).Methods("GET")
	routerAPI.HandleFunc("/forum/{slug:[A-Za-z0-9._-]+}/create", handler.ForumCreateThread).Methods("POST")
	routerAPI.HandleFunc("/forum/{slug:[A-Za-z0-9._-]+}/users", handler.ForumUsers).Methods("GET")
	routerAPI.HandleFunc("/forum/{slug:[A-Za-z0-9._-]+}/threads", handler.ForumThreads).Methods("GET")

	/*====================== USER ======================*/
	routerAPI.HandleFunc("/user/{nickname:[A-Za-z0-9.]+}/create", handler.UserCreate).Methods("POST")
	routerAPI.HandleFunc("/user/{nickname:[A-Za-z0-9.]+}/profile", handler.UserDetails).Methods("GET")
	routerAPI.HandleFunc("/user/{nickname:[A-Za-z0-9.]+}/profile", handler.UserUpdate).Methods("POST")

	log.Info("Start server at port 5000...")
	if err := http.ListenAndServe(":5000", router); err != nil {
		log.Fatal(err)
	}
}
