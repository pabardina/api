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
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/hirondelle-app/api/models"
)

type AppHandlers struct {
	*gorm.DB                     `inject:""`
	*jwtmiddleware.JWTMiddleware `inject:""`
}

type userKey int

func (ah *AppHandlers) TwitterHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, err := ah.CheckJWT(w, r)

		if err != nil {
			return
		}

		ctx := req.Context()
		auth0User := ctx.Value("user")

		user := models.User{}
		twitterID := auth0User.(*jwt.Token).Claims.(jwt.MapClaims)["sub"].(string)
		ah.Where("twitter_id = ?", twitterID).First(&user)

		// no user in DB
		if user.TwitterID == "" {
			// create user
			user.TwitterID = twitterID
			ah.Create(&user)
		}

		ctx = context.WithValue(ctx, "user", &user)
		toto := r.WithContext(ctx)

		h.ServeHTTP(w, toto)
	})
}

func (h *AppHandlers) GetTweetsEndpoint(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	user := ctx.Value("user")

	fmt.Print(user)

	tweets, _ := models.GetTweetsByUser(h.DB)

	if err := writeJSON(w, tweets, 200); err != nil {
		log.Fatal(err)
	}
	fmt.Print("hello")
}

func (h *AppHandlers) PostKeywordEndpoint(w http.ResponseWriter, req *http.Request) {

	keyword := models.Keyword{}
	keyword.Label = "toto"
	err := models.CreateKeyword(h.DB, &keyword)

	if err != nil {
		log.Fatal(err)
	}

	if err := writeJSON(w, keyword, 204); err != nil {
		log.Fatal(err)
	}
}

func (h *AppHandlers) GetAllKeywordsEndpoint(w http.ResponseWriter, req *http.Request) {

	keywords, _ := models.GetKeywords(h.DB)

	if err := writeJSON(w, keywords, 200); err != nil {
		log.Fatal(err)
	}
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

func main() {

	db, err := gorm.Open("postgres", "host=localhost user=toto dbname=toto password=toto sslmode=disable")
	if err != nil {
		panic("failed to connect database")
	}

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("xxxxxxx"), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	appHandlers := &AppHandlers{}

	if err := inject.Populate(db, jwtMiddleware, appHandlers); err != nil {
		panic(err)
	}

	// Migrate the schema
	db.AutoMigrate(&models.Tweet{}, &models.Keyword{}, &models.User{})

	// delete
	keyword := models.Keyword{}
	keyword.Label = "dsd"

	err = models.CreateKeyword(db, &keyword)

	db.Create(&models.Tweet{TweetID: "2", Likes: 2, Retweets: 3, KeywordID: keyword.ID})
	db.Create(&models.Tweet{TweetID: "3", Likes: 3, Retweets: 4, KeywordID: keyword.ID})
	// delete

	router := mux.NewRouter()
	router.Handle("/tweets", appHandlers.TwitterHandler(http.HandlerFunc(appHandlers.GetTweetsEndpoint))).Methods("GET")

	//router.HandleFunc("/keywords", appHandlers.PostKeywordEndpoint).Methods("POST")
	//router.HandleFunc("/keywords", appHandlers.GetAllKeywordsEndpoint).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}
