package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jozuacydsagradojuariojcsj/go_api_practice/internal/handlers"
	log "github.com/sirupsen/logrus"
)

func main(){
	url := "localhost:8000"
	log.SetReportCaller(true)
	var r *chi.Mux = chi.NewRouter()
	handlers.Handler(r)

	fmt.Println("Starting GO API service...")

	fmt.Printf(`
 ______     ______        ______     ______   __    
/\  ___\   /\  __ \      /\  __ \   /\  == \ /\ \   
\ \ \__ \  \ \ \/\ \     \ \  __ \  \ \  _-/ \ \ \  
 \ \_____\  \ \_____\     \ \_\ \_\  \ \_\    \ \_\ 
  \/_____/   \/_____/      \/_/\/_/   \/_/     \/_/ 
  
 Started at %v`,url)

  err := http.ListenAndServe(url, r)
  if err != nil {
	  log.Error(err)
  }
}