package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"practice/storage/postgres"
)

func main() {
	//connect to database
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//create the table if it doesn't exist
	log.Println("im working")
	//_, err = db.Exec("DROP TABLE IF EXISTS telematic")

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS telematic (id SERIAL PRIMARY KEY, number INT, speed INT, coords INT)")
	// , pub_date BIGINT DEFAULT 0)
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/data", postgres.GetAllData(db)).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", jsonContentTypeMiddleware(router)))
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
