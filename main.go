package main

import (
	"context"
	"encoding/hex"
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
	"github.com/gorilla/securecookie"
	"github.com/joho/godotenv"
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

	sessionHash := os.Getenv("SESSION_HASH_KEY")
	sessionBlockKeyHex := os.Getenv("SESSION_BLOCK_KEY")
	sessionBlockKey, err := hex.DecodeString(sessionBlockKeyHex)
	if err != nil {
		log.Fatalf("Invalid SESSION_BLOCK_KEY: %v", err)
	}
	cookieHandler := securecookie.New([]byte(sessionHash), sessionBlockKey)

	handler := handlers.Handler{
		Model:         model,
		Logger:        logger,
		Validator:     validate,
		Mailer:        sesMailer,
		CookieHandler: cookieHandler,
	}

	middleWare := middleware.Middleware{
		Model:         model,
		CookieHandler: cookieHandler,
		Logger:        logger,
	}

	// Guest Routes
	r.HandleFunc("GET /admin/login", handler.RenderAdminLoginPage)
	r.HandleFunc("POST /auth/login", handler.Login)
	r.HandleFunc("GET /", handler.Pong)

	// auth routes
	r.Handle("PATCH /api/v1/users/{id}", middleWare.Auth(http.HandlerFunc(handler.UpdateUser)))
	r.Handle("GET /api/v1/departments", middleWare.Auth(http.HandlerFunc(handler.GetDepartments)))

	// admin routes
	r.Handle("GET /admin/dashboard", middleWare.SessionAuth(middleWare.Admin(http.HandlerFunc(handler.RenderAdminDashboard))))
	r.Handle("POST /admin/logout", middleWare.SessionAuth(middleWare.Admin(http.HandlerFunc(handler.Logout))))

	// topics
	r.Handle("POST /admin/topics", middleWare.SessionAuth(middleWare.Admin(http.HandlerFunc(handler.CreateTopic))))

	// questions
	r.Handle("POST /admin/questions", middleWare.SessionAuth(middleWare.Admin(http.HandlerFunc(handler.CreateQuestion))))
	r.Handle("GET /admin/questions", middleWare.SessionAuth(middleWare.Admin(http.HandlerFunc(handler.RenderAdminQuestions))))
	r.Handle("GET /admin/questions/{id}", middleWare.SessionAuth(middleWare.Admin(http.HandlerFunc(handler.RenderSingleAdminQuestion))))
	r.Handle("GET /admin/questions/create", middleWare.SessionAuth(middleWare.Admin(http.HandlerFunc(handler.RenderAdminCreateQuestion))))
	r.Handle("POST /admin/questions/{id}/delete", middleWare.SessionAuth(middleWare.Admin(http.HandlerFunc(handler.DeleteQuestion))))

	static := http.FileServer(http.Dir("./static"))
	r.Handle("GET /static/", http.StripPrefix("/static/", static))

	fmt.Println("Server started on port 8080")

	err = http.ListenAndServe(":8080", gHandlers.CombinedLoggingHandler(os.Stdout, middleware.Brotli(r)))
	if err != nil {
		log.Printf("failed to run the server: %v", err)
	}
}
