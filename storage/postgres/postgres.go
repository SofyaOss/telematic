package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"practice/storage"

	//"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"log"
	//"practice/storage"
)

//const (
//	databaseUrl string = "postgres://postgres:postgres@127.0.0.1:8081/telematicDB"
//)

type postgresDB struct {
	ctx  context.Context
	pool *pgxpool.Pool
}

func New() error {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS telematic (id SERIAL PRIMARY KEY, name TEXT, email TEXT)")

	if err != nil {
		log.Fatal(err)
	}
	//pool, err := pgxpool.Connect(context.Background(), databaseUrl)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//err = pool.Ping(context.Background())
	//if err != nil {
	//	log.Fatal(err)
	//}
	//return &postgresDB{ctx: context.Background(), pool: pool}, nil
	return nil
}

// get all data
func GetAllData(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM telematic")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		cars := []storage.Car{}
		for rows.Next() {
			var c storage.Car
			if err := rows.Scan(&c.ID, &c.Number, &c.Speed, &c.Coords, c.Date); err != nil {
				log.Fatal(err)
			}
			cars = append(cars, c)
		}
		if err := rows.Err(); err != nil {
			log.Fatal("4444444444444444", err)
		}

		json.NewEncoder(w).Encode(cars)
	}
}

// get user by id
func getUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var c storage.Car
		err := db.QueryRow("SELECT * FROM telematic WHERE id = $1", id).Scan(&c.ID, &c.Number, &c.Speed, &c.Coords, c.Date)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(c)
	}
}
