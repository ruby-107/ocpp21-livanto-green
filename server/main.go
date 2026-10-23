package main

import (
	"log"
	"net/http"

	"ocpp21-livantogreen/internal/database"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load("../.env")
	_ = godotenv.Load(".env")

	db, err := database.Connect()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()

	repo := database.NewRepository(db)
	server := &Server{Repo: repo}

	http.HandleFunc("/ocpp/", server.websocketHandler)

	log.Println("CSMS server running on :9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
}
