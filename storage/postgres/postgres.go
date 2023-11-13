package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	pb "practice/internal/grpc"
	"practice/storage"
	"time"

	//"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	//"practice/storage"
)

//const (
//	databaseUrl string = "postgres://postgres:postgres@127.0.0.1:8081/telematicDB"
//)

type postgresDB struct {
	ctx  context.Context
	pool *pgxpool.Pool
}

type TelematicDB struct {
	db *pgxpool.Pool
}

func New(ctx context.Context, databaseUrl string) (*TelematicDB, error) {
	//db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer db.Close()
	//
	//_, err = db.Exec("CREATE TABLE IF NOT EXISTS $1 (id SERIAL PRIMARY KEY, name TEXT, email TEXT)", tableName)
	//
	//if err != nil {
	//	log.Fatal(err)
	//}

	//pool, err := pgxpool.Connect(context.Background(), databaseUrl)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//err = pool.Ping(context.Background())
	//if err != nil {
	//	log.Fatal(err)
	//}
	for {
		_, err := pgxpool.Connect(ctx, databaseUrl)
		if err == nil {
			break
		}
	}
	db, err := pgxpool.Connect(ctx, databaseUrl)
	if err != nil {
		return nil, err
	}
	t := TelematicDB{
		db: db,
	}
	return &t, nil
	//return &postgresDB{ctx: context.Background(), pool: pool}, nil
	//return nil
}

func (t *TelematicDB) CreateTable() error {
	_, err := t.db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS telematic (
			id SERIAL PRIMARY KEY, 
			car_number INT NOT NULL DEFAULT 1,
			speed INT NOT NULL DEFAULT 0,
			latitude FLOAT NOT NULL DEFAULT 0,
			longitude FLOAT NOT NULL DEFAULT 0,
			date DATE NOT NULL DEFAULT CURRENT_DATE
		);
	`)
	if err != nil {
		return err
	}
	return nil
}

func (t *TelematicDB) DropTable() error {
	_, err := t.db.Exec(context.Background(), `DROP TABLE IF EXISTS telematic;`)
	if err != nil {
		return err
	}
	return nil
}

// get all data
func (t *TelematicDB) GetAllData() ([]storage.Car, error) {
	rows, err := t.db.Query(context.Background(), `SELECT * FROM telematic`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var cars []storage.Car
	for rows.Next() {
		var car storage.Car
		err = rows.Scan(
			&car.ID,
			&car.Number,
			&car.Speed,
			&car.Coordinates.Latitude,
			&car.Coordinates.Longitude,
			&car.Date,
		)
		if err != nil {
			return nil, err
		}
		cars = append(cars, car)
	}
	return cars, nil
	//return func(w http.ResponseWriter, r *http.Request) {
	//	rows, err := db.Query("SELECT * FROM telematic")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	defer rows.Close()
	//
	//	cars := []storage.Car{}
	//	for rows.Next() {
	//		var c storage.Car
	//		if err := rows.Scan(&c.ID, &c.Number, &c.Speed, &c.Coords); err != nil {
	//			log.Fatal(err)
	//		}
	//		cars = append(cars, c)
	//	}
	//	if err := rows.Err(); err != nil {
	//		log.Fatal("4444444444444444", err)
	//	}
	//
	//	json.NewEncoder(w).Encode(cars)
	//}
}

func (t *TelematicDB) GetByDate(d1s, d2s string) ([]*pb.Car, error) {
	d1, err1 := time.Parse("2006-01-02", d1s)
	if err1 != nil {
		return nil, err1
	}
	d2, err2 := time.Parse("2006-01-02", d2s)
	if err2 != nil {
		return nil, err2
	}
	rows, err := t.db.Query(context.Background(), `SELECT * FROM telematic WHERE date BETWEEN $1 AND $2;`, d1, d2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var cars []*pb.Car
	for rows.Next() {
		var car storage.Car
		err = rows.Scan(
			&car.ID,
			&car.Number,
			&car.Speed,
			&car.Coordinates.Latitude,
			&car.Coordinates.Longitude,
			&car.Date,
		)
		if err != nil {
			return nil, err
		}
		cars = append(cars, CarPostgresToProto(&car))
	}
	return cars, nil
}

func (t *TelematicDB) GetByCarNumber(carNums []int) ([]*pb.Car, error) {
	var res []*pb.Car
	for _, num := range carNums {
		rows, err := t.db.Query(context.Background(), `SELECT * FROM telematic WHERE car_number=$1 ORDER BY id desc LIMIT 1;`, num)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var car storage.Car
			err = rows.Scan(
				&car.ID,
				&car.Number,
				&car.Speed,
				&car.Coordinates.Latitude,
				&car.Coordinates.Longitude,
				&car.Date,
			)
			if err != nil {
				return nil, err
			}
			res = append(res, CarPostgresToProto(&car))
		}
	}
	return res, nil
}

func GetByID(db *sql.DB, num, speed, coords int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var c storage.Car
		err := db.QueryRow("SELECT * FROM telematic WHERE number = $1", id).Scan(&c.ID, &c.Number, &c.Speed, &c.Coordinates.Latitude, &c.Coordinates.Longitude)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(c)
	}
}

func (t *TelematicDB) AddData(c *storage.Car) error {
	err := t.db.QueryRow(context.Background(),
		`INSERT INTO telematic (car_number, speed, latitude, longitude, date) VALUES ($1, $2, $3, $4, $5);`,
		c.Number, c.Speed, fmt.Sprint(c.Coordinates.Latitude), fmt.Sprint(c.Coordinates.Longitude), c.Date).Scan()
	return err
}

//func AddData(db *sql.DB) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		var c storage.Car
//		json.NewDecoder(r.Body).Decode(&c)
//
//		err := db.QueryRow("INSERT INTO telematic (car_number, speed, latitude, longitude, date) " +
//			"VALUES ($1, $2, $3, $4, $5) RETURNING id",
//			c.Number, c.Speed, c.Coords.Latitude, c.Coords.Longitude, c.Date).Scan(&c.ID)
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		json.NewEncoder(w).Encode(c)
//	}
//}

// get user by id
func getUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var c storage.Car
		err := db.QueryRow("SELECT * FROM telematic WHERE id = $1", id).Scan(&c.ID, &c.Number, &c.Speed, &c.Coordinates.Latitude, &c.Coordinates.Longitude)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(c)
	}
}
