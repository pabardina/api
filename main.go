package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aiden0z/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/facebookgo/inject"
	"github.com/gorilla/mux"
	"github.com/jessevdk/go-flags"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/hirondelle-app/api/api"
	"github.com/hirondelle-app/api/common"
	"github.com/hirondelle-app/api/tweets"
	"github.com/hirondelle-app/api/users"
)

func init() {
	_, err := flags.Parse(&common.Config)
	if err != nil {
		panic(err)
	}
}

func main() {
	db := initDatabase()
	jwtMiddleware := initJwtMiddleware()
	tweetsHandlers := &api.TweetsHandlers{}
	authMiddleware := &api.AuthMiddleware{}
	tweetsManager := &tweets.Manager{}
	usersManager := &users.Manager{}
	if err := inject.Populate(db, jwtMiddleware, tweetsHandlers, authMiddleware, tweetsManager, usersManager); err != nil {
		panic(err)
	}

	router := mux.NewRouter()
	//router.Handle("/tweets", authMiddleware.Use(http.HandlerFunc(tweetsHandlers.GetTweetsEndpoint))).Methods("GET")
	router.Handle("/tweets", http.HandlerFunc(tweetsHandlers.GetTweetsEndpoint)).Methods("GET")
	router.Handle("/tweets", http.HandlerFunc(tweetsHandlers.PostTweetEndpoint)).Methods("POST")

	router.HandleFunc("/keywords", tweetsHandlers.PostKeywordEndpoint).Methods("POST")
	router.HandleFunc("/keywords", tweetsHandlers.GetAllKeywordsEndpoint).Methods("GET")

	fmt.Printf("Starting server on port %v\n", common.Config.ServerPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", common.Config.ServerPort), router))
}

func initDatabase() *gorm.DB {
	db, err := gorm.Open("postgres",
		fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable",
			common.Config.Database.Address,
			common.Config.Database.Username,
			common.Config.Database.Name,
			common.Config.Database.Password))
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&tweets.Tweet{}, &tweets.Keyword{}, &users.User{})
	return db
}

func initJwtMiddleware() *jwtmiddleware.JWTMiddleware {
	return jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(common.Config.Auth0Secret), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
}
