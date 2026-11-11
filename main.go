package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/jozuacydsagradojuariojcsj/go_api_practice/internal/api"
	"github.com/jozuacydsagradojuariojcsj/go_api_practice/internal/postgres"
	"github.com/jozuacydsagradojuariojcsj/go_api_practice/internal/store"
)

func main() {
	_ = godotenv.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := postgres.Connect(ctx)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer pool.Close()
	log.Println("Starting GO API service...")

	log.Printf(`
 ______     ______        ______     ______   __    
/\  ___\   /\  __ \      /\  __ \   /\  == \ /\ \   
\ \ \__ \  \ \ \/\ \     \ \  __ \  \ \  _-/ \ \ \  
 \ \_____\  \ \_____\     \ \_\ \_\  \ \_\    \ \_\ 
  \/_____/   \/_____/      \/_/\/_/   \/_/     \/_/ `)

	s := store.New(pool)
	s := api.New(s)

	err := http.ListenAndServe(url, r)
	if err != nil {
		log.Error(err)
	}
}
