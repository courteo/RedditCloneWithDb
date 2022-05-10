package main

//git config --global --add url."git@gitlab.com:courteo/lectures-2022-1.git".insteadOf "https://gitlab.com/courteo/lectures-2022-1.git"

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"redditclone/pkg/handlers"
	"redditclone/pkg/middleware"
	"redditclone/pkg/posts"
	"redditclone/pkg/posts/repo"
	"redditclone/pkg/session"
	"redditclone/pkg/user"
)

func main() {
	// основные настройки к базе
	dsn := "root:123456789@tcp(localhost:3306)/golang?"
	// указываем кодировку
	dsn += "&charset=utf8"
	// отказываемся от prapared statements
	// параметры подставляются сразу
	dsn += "&interpolateParams=true"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println(err.Error(), " ASDASDA")
	}

	db.SetMaxOpenConns(10)

	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		fmt.Println(err.Error(), " HHJFGGF")
	}

	// основные настройки к базе
	dsn1 := "root:123456789@tcp(localhost:3306)/golang?"
	// указываем кодировку
	dsn1 += "&charset=utf8"
	// отказываемся от prapared statements
	// параметры подставляются сразу
	dsn1 += "&interpolateParams=true"

	db1, err := sql.Open("mysql", dsn1)
	if err != nil {
		fmt.Println(err.Error(), " ASDASDA")
	}

	db1.SetMaxOpenConns(10)

	err = db1.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		fmt.Println(err.Error(), " HHJFGGF")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	// если коллекции не будет, то она создасться автоматически
	collection := client.Database("sample_training").Collection("posts")

	zapLogger, _ := zap.NewProduction()

	defer zapLogger.Sync()
	logger := zapLogger.Sugar()

	sessionManager := session.NewSessionsManager(db1)
	userRepo := user.NewMemoryRepo(db)
	userHandler := &handlers.UserHandler{
		UserRepo:       userRepo,
		Logger:         logger,
		SessionManager: sessionManager,
	}

	postRepo := posts.NewMemoryRepo(collection)
	Repo := repo.MyRepo{Db: postRepo}
	postHandler := &handlers.PostsHandler{
		PostRepo:       Repo,
		Logger:         logger,
		SessionManager: sessionManager,
	}

	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	r.Handle("/", http.FileServer(http.Dir("static/html/")))

	r.HandleFunc("/api/register", userHandler.Register).Methods("POST")
	r.HandleFunc("/api/login", userHandler.Login).Methods("POST")

	r.HandleFunc("/api/posts", middleware.Auth(postHandler.Add)).Methods("POST")
	r.HandleFunc("/api/post/{POST_ID:[0-9]+}", middleware.Auth(postHandler.AddComment)).Methods("POST")

	r.HandleFunc("/api/posts/", postHandler.GetAllPosts).Methods("GET")
	r.HandleFunc("/api/posts/{CATEGORY_NAME}", postHandler.GetCategory).Methods("GET")
	r.HandleFunc("/api/post/{POST_ID:[0-9]+}", postHandler.GetPost).Methods("GET")
	r.HandleFunc("/api/post/{POST_ID:[0-9]+}/upvote", middleware.Auth(postHandler.Upvote)).Methods("GET")
	r.HandleFunc("/api/post/{POST_ID:[0-9]+}/unvote", middleware.Auth(postHandler.Unvote)).Methods("GET")
	r.HandleFunc("/api/post/{POST_ID:[0-9]+}/downvote", middleware.Auth(postHandler.Downvote)).Methods("GET")
	r.HandleFunc("/api/user/{USER_LOGIN}", postHandler.GetUserPost).Methods("GET")

	r.HandleFunc("/api/post/{POST_ID:[0-9]+}/{COMMENT_ID:[0-9]+}", postHandler.DeleteComment).Methods("DELETE")
	r.HandleFunc("/api/post/{POST_ID:[0-9]+}", postHandler.Delete).Methods("DELETE")
	r.NotFoundHandler = http.HandlerFunc(NotHandler)
	//mux := middleware.Auth(r)

	addr := ":8080"
	logger.Infow("starting server",
		"type", "START",
		"addr", addr,
	)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", r)
}

func NotHandler(w http.ResponseWriter, r *http.Request) {
	data, _ := ioutil.ReadFile("static/html/index.html")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
