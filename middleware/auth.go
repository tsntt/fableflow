package middleware

import (
	"context"
	"errors"
	"net/http"

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
		// origin := r.Header.Get("Origin")
		// if origin == "" {
		// 	return next(w, r)
		// }

		AID := r.Header.Get("Fableflowaid")

		claims, err := util.VerifyToken(AID)
		if err != nil {
			return err
		}

		if claims["id"].(string) == "" || claims["bid"].(string) == "" {
			return errors.New("invalid token")
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, AccID, claims["id"].(string))

		// check host and domain
		return next(w, r.WithContext(ctx))
	}
}
