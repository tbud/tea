package tea

import (
	"github.com/tbud/tea/context"
	"github.com/tbud/x/log"
)

type AppConfigStruct struct {
	Name   string
	Secret string

	HttpPort    int
	HttpAddr    string
	HttpSsl     bool
	HttpSslCert string
	HttpSslKey  string

	CookiePrefix   string
	CookieHttpOnly bool
	CookieSecure   bool
}

var App = &AppConfigStruct{
	HttpPort: 9000,
	HttpAddr: "127.0.0.1",
	HttpSsl:  false,
}

var Log *log.Logger

func init() {
	conf, err := config.Load("conf/app.conf")
	if err != nil {
		panic(err)
	}

	Log, err = log.New(conf.SubConfig("log"))
	if err != nil {
		panic(err)
	}

	if err = conf.SubConfig("app").SetStruct(App); err != nil {
		Log.Error("%v", err)
	}
}

func Run(port int) {
	context.Run(port)
}
