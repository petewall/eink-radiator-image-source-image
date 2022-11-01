package internal

import (
	"net/http"
)

//counterfeiter:generate . HttpGetter
type HttpGetter func(path string) (*http.Response, error)

var HttpGet = http.Get
