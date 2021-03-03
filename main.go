package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"

	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/manabie-com/togo/internal/authorizers"
	"github.com/manabie-com/togo/internal/handlers"
	"github.com/manabie-com/togo/internal/seed"
	"github.com/manabie-com/togo/internal/services"
	"github.com/manabie-com/togo/internal/storages/postgres"
)

func main() {
	// db, err := sql.Open("sqlite3", "./data.db")
	// if err != nil {
	// 	log.Fatal("error opening db", err)
	// }

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	}

	dbURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"))
	pgDB, err := gorm.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Cannot connect to postgres database. Error: ", err)
	}

	seed.Load(pgDB)

	mux := http.NewServeMux()

	jwtAuthorizer := &authorizers.JWTAuthorizer{
		JWTKey: "wqGyEBBfPK9w3Lxw",
	}

	taskService := &services.TaskServiceImpl{
		Store: &postgres.PostgresDB{
			DB: pgDB,
		},
	}

	userService := &services.UserServiceImpl{
		Store: &postgres.PostgresDB{
			DB: pgDB,
		},
		Authorizer: jwtAuthorizer,
	}

	taskHandler := &handlers.TaskHandler{TaskService: taskService}

	authorizationHandler := &handlers.AuthorizationHandler{Authorizer: jwtAuthorizer, UserService: userService}

	mux.Handle("/login", authorizationHandler)
	mux.Handle("/tasks", authorizationHandler.HandleAuthorization(taskHandler))
	log.Println("Listening on :5050...")
	err = http.ListenAndServe(":5050", mux)
	log.Fatal(err)
}
