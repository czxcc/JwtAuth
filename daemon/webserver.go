package daemon

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/netutil"
	"github.com/BluePecker/JwtAuth/daemon/webserver"
	"github.com/BluePecker/JwtAuth/dialog/server/router/signal"
	"github.com/BluePecker/JwtAuth/dialog/server/router/coder"
	"github.com/Sirupsen/logrus"
)

func (d *Daemon) WebServer(ch chan struct{}) error {
	go func() {
		l, err := netutil.UNIX(d.Options.SockFile, 0666)
		if err == nil {
			(&webserver.Web{}).Listen(iris.Listener(l), nil,
				signal.NewRoute(d),
			)
		}
		logrus.Error(err)
	}()
	addr := fmt.Sprintf("%s:%d", d.Options.Host, d.Options.Port)
	runner := iris.Addr(addr)
	if d.Options.TLS.Key != "" && d.Options.TLS.Cert != "" {
		runner = iris.TLS(addr, d.Options.TLS.Cert, d.Options.TLS.Key)
	}
	return (&webserver.Web{}).Listen(runner, ch,
		coder.NewRoute(d),
	)
}
