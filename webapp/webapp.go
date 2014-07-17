package webapp

import (
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/mix3/go-rocket-sample-app/webapp/controller"
	"github.com/mix3/go-rocket-sample-app/webapp/models/db"
	"github.com/mix3/go-rocket-sample-app/webapp/models/mailgun"
	"github.com/mix3/go-rocket-sample-app/webapp/view"
	"github.com/mix3/rocket"
)

type WebApp struct {
	rocket.WebApp
}

func NewWebApp() WebApp {
	app := WebApp{}
	app.Init()

	app.AddRoute(
		"/",
		//func(c rocket.CtxData) {
		//	controller.TopPage(c)
		//},
		controller.TopPage,
		view.Render{},
	)

	app.AddRoute(
		"/ping",
		controller.PingPage,
		view.Render{},
	)

	app.AddRoute(
		"/signup",
		controller.InterimRegisterEmailPage,
		view.Render{},
	)

	app.AddRoute(
		"/signup/:hash",
		controller.RegisterEmailPage,
		view.Render{},
	)

	app.AddRoute(
		"/register",
		controller.RegisterPage,
		view.Render{},
	)

	app.BuildRouter()

	return app
}

func Start(listener net.Listener) {
	go ping()
	go kick()
	app := NewWebApp()
	app.Start(listener)
}

func ping() {
	ticker := time.NewTicker(time.Minute * 50)
	for {
		select {
		case <-ticker.C:
			http.Get(os.Getenv("PING_URL"))
		}
	}
}

func kick() {
	ticker := time.NewTicker(time.Minute * 1)
	for {
		select {
		case <-ticker.C:
			// remind
			log.Printf("kick")
			ret := db.GetDB().RemindList()
			client := mailgun.NewClient()
			for _, v := range ret {
				log.Printf("%#v", v)
				res, err := client.Send(
					os.Getenv("APP_NAME"),
					os.Getenv("APP_ADDRESS"),
					v.To,
					"remind",
					v.Message,
				)
				log.Printf("%#v", res)
				log.Printf("%#v", err)
				if err != nil {
					log.Printf("send error %#v", err)
				}
				err = db.GetDB().DeleteRemind(v)
				if err != nil {
					log.Printf("delete error %#v", err)
				}
			}
		}
	}
}
