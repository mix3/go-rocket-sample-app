package controller

import (
	"net/http"

	"github.com/acidlemon/rocket"
)

func TopPage(c rocket.CtxData) {
	c.Res().StatusCode = http.StatusOK
	c.RenderText("Hello World")
}

func PingPage(c rocket.CtxData) {
	c.Res().StatusCode = http.StatusOK
	c.RenderText("OK")
}
