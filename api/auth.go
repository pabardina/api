package api

import (
	"context"
	"net/http"

	jwtmiddleware "github.com/aiden0z/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/hirondelle-app/api/users"
)

type contextUser string

type AuthMiddleware struct {
	*jwtmiddleware.JWTMiddleware `inject:""`
	*users.Manager               `inject:""`
}

func (auth *AuthMiddleware) Use(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, _ := auth.CheckJWT(w, r)

		ctx := req.Context()
		auth0User := ctx.Value("user")

		twitterID := auth0User.(*jwt.Token).Claims.(jwt.MapClaims)["sub"].(string)
		user, err := auth.FindOrCreateUserForTwitterID(twitterID)
		if err != nil {
			// TODO
			panic(err)
		}

		ctx = context.WithValue(ctx, contextUser("user"), &user)

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
