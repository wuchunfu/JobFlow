package grpcpool

import (
	"context"
	"gin-vue/config"
	"gin-vue/middleware/rpc/auth"
	pb "gin-vue/middleware/rpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"strings"
	"sync"
	"time"
)

const (
	backOffMaxDelay = 3 * time.Second
	dialTimeout     = 2 * time.Second
)

var (
	Pool = &GRPCPool{
		conns: make(map[string]*Client),
	}

	keepAliveParams = keepalive.ClientParameters{
		Time:                20 * time.Second,
		Timeout:             3 * time.Second,
		PermitWithoutStream: true,
	}

	//connectParams = grpc.ConnectParams{
	//	Backoff: backoff.Config{
	//		MaxDelay: backOffMaxDelay,
	//	},
	//}
	//connectParams = grpc.ConnectParams{
	//	MinConnectTimeout: dialTimeout,
	//}
)

type Client struct {
	conn      *grpc.ClientConn
	rpcClient pb.TaskClient
}

type GRPCPool struct {
	// map key格式 ip:port
	conns map[string]*Client
	mu    sync.RWMutex
}

func (pool *GRPCPool) Get(addr string) (pb.TaskClient, error) {
	pool.mu.RLock()
	client, ok := pool.conns[addr]
	pool.mu.RUnlock()
	if ok {
		return client.rpcClient, nil
	}

	client, err := pool.factory(addr)
	if err != nil {
		return nil, err
	}
	return client.rpcClient, nil
}

// 释放连接
func (pool *GRPCPool) Release(addr string) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	client, ok := pool.conns[addr]
	if !ok {
		return
	}
	delete(pool.conns, addr)
	client.conn.Close()
}

// 创建连接
func (pool *GRPCPool) factory(addr string) (*Client, error) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	client, ok := pool.conns[addr]
	if ok {
		return client, nil
	}

	//opts := []grpc.DialOption{
	//	grpc.WithKeepaliveParams(keepAliveParams),
	//	grpc.WithConnectParams(connectParams),
	//}

	opts := []grpc.DialOption{
		grpc.WithKeepaliveParams(keepAliveParams),
		grpc.WithDisableRetry(),
	}

	setting := config.InitConfig
	if !setting.System.EnableTls {
		opts = append(opts, grpc.WithInsecure(), grpc.WithBlock())
	} else {
		server := strings.Split(addr, ":")
		certificate := auth.Certificate{
			CAFile:     setting.System.CaFile,
			CertFile:   setting.System.CertFile,
			KeyFile:    setting.System.KeyFile,
			ServerName: server[0],
		}

		transportCreds, err := certificate.GetTransportCredsForClient()
		if err != nil {
			return nil, err
		}
		opts = append(opts, grpc.WithTransportCredentials(transportCreds))
	}

	ctx, cancel := context.WithTimeout(context.Background(), dialTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr, opts...)
	if err != nil {
		return nil, err
	}

	client = &Client{
		conn:      conn,
		rpcClient: pb.NewTaskClient(conn),
	}
	pool.conns[addr] = client
	return client, nil
}
