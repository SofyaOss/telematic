package postgres

import (
	"context"
	"log"
	"os"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	_, err := New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Println(err)
	}
}

func TestTelematicDB_DropTable(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	newDB, err := New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Println(err)
	}
	err = newDB.DropTable()
	if err != nil {
		log.Println("failed to drop table with error:", err)
	}
}

func TestTelematicDB_CreateTable(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	newDB, err := New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Println(err)
	}
	err = newDB.CreateTable()
	if err != nil {
		log.Println("failed to create table with error:", err)
	}
}

func TestTelematicDB_GetAllData(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	newDB, err := New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Println(err)
	}
	res, err := newDB.GetAllData()
	if err != nil {
		log.Println("failed to get all data with error:", err)
	} else {
		log.Println("GetAllData() result:", res)
	}
}
