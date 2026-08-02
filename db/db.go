package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

var Conn *pgx.Conn

func Connect() {
	dbURL := os.Getenv("DB_URL")

	conn, err := pgx.Connect(context.Background(), dbURL)

	if err != nil {
		log.Fatal("Error connecting to db", err)
	}

	Conn = conn
	fmt.Println("Succesfully Connected to the DB")
}

func Close() {
	if Conn != nil {
		Conn.Close(context.Background())
	}
}
