package webapp

import (
	"net"
	"net/http"
	"os"
	"time"

	"github.com/acidlemon/rocket"
	"github.com/mix3/go-rocket-sample-app/webapp/controller"
	"github.com/mix3/go-rocket-sample-app/webapp/view"
)

type WebApp struct {
	rocket.WebApp
}

func NewWebApp() WebApp {
	app := WebApp{}
	app.Init()

	app.AddRoute(
		"/",
		//		func(c rocket.CtxData) {
		//			controller.TopPage(c)
		//		},
		controller.TopPage,
		view.Render{},
	)

	app.AddRoute(
		"/ping",
		controller.PingPage,
		view.Render{},
	)

	app.BuildRouter()

	return app
}

func Start(listener net.Listener) {
	go ping()
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
