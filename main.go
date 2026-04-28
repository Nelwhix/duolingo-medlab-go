package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Nelwhix/duolingo-medlab-go/handlers"
	"github.com/Nelwhix/duolingo-medlab-go/pkg"
	"github.com/Nelwhix/duolingo-medlab-go/pkg/mailer"
	"github.com/Nelwhix/duolingo-medlab-go/pkg/middleware"
	"github.com/Nelwhix/duolingo-medlab-go/pkg/models"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/go-playground/validator/v10"
	gHandlers "github.com/gorilla/handlers"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

var validate *validator.Validate

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := http.NewServeMux()

	dailyLogWriter := &pkg.DailyLogger{LogDir: "logs"}

	logger, err := pkg.CreateNewLogger(dailyLogWriter)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}

	validate = validator.New(validator.WithRequiredStructEnabled())

	pool, err := pkg.ConnectToDB()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()
	model := &models.Model{
		Conn: pool,
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		log.Fatalf("Unable to load sdk config, %v", err)
	}

	sesClient := sesv2.NewFromConfig(cfg)
	sesMailer := mailer.NewSESMailer(sesClient, os.Getenv("MAIL_FROM_ADDRESS"))

	handler := handlers.Handler{
		Model:     model,
		Logger:    logger,
		Validator: validate,
		Mailer:    sesMailer,
	}

	middleWare := middleware.Middleware{
		Model: model,
	}

	// Guest Routes
	r.HandleFunc("GET /api/v1/ping", handler.Pong)
	r.HandleFunc("POST /api/v1/auth/signup", handler.SignUp)
	r.HandleFunc("POST /api/v1/auth/login", handler.Login)
	r.HandleFunc("POST /api/v1/auth/forgot-password", handler.ForgotPassword)
	r.HandleFunc("POST /api/v1/auth/reset-password", handler.ResetPassword)
	r.Handle("PATCH /api/v1/users/{id}", middleWare.Auth(http.HandlerFunc(handler.UpdateUser)))
	r.Handle("GET /api/v1/departments", middleWare.Auth(http.HandlerFunc(handler.GetDepartments)))

	fs := http.FileServer(http.Dir("./swagger-ui"))
	r.Handle("GET /docs/", http.StripPrefix("/docs/", fs))

	fmt.Println("Server started on port 8080")
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PATCH", "PUT", "DELETE"},
		AllowCredentials: true,
		Debug:            false,
	})

	err = http.ListenAndServe(":8080", gHandlers.CombinedLoggingHandler(os.Stdout, c.Handler(middleware.ContentTypeMiddleware(r))))
	if err != nil {
		log.Printf("failed to run the server: %v", err)
	}
}
