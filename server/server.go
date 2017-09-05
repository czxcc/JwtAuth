package server

import (
    "strconv"
    "time"
    "fmt"
    "github.com/kataras/iris"
    "github.com/BluePecker/JwtAuth/server/router"
    "github.com/Sirupsen/logrus"
    "github.com/kataras/iris/context"
)

type TLS struct {
    Cert string
    Key  string
}

type Options struct {
    Host string
    Port int
    Tls  *TLS
}

type Server struct {
    app *iris.Application
}

func (s *Server) initHttpApp() {
    if s.app == nil {
        s.app = iris.New()
        s.app.Use(func(ctx context.Context) {
            startTime := time.Now()
            ctx.Next()
            endTime := time.Now()
            latency := endTime.Sub(startTime)
            logrus.Infof("%v %4v %s %s %s", strconv.Itoa(ctx.GetStatusCode()), latency, ctx.RemoteAddr(), ctx.Method(), ctx.Path())
        })
    }
}

func (s *Server) Accept(options Options) {
    s.initHttpApp()
    
    if options.Tls != nil {
        //是否将80端口的请求转发到443
        //target, _ := url.Parse("https://127.0.0.1:443")
        //go host.NewProxy("127.0.0.1:80", target).ListenAndServe()
        var addr string = fmt.Sprintf("%s:%d", options.Host, options.Port)
        err := s.app.Run(iris.TLS(addr, options.Tls.Cert, options.Tls.Key))
        if err != nil {
            logrus.Error(err)
        }
        
    } else {
        var addr string = fmt.Sprintf("%s:%d", options.Host, options.Port)
        err := s.app.Run(iris.Addr(addr))
        if err != nil {
            logrus.Error(err)
        }
    }
}

func (s *Server) AddRouter(routers... router.Router) {
    s.initHttpApp()
    
    for _, item := range routers {
        item.Routes(s.app)
    }
}