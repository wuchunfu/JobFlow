package client

import (
	"errors"
	"fmt"
	"gin-vue/middleware/rpc/grpcpool"
	pb "gin-vue/middleware/rpc/proto"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
	"time"
)

var (
	taskMap sync.Map
)

func generateTaskUniqueKey(ip string, port int, id int64) string {
	return fmt.Sprintf("%s:%d:%d", ip, port, id)
}

func Stop(ip string, port int, id int64) {
	key := generateTaskUniqueKey(ip, port, id)
	cancel, ok := taskMap.Load(key)
	if !ok {
		return
	}
	cancel.(context.CancelFunc)()
}

func Exec(ip string, port int, taskReq *pb.TaskRequest) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Error("panic#rpc/client.go:Exec#", err)
		}
	}()
	addr := fmt.Sprintf("%s:%d", ip, port)
	// Set up a connection to the server.
	//conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	//if err != nil {
	//	log.Fatalf("did not connect: %v", err)
	//}
	//defer conn.Close()
	//pool := pb.NewTaskClient(conn)

	pool, err := grpcpool.Pool.Get(addr)
	if err != nil {
		return "", errors.New("执行超时, 强制结束")
	}
	if taskReq.Timeout <= 0 || taskReq.Timeout > 86400 {
		taskReq.Timeout = 86400
	}
	timeout := time.Duration(taskReq.Timeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	taskUniqueKey := generateTaskUniqueKey(ip, port, taskReq.Id)
	taskMap.Store(taskUniqueKey, cancel)
	defer taskMap.Delete(taskUniqueKey)

	resp, err := pool.Run(ctx, taskReq)
	if err != nil {
		return parseGRPCError(err)
	}

	if resp.Error == "" {
		return resp.Output, nil
	}
	return resp.Output, errors.New(resp.Error)
}

func parseGRPCError(err error) (string, error) {
	switch status.Code(err) {
	case codes.Unavailable:
		return "", errors.New("无法连接远程服务器")
	case codes.DeadlineExceeded:
		return "", errors.New("执行超时, 强制结束")
	case codes.Canceled:
		return "", errors.New("手动停止")
	}
	return "", err
}
