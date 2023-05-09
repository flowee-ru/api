package users

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/flowee-ru/flowee-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/gomail.v2"
)

func ResendEmail(ctx context.Context, w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	email := r.FormValue("email")

	if email == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"success": false}`)
		return
	}

	var account models.Account
	err := db.Collection("accounts").FindOne(ctx, bson.D{primitive.E{Key: "email", Value: email}, primitive.E{Key: "isActive", Value: false}}).Decode(&account)
	if err == mongo.ErrNoDocuments {
		fmt.Fprintf(w, `{"success": false, "errorCode": 1}`)
		return
	}

	if account.LastEmailSend + 300 > int32(time.Now().Unix()) {
		fmt.Fprintf(w, `{"success": false, "errorCode": 2}`)
		return
	}

	verifyLink := os.Getenv("APP_HOST") + "/verify?token=" + account.VerifyToken

	mail := gomail.NewMessage()
	mail.SetHeader("From", os.Getenv("SMTP_USER"))
	mail.SetHeader("To", email)
	mail.SetHeader("Subject", "Verify your account on Flowee")
	mail.SetBody("text/html", `<div align="center">Hello ` + account.Username + `! Please follow this link to activate your account on Flowee:<br><a href="` + verifyLink + `">` + verifyLink + `</div>`)

	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	dialer := gomail.NewDialer(os.Getenv("SMTP_HOST"), smtpPort, os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASSWORD"))

	err = dialer.DialAndSend(mail)
	if err != nil {
		fmt.Fprintf(w, `{"success": false, "errorCode": 3}`)
		return
	}

	db.Collection("accounts").UpdateOne(ctx, bson.D{primitive.E{Key: "email", Value: email}},
	bson.D{primitive.E{
		Key: "$set", Value: bson.D{primitive.E{
			Key: "lastEmailSend", Value: int32(time.Now().Unix()),
		}},
	}})

	fmt.Fprintf(w, `{"success": true}`)
}