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
		req, err := auth.CheckJWT(w, r)

		if err != nil {
			return
		}

		ctx := req.Context()
		auth0User := ctx.Value("user")

		userAuthID := auth0User.(*jwt.Token).Claims.(jwt.MapClaims)["sub"].(string)
		user, err := auth.FindOrCreateUser(userAuthID)

		if err != nil {
			httpError(w, 403, "auth", err.Error())
		}

		if user.IsAdmin == false {
			httpError(w, 403, "auth", "Admin only")
			return
		}

		ctx = context.WithValue(ctx, contextUser("user"), &user)

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
