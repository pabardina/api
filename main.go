package main

import (
	"context"
	"encoding/json"
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

	"github.com/hirondelle-app/api/models"
)

type contextUser string

type appHandlers struct {
	*gorm.DB                     `inject:""`
	*jwtmiddleware.JWTMiddleware `inject:""`
}

var config struct {
	Auth0Secret string `short:"s" long:"auth-secret" description:"The secret from Auth0" required:"true"`
	ServerPort  int    `short:"p" long:"server-port" description:"The server port" default:"8000" required:"true"`
	Database    struct {
		Address  string `long:"db-address" description:"The database address" default:"localhost" required:"true"`
		Username string `long:"db-user" description:"The database username" required:"true"`
		Password string `long:"db-password" description:"The database password" required:"true"`
		Name     string `long:"db-name" description:"The database name" required:"true"`
	}
}

func (customHandler *appHandlers) TwitterHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, _ := customHandler.CheckJWT(w, r)

		ctx := req.Context()
		auth0User := ctx.Value("user")

		user := models.User{}
		twitterID := auth0User.(*jwt.Token).Claims.(jwt.MapClaims)["sub"].(string)
		customHandler.Where("twitter_id = ?", twitterID).First(&user)

		// no user in DB
		if user.TwitterID == "" {
			// create user
			user.TwitterID = twitterID
			customHandler.Create(&user)
		}

		ctx = context.WithValue(ctx, contextUser("user"), &user)
		updatedReq := r.WithContext(ctx)

		h.ServeHTTP(w, updatedReq)
	})
}

func (customHandler *appHandlers) GetTweetsEndpoint(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	user := ctx.Value("user")

	fmt.Print(user)

	tweets, _ := models.GetTweetsByUser(customHandler.DB)

	if err := writeJSON(w, tweets, 200); err != nil {
		log.Fatal(err)
	}
}

func (customHandler *appHandlers) PostTweetEndpoint(w http.ResponseWriter, req *http.Request) {
	tweet := models.Tweet{}
	json.NewDecoder(req.Body).Decode(&tweet)

	//req.ParseForm()
	//KeywordLabel := req.Form.Get("hashtag")

	// check if tweet is correct
	if err := models.ValidateTweet(&tweet); err != nil {
		httpError(w, 400, "invalid_tweet", err.Error())
		return
	}

	// save tweet
	if err := models.CreateTweet(customHandler.DB, &tweet); err != nil {
		httpError(w, 400, "db_error", err.Error())
		return
	}

	// return tweet in response
	writeJSON(w, tweet, 201)
}

func (customHandler *appHandlers) PostKeywordEndpoint(w http.ResponseWriter, req *http.Request) {

	keyword := models.Keyword{}
	json.NewDecoder(req.Body).Decode(&keyword)

	if keyword.Label == "" {
		httpError(w, 400, "invalid_keyword", "Label must not be empty")
		return
	}

	if err := models.CreateKeyword(customHandler.DB, &keyword); err != nil {
		httpError(w, 400, "db_error", err.Error())
		return
	}

	writeJSON(w, keyword, 201)
}

func (customHandler *appHandlers) GetAllKeywordsEndpoint(w http.ResponseWriter, req *http.Request) {

	keywords, _ := models.GetKeywords(customHandler.DB)

	if err := writeJSON(w, keywords, 200); err != nil {
		log.Fatal(err)
	}
}

func httpError(w http.ResponseWriter, code int, msg, description string) {
	writeJSON(w, struct {
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
	}{msg, description}, code)
}

func writeJSON(w http.ResponseWriter, data interface{}, code int) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(b)

	return nil
}

func init() {
	_, err := flags.Parse(&config)

	if err != nil {
		panic(err)
	}

}

func main() {
	db, err := gorm.Open("postgres",
		fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable",
			config.Database.Address, config.Database.Username, config.Database.Name, config.Database.Password))

	if err != nil {
		panic("failed to connect database")
	}

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Auth0Secret), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	appHandlers := &appHandlers{}

	if err := inject.Populate(db, jwtMiddleware, appHandlers); err != nil {
		panic(err)
	}

	// Migrate the schema
	db.AutoMigrate(&models.Tweet{}, &models.Keyword{}, &models.User{})

	// delete
	keyword := models.Keyword{}
	keyword.Label = "dsd"

	models.CreateKeyword(db, &keyword)

	db.Create(&models.Tweet{TweetID: "2", Likes: 2, Retweets: 3, KeywordID: keyword.ID})
	db.Create(&models.Tweet{TweetID: "3", Likes: 3, Retweets: 4, KeywordID: keyword.ID})
	// delete

	router := mux.NewRouter()
	//router.Handle("/tweets", appHandlers.TwitterHandler(http.HandlerFunc(appHandlers.GetTweetsEndpoint))).Methods("GET")
	router.Handle("/tweets", http.HandlerFunc(appHandlers.GetTweetsEndpoint)).Methods("GET")
	router.Handle("/tweets", http.HandlerFunc(appHandlers.PostTweetEndpoint)).Methods("POST")

	router.HandleFunc("/keywords", appHandlers.PostKeywordEndpoint).Methods("POST")
	router.HandleFunc("/keywords", appHandlers.GetAllKeywordsEndpoint).Methods("GET")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.ServerPort), router))
}
