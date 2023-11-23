package postgres

import (
	"context"
	"fmt"
	"time"

	pb "practice/internal/grpc"
	"practice/storage"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type TelematicDB struct {
	db *pgxpool.Pool
}

func New(ctx context.Context, databaseUrl string) (*TelematicDB, error) {
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
}

func (t *TelematicDB) CreateTable(ctx context.Context) error {
	_, err := t.db.Exec(ctx, `
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

func (t *TelematicDB) DropTable(ctx context.Context) error {
	_, err := t.db.Exec(ctx, `DROP TABLE IF EXISTS telematic;`)
	if err != nil {
		return err
	}
	return nil
}

// get all data
func (t *TelematicDB) GetAllData(ctx context.Context) ([]storage.Car, error) {
	rows, err := t.db.Query(ctx, `SELECT * FROM telematic`)
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
}

func (t *TelematicDB) GetByDate(d1s, d2s string, nums []int64, ctx context.Context) ([]*pb.Car, error) {
	d1, err1 := time.Parse("2006-01-02", d1s)
	if err1 != nil {
		return nil, err1
	}
	d2, err2 := time.Parse("2006-01-02", d2s)
	if err2 != nil {
		return nil, err2
	}
	rows, err := t.db.Query(ctx, `SELECT * FROM telematic WHERE (date BETWEEN $1 AND $2) AND (car_number = ANY($3));`, d1, d2, pq.Array(nums))
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

func (t *TelematicDB) GetByCarNumber(ctx context.Context, carNums []int64) ([]*pb.Car, error) {
	var cars []*pb.Car
	for _, num := range carNums {
		rows, err := t.db.Query(ctx, `SELECT * FROM telematic WHERE car_number=$1 ORDER BY id desc LIMIT 1;`, num)
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
			cars = append(cars, CarPostgresToProto(&car))
		}
	}
	return cars, nil
}

func (t *TelematicDB) AddData(c *storage.Car, ctx context.Context) error {
	err := t.db.QueryRow(ctx,
		`INSERT INTO telematic (car_number, speed, latitude, longitude, date) VALUES ($1, $2, $3, $4, $5);`,
		c.Number, c.Speed, fmt.Sprint(c.Coordinates.Latitude), fmt.Sprint(c.Coordinates.Longitude), c.Date).Scan()
	return err
}
