package websocketController

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/hpcloud/tail"
	"net/http"
	"time"
)

var LogFilepath = "./logs/access.log"

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 接口：WebSocket 请求，持续返回日志内容
func LogTail(ctx *gin.Context) {
	fmt.Println("=================")
	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer ws.Close()

	// 监控日志文件
	tails, err := tail.TailFile(LogFilepath, tail.Config{
		ReOpen: true,
		Follow: true,
		// Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "日志文件打开失败",
			"reason":  err.Error(),
		})
		return
	}

	for {
		// 读取 ws 中的数据，接收到客户端链接
		mt, message, err := ws.ReadMessage()
		if err != nil {
			// 客户端关闭连接时也会进入
			fmt.Println(err.Error())
			break
		}
		fmt.Println(string(message))

		// 开始
		var msg *tail.Line
		var ok bool
		for {
			msg, ok = <-tails.Lines
			if !ok {
				fmt.Printf("tail file close reopen, filename:%s\n", tails.Filename)
				time.Sleep(500 * time.Millisecond)
				continue
			}
			fmt.Println("msg:", msg)

			// 写入 websocket 数据
			err = ws.WriteMessage(mt, []byte(msg.Text))
			if err != nil {
				fmt.Println(err.Error())
				break
			}
		}
	}
}
