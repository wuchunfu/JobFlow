package server

import (
	"github.com/wuchunfu/JobFlow/middleware/rpc/auth"
	pb "github.com/wuchunfu/JobFlow/middleware/rpc/proto"
	"github.com/wuchunfu/JobFlow/utils/unixUtils"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

type Server struct{}

var keepAlivePolicy = keepalive.EnforcementPolicy{
	MinTime:             10 * time.Second,
	PermitWithoutStream: true,
}

var keepAliveParams = keepalive.ServerParameters{
	MaxConnectionIdle: 30 * time.Second,
	Time:              30 * time.Second,
	Timeout:           3 * time.Second,
}

func (s Server) Run(ctx context.Context, req *pb.TaskRequest) (*pb.TaskResponse, error) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Error(err)
		}
	}()
	logrus.Infof("execute cmd start: [id: %d cmd: %s]", req.Id, req.Command)
	output, err := unixUtils.ExecShell(ctx, req.Command)
	resp := new(pb.TaskResponse)
	resp.Output = output
	if err != nil {
		resp.Error = err.Error()
	} else {
		resp.Error = ""
	}
	logrus.Infof("execute cmd end: [id: %d cmd: %s err: %s]", req.Id, req.Command, resp.Error)
	return resp, nil
}

func Start(addr string, enableTLS bool, certificate auth.Certificate) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		logrus.Fatal(err)
	}
	opts := []grpc.ServerOption{
		grpc.KeepaliveParams(keepAliveParams),
		grpc.KeepaliveEnforcementPolicy(keepAlivePolicy),
	}
	if enableTLS {
		tlsConfig, err := certificate.GetTLSConfigForServer()
		if err != nil {
			logrus.Fatal(err)
		}
		opt := grpc.Creds(credentials.NewTLS(tlsConfig))
		opts = append(opts, opt)
	}
	server := grpc.NewServer(opts...)
	pb.RegisterTaskServer(server, Server{})
	logrus.Infof("server listen on %s", addr)

	go func() {
		err = server.Serve(l)
		if err != nil {
			logrus.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	for {
		s := <-c
		logrus.Infoln("收到信号 -- ", s)
		switch s {
		case syscall.SIGHUP:
			logrus.Infoln("收到终端断开信号, 忽略")
		case syscall.SIGINT, syscall.SIGTERM:
			logrus.Info("应用准备退出")
			server.GracefulStop()
			return
		}
	}
}
