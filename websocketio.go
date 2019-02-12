package cogutils

import (
	"log"
	"net/http"

	"github.com/googollee/go-socket.io"
)

type cogSocketIo struct {
	wsserver *socketio.Server
}

func NewCogSocketIo() *cogSocketIo {
	var ws = cogSocketIo{}
	if !Args.WEB_SOCKET.Enabled {
		log.Fatalln("web socket disabled.")
		return nil
	}
	var err error
	ws.wsserver, err = socketio.NewServer(nil)
	if err != nil {
		log.Fatalln("open socker io server error.", err)
		return nil
	}
	return &ws
}

func (ws *cogSocketIo) InitWebSocket() {
	http.Handle(Args.WEB_SOCKET.Path, ws.wsserver)
}

func (ws *cogSocketIo) On(event string, f interface{}) {
	err := ws.wsserver.On(event, f)
	if err != nil {
		log.Fatalln("can not register the event on the socketio.", event, f)
	}
}

func init() {
	log.Println("init websocktio module.")
}
