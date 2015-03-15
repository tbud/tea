package tea

import (
	"github.com/tbud/tea/context"
	"github.com/tbud/x/log"
)

var App *context.AppConfigStruct

var Log *log.Logger

func init() {
	App = context.App
	Log = context.Log
}

func Run(port int) {
	context.Run(port)
}
