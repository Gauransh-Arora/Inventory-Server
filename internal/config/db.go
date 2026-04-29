package config

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn

func ConnectDB(){
	var err error
	dbURL := os.Getenv("DB_URL")
    log.Println("DB_URL:", dbURL) // add this
	DB, err = pgx.Connect(context.Background(), os.Getenv("DB_URL"))
	if err != nil{
		log.Fatal("DB Connection failed: ",err)
	}
	log.Println("Connected to DB")
}