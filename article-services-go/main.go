package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func main() {
    dsn := "root:@tcp(127.0.0.1:3306)/article?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("Gagal Connect ke Database")
    }
    DB = db

    // migrate (sync model ke db)
    err = DB.AutoMigrate(&Post{})
    if err != nil {
        fmt.Println("Migrate Gagal:", err)
    } else {
        fmt.Println("Database Migrate Sukses")
    }

	router := mux.NewRouter()
    router.HandleFunc("/article", CreatePost).Methods("POST")
    router.HandleFunc("/article/{limit:[0-9]+}/{offset:[0-9]+}", GetPosts).Methods("GET")
    router.HandleFunc("/article/status/{status}/{limit:[0-9]+}/{offset:[0-9]+}", GetPostsByStatus).Methods("GET")
    router.HandleFunc("/article/{id:[0-9]+}", GetPost).Methods("GET")
    router.HandleFunc("/article/{id:[0-9]+}", UpdatePost).Methods("PUT")
    router.HandleFunc("/article/{id:[0-9]+}/status", UpdatePostStatus).Methods("PATCH")
    router.HandleFunc("/article/{id:[0-9]+}", DeletePost).Methods("DELETE")

    fmt.Println("Server running on port 8080")

    headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
    originsOk := handlers.AllowedOrigins([]string{"http://localhost:3000"})
    methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"})

    http.ListenAndServe(":8080", handlers.CORS(originsOk, headersOk, methodsOk)(router))

}