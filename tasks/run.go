package tasks

import (
	"fmt"
	. "github.com/tbud/bud/context"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
)

type runTask struct {
	ServerHost string
	Port       int
	Protocal   string

	proxy *httputil.ReverseProxy
}

func (r *runTask) Execute() error {
	go func() {
		if err := http.ListenAndServe(r.ServerHost, r); err != nil {
			Log.Error("%v", err)
			os.Exit(1)
		}
	}()

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch

	return nil
}

func (r *runTask) Validate() (err error) {
	if r.Port == 0 {
		r.Port, err = getFreePort()
		if err != nil {
			return
		}
	}

	if len(r.ServerHost) == 0 || r.proxy == nil {
		var serverUrl *url.URL
		serverUrl, err = url.ParseRequestURI(fmt.Sprintf(r.Protocal+"://%s:%d", "localhost", r.Port))
		if err != nil {
			return
		}
		r.ServerHost = serverUrl.String()[len(r.Protocal+"://"):]

		r.proxy = httputil.NewSingleHostReverseProxy(serverUrl)
	}

	return nil
}

func (r *runTask) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	r.proxy.ServeHTTP(rw, req)
}

func init() {
	run := runTask{
		Protocal: "http",
	}

	Task("run", &run, Group("tea"), Usage("Run tea framework application."))
}

func getFreePort() (port int, err error) {
	var conn net.Listener
	conn, err = net.Listen("tcp", ":0")
	if err != nil {
		return
	}

	port = conn.Addr().(*net.TCPAddr).Port
	err = conn.Close()
	return
}
