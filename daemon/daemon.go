package daemon

import (
    "github.com/sevlyar/go-daemon"
    "github.com/Sirupsen/logrus"
    "github.com/BluePecker/JwtAuth/server"
    "github.com/BluePecker/JwtAuth/server/router/jwt"
)

type Conf struct {
    Daemon  bool
    
    PidFile string
    LogFile string
    
    Port    int
}

type Daemon struct {
    
}

func (d *Daemon) Start(conf Conf) {
    if (conf.Daemon == true) {
        dCtx := daemon.Context{
            PidFileName: conf.PidFile,
            PidFilePerm: 0644,
            LogFileName: conf.LogFile,
            LogFilePerm: 0640,
            Umask:       027,
            WorkDir:     "/",
        }
        
        defer dCtx.Release()
        
        if child, err := dCtx.Reborn(); err != nil {
            logrus.Fatal(err)
        } else if child != nil {
            return
        }
    }
    
    api := &server.Server{}
    api.AddRouter(jwt.NewRouter(nil))
    
    api.Accept(server.Options{Host: "", Port: conf.Port})
}