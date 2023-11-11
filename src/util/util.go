package util

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"html/template"
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
		"exp": time.Now().Add(time.Second * 24).Unix(),
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

func SendSimpleMessage(link, recipient string) {
	domain := os.Getenv("MAILDOMAIN")
	apiKey := os.Getenv("MAILKEY")

	mg := mailgun.NewMailgun(domain, apiKey)

	sender := "Mailgun Sandbox <postmaster@sandbox8294e2fcbc994631a26f60f90df9a173.mailgun.org>"
	subject := "FableFlow Conrfimation Request"

	m := mg.NewMessage(sender, subject, "", recipient)

	t, err := template.ParseFiles("src/tmpl/email.html")
	if err != nil {
		log.Println(err)
		return
	}

	data := struct {
		Link string
	}{
		Link: link,
	}

	buf := new(bytes.Buffer)
	if err := t.Execute(buf, data); err != nil {
		log.Println(err, buf)
		return
	}

	m.SetHtml(buf.String())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, id, err := mg.Send(ctx, m)
	if err != nil {
		log.Printf("resp: %s\nid: %s\nerr: %v", resp, id, err)
	}
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

func WriteJson(w http.ResponseWriter, status int, msg any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(msg)
	if err != nil {
		log.Printf("%+v", err)
	}
}
