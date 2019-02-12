//	初始化web环境，包括普通http/https环境以及websocket环境
package cogutils

import (
	// "flag"
	"log"
	// "net/http"
	// "os"
	// "strings"
	// "github.com/gin-gonic/gin"
	// "github.com/gorilla/websocket"
)

var wh *WebginHelper
var ws *cogSocketIo

// func (w WebginHelper)initIndexFunc(router *gin.RouterGroup) {
// 	router.GET("/abc", func(c *gin.Context) {
// 		//升级get请求为webSocket协议
// 		ws, err := w.upgrader.Upgrade(c.Writer, c.Request, nil)
// 		if err != nil {
// 			return
// 		}
// 		defer  ws.Close()
// 		for {
// 		//读取ws中的数据
// 		mt, message, err := ws.ReadMessage()
// 		if err != nil {
// 			break
// 		}

// 		if string(message) == "ping" {
// 			message = []byte("pong")
// 		}
// 		//写入ws数据
// 		err = ws.WriteMessage(mt, message)
// 		if err != nil {
// 			break
// 		}
// 	}
// }

func InitWeb() {
	log.Println(".....")
	wh = NewWebginHelper()
	ws = NewCogSocketIo()
	wh.InitHttp()
	ws.InitWebSocket()
}

func GetCogSocketIo() *cogSocketIo {
	return ws
}

func init() {
	log.Println("init webinit module.")
}
