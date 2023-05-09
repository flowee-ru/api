package users

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/flowee-ru/flowee-api/models"
	"github.com/flowee-ru/flowee-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

func contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

func Register(ctx context.Context, w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	captcha := r.FormValue("captcha")

	if username == "" || password == "" || email == "" || captcha == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"success": false}`)
		return
	}

	if len(username) < 3 || len(username) > 15 {
		fmt.Fprintf(w, `{"success": false, "errorCode": 1}`)
		return
	}

	test, _ := regexp.MatchString("^\\s*$", username)
	if test {
		fmt.Fprintf(w, `{"success": false, "errorCode": 1}`)
		return
	}

	badUsernames := []string {
		"settings",
		"verify",
	}

	if contains(badUsernames, username) {
		fmt.Fprintf(w, `{"success": false, "errorCode": 2}`)
		return
	}

	check, err := utils.VerifyCaptcha(captcha)
	if err != nil || !check {
		fmt.Fprintf(w, `{"success": false, "errorCode": 3}`)
		return
	}

	err = db.Collection("accounts").FindOne(ctx, bson.D{primitive.E{Key: "username", Value: username}}).Decode(nil)
	if err != mongo.ErrNoDocuments {
		fmt.Fprintf(w, `{"success": false, "errorCode": 4}`)
		return
	}

	err = db.Collection("accounts").FindOne(ctx, bson.D{primitive.E{Key: "email", Value: email}}).Decode(nil)
	if err != mongo.ErrNoDocuments {
		fmt.Fprintf(w, `{"success": false, "errorCode": 5}`)
		return
	}

	verifyToken := utils.GenerateToken(10)
	verifyLink := os.Getenv("APP_HOST") + "/verify?token=" + verifyToken

	mail := gomail.NewMessage()
	mail.SetHeader("From", os.Getenv("SMTP_USER"))
	mail.SetHeader("To", email)
	mail.SetHeader("Subject", "Verify your account on Flowee")
	mail.SetBody("text/html", `<div align="center">Hello ` + username + `! Please follow this link to activate your account on Flowee:<br><a href="` + verifyLink + `">` + verifyLink + `</div>`)

	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	dialer := gomail.NewDialer(os.Getenv("SMTP_HOST"), smtpPort, os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASSWORD"))

	err = dialer.DialAndSend(mail)
	if err != nil {
		fmt.Fprintf(w, `{"success": false, "errorCode": 6}`)
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	db.Collection("accounts").InsertOne(ctx, models.Account{
		ID: primitive.NewObjectID(),
		Username: username,
		Password: string(hash),
		Email: email,
		Timestamp: int32(time.Now().Unix()),
		Avatar: "default",
		LastStream: 0,
		LastEmailSend: 0,
		Token: utils.GenerateToken(30),
		VerifyToken: verifyToken,
		StreamToken: utils.GenerateToken(30),
		StreamName: username + "'s stream",
		Viewers: 0,
		IsLive: false,
		IsActive: false,
	})

	fmt.Fprintf(w, `{"success": true}`)
}