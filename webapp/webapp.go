package webapp

import (
	"net"

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

	app.BuildRouter()

	return app
}

func Start(listener net.Listener) {
	app := NewWebApp()
	app.Start(listener)
}
