package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tsntt/fableflow/src/util"
	"github.com/uptrace/bunrouter"
)

type CtxKey string

func (k CtxKey) String() string { return string(k) }

var (
	AccID = CtxKey("id")
)

func AccountAuth(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, r bunrouter.Request) error {
		AID := r.Header.Get("Fableflowaid")

		claims, err := util.VerifyToken(AID)

		if claims["id"].(string) == "" || claims["bid"].(string) == "" {
			util.WriteJson(w, http.StatusUnauthorized, "invalid token")
			return errors.New("invalid token")
		}

		if errors.Is(err, jwt.ErrTokenExpired) {
			token, err := util.NewAccountToken(claims["bid"].(string), claims["id"].(string))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return err
			}

			w.Header().Set("Fableflowaid", token)
		} else if err != nil {
			util.WriteJson(w, http.StatusUnauthorized, "invalid token")
			return err
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, AccID, claims["id"].(string))

		return next(w, r.WithContext(ctx))
	}
}
