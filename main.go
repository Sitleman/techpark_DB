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

	/*====================== THREAD ======================*/
	routerAPI.HandleFunc("/thread/{slug_or_id:[A-Za-z0-9._-]+}/create", handler.ThreadCreatePosts).Methods("POST")
	routerAPI.HandleFunc("/thread/{slug_or_id:[A-Za-z0-9._-]+}/vote", handler.ThreadVote).Methods("POST")
	routerAPI.HandleFunc("/thread/{slug_or_id:[A-Za-z0-9._-]+}/details", handler.ThreadDetails).Methods("GET")
	routerAPI.HandleFunc("/thread/{slug_or_id:[A-Za-z0-9._-]+}/details", handler.ThreadUpdate).Methods("POST")
	routerAPI.HandleFunc("/thread/{slug_or_id:[A-Za-z0-9._-]+}/posts", handler.ThreadPosts).Methods("GET")

	/*====================== POST ======================*/
	routerAPI.HandleFunc("/post/{id:[0-9]+}/details", handler.PostGet).Methods("GET")
	routerAPI.HandleFunc("/post/{id:[0-9]+}/details", handler.PostUpdate).Methods("POST")

	/*====================== USER ======================*/
	routerAPI.HandleFunc("/user/{nickname:[A-Za-z0-9._-]+}/create", handler.UserCreate).Methods("POST")
	routerAPI.HandleFunc("/user/{nickname:[A-Za-z0-9._-]+}/profile", handler.UserDetails).Methods("GET")
	routerAPI.HandleFunc("/user/{nickname:[A-Za-z0-9._-]+}/profile", handler.UserUpdate).Methods("POST")

	/*====================== SERVICE ======================*/
	routerAPI.HandleFunc("/service/status", handler.ServiceStatus).Methods("GET")
	routerAPI.HandleFunc("/service/clear", handler.ServiceClear).Methods("POST")

	log.Info("Start server at port 5000...")
	if err := http.ListenAndServe(":5000", router); err != nil {
		log.Fatal(err)
	}
}
