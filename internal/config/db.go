package config

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectDB(){
	var err error
	DB, err = pgxpool.New(context.Background(), os.Getenv("DB_URL"))
	if err != nil{
		log.Fatal("DB Connection failed: ",err)
	}
	log.Println("Connected to DB pool")
}