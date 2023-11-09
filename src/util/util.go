package util

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mailgun/mailgun-go/v4"
)

func NewAccountToken(bid, id string) (string, error) {
	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"bid": bid,
		"id":  id,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tkString, err := tk.SignedString([]byte(os.Getenv("JWTKEY")))
	if err != nil {
		return "", err
	}

	return tkString, nil
}

func VerifyToken(tkstring string) (jwt.MapClaims, error) {
	tk, err := jwt.Parse(tkstring, func(tk *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWTKEY")), nil
	})

	if err != nil {
		return nil, err
	}

	if !tk.Valid {
		return nil, errors.New("invalid")
	}

	claims := tk.Claims.(jwt.MapClaims)

	return claims, nil
}

func SendSimpleMessage(link string) (string, error) {
	// TODO: this must be env var
	domain := os.Getenv("MAILDOMAIN")
	apiKey := os.Getenv("MAILKEY")

	mg := mailgun.NewMailgun(domain, apiKey)

	sender := "Mailgun Sandbox <postmaster@sandbox8294e2fcbc994631a26f60f90df9a173.mailgun.org>"
	recipient := "tsn.teo@proton.me"

	subject := "FableFlow Conrfimation Request"
	body := `<!DOCTYPE html>
	<html lang="en">
	<body>
		<h1>Please activate your api access</h1>
		<a href="` + link + `">Click here!</a>
	</body>
	</html>`

	m := mg.NewMessage(sender, subject, "", recipient)

	m.SetHtml(body)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, id, err := mg.Send(ctx, m)

	fmt.Printf("%+v", resp)

	// TODO: remove returns and deal with errors
	return id, err
}

func RandHash() string {
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	charset := "abcdefghijklmnopqrstuvwxyz"

	var hash string
	for i := 0; i < 10; i++ {
		c := charset[seed.Intn(len(charset))]
		hash += string(c)
	}

	return hash
}

func addCorsHeader(res http.ResponseWriter) {
	headers := res.Header()
	headers.Add("Access-Control-Allow-Origin", "*")
	headers.Add("Vary", "Origin")
	headers.Add("Vary", "Access-Control-Request-Method")
	headers.Add("Vary", "Access-Control-Request-Headers")
	headers.Add("Access-Control-Allow-Headers", "*")
	headers.Add("Access-Control-Allow-Methods", "*")
}

func WriteJson(w http.ResponseWriter, status int, msg any) {
	w.Header().Set("Content-Type", "application/json")
	addCorsHeader(w)
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(msg)
	if err != nil {
		log.Printf("%+v", err)
	}
}
