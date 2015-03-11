package tasks

import (
	"crypto/tls"
	"fmt"
	. "github.com/tbud/bud/context"
	"github.com/tbud/tea"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
)

type runTask struct {
	commonCfg
	ServerHost string
	Port       int

	proxy *httputil.ReverseProxy
	app   *App
}

func (r *runTask) Execute() error {
	go func() {
		addr := fmt.Sprintf("%s:%d", tea.App.HttpAddr, tea.App.HttpPort)
		Log.Info("Listening on %s", addr)

		var err error
		if tea.App.HttpSsl {
			err = http.ListenAndServeTLS(addr, tea.App.HttpSslCert, tea.App.HttpSslKey, r)
		} else {
			err = http.ListenAndServe(addr, r)
		}

		if err != nil {
			Log.Fatal("Failed to start reverse proxy: %v", err)
		}
	}()

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch

	return nil
}

func (r *runTask) Validate() (err error) {
	if err = r.commonCfg.Validate(); err != nil {
		return err
	}

	addr := tea.App.HttpAddr
	if len(addr) == 0 {
		addr = "localhost"
	}

	if r.Port == 0 {
		r.Port, err = getFreePort()
		if err != nil {
			return
		}
	}

	scheme := "http"
	if tea.App.HttpSsl {
		scheme = "https"
	}

	if len(r.ServerHost) == 0 || r.proxy == nil {
		var serverUrl *url.URL
		serverUrl, err = url.ParseRequestURI(fmt.Sprintf(scheme+"://%s:%d", addr, r.Port))
		if err != nil {
			return
		}

		r.ServerHost = serverUrl.String()[len(scheme+"://"):]

		r.proxy = httputil.NewSingleHostReverseProxy(serverUrl)

		if tea.App.HttpSsl {
			r.proxy.Transport = &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
		}
	}

	return nil
}

func (r *runTask) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/favicon.ico" {
		return
	}

	if r.app == nil {
		r.app = NewApp(r.binPath, r.Port)
		if err := r.app.Start(); err != nil {
			Log.Error("%v", err)
		}
	}

	r.proxy.ServeHTTP(rw, req)
}

func init() {
	run := runTask{}

	Task("run", TEA_TASK_GROUP, &run, Tasks("tea.build"), Usage("Run tea framework application."))
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
