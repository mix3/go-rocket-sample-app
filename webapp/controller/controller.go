package controller

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mix3/go-rocket-sample-app/webapp/models/cloudmailin"
	"github.com/mix3/go-rocket-sample-app/webapp/models/db"
	"github.com/mix3/go-rocket-sample-app/webapp/models/mailgun"
	"github.com/mix3/go-rocket-sample-app/webapp/models/timeparser"
	"github.com/mix3/rocket"
)

func serverError(c rocket.CtxData, err error) {
	c.Res().StatusCode = http.StatusInternalServerError
	c.RenderText(err.Error())
}

func TopPage(c rocket.CtxData) {
	c.Res().StatusCode = http.StatusOK
	value := rocket.RenderVars{
		"APP_NAME":         os.Getenv("APP_NAME"),
		"SIGNUP_ADDRESS":   os.Getenv("SIGNUP_ADDRESS"),
		"REGISTER_ADDRESS": os.Getenv("REGISTER_ADDRESS"),
	}
	c.Render("webapp/template/top.html", value)
}

func PingPage(c rocket.CtxData) {
	c.Res().StatusCode = http.StatusOK
	c.RenderText("OK")
}

func InterimRegisterEmailPage(c rocket.CtxData) {
	log.Printf("%#v", c)
	data, err := cloudmailin.Decode(c.Req().Body)
	if err != nil {
		serverError(c, err)
		return
	}

	hash, err := db.GetDB().InterimRegisterEmail(data.Envelope.From)
	if err != nil {
		client := mailgun.NewClient()
		res, err := client.Send(
			os.Getenv("APP_NAME"),
			os.Getenv("APP_ADDRESS"),
			data.Envelope.From,
			"ERROR",
			err.Error(),
		)
		log.Printf("res %v", res)
		log.Printf("err %v", err)
	} else {
		client := mailgun.NewClient()
		res, err := client.Send(
			os.Getenv("APP_NAME"),
			os.Getenv("APP_ADDRESS"),
			data.Envelope.From,
			"SUCCESS",
			os.Getenv("APP_URL")+"signup/"+hash,
		)
		log.Printf("res %v", res)
		log.Printf("err %v", err)
	}

	c.Res().StatusCode = http.StatusOK
	c.RenderText("OK")
}

func RegisterEmailPage(c rocket.CtxData) {
	err := db.GetDB().RegisterEmail(c.Params().Get("hash"))
	if err != nil {
		serverError(c, err)
		return
	}
	c.Res().StatusCode = http.StatusOK
	c.RenderText("OK")

}

func RegisterPage(c rocket.CtxData) {
	data, err := cloudmailin.Decode(c.Req().Body)
	if err != nil {
		serverError(c, err)
		return
	}
	log.Printf("%#v", data)

	at, err := timeparser.Parse(data.Headers.Subject)
	if err != nil {
		client := mailgun.NewClient()
		res, err := client.Send(
			os.Getenv("APP_NAME"),
			os.Getenv("APP_ADDRESS"),
			data.Envelope.From,
			"ERROR",
			err.Error(),
		)
		log.Printf("res %v", res)
		log.Printf("err %v", err)
	} else {
		db.GetDB().RegisterRemind(data.Envelope.From, data.Plain, at)
		client := mailgun.NewClient()
		res, err := client.Send(
			os.Getenv("APP_NAME"),
			os.Getenv("APP_ADDRESS"),
			data.Envelope.From,
			"SUCCESS",
			fmt.Sprintf("body: %v\n at: %v\n", data.Plain, at),
		)
		log.Printf("res %v", res)
		log.Printf("err %v", err)
	}

	c.Res().StatusCode = http.StatusOK
	c.RenderText("OK")
}
