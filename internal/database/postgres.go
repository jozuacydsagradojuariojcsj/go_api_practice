package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)
type User struct {
	ID int
	Name string
	Coins string
}


func main() {
	ctx := context.Background()

	_ = godotenv.Load()

	cfg,err := pgxpool.ParseConfig("")
	if err != nil {
		log.Errorf("parse config %v", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Errorf("create pool %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil{
		log.Errorf("ping %v", err)
	}

	fmt.Print("Connected")
}