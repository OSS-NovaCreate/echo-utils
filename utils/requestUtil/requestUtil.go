package requestUtil

import (
	"github.com/labstack/echo/v4"
)

type RequestData struct {
	Context   echo.Context
	UserAgent string
}

var Request RequestData

func SetContext(c echo.Context) {
	Request.Context = c
}

func GetContext() echo.Context {
	return Request.Context
}

func SetUserAgent(u string) {
	Request.UserAgent = u
}

func GetUserAgent() string {
	return Request.UserAgent
}
