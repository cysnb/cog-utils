//	从初始化gin的web环境，加载static资源和template资源
package cogutils

import (
	// "flag"
	"log"
	// "net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebginHelper struct {
	ZdbcWebRouter *gin.Engine
	initFuncMap   map[string][]func(*gin.RouterGroup)
	upgrader      websocket.Upgrader
}

func NewWebginHelper() *WebginHelper {
	w := WebginHelper{
		ZdbcWebRouter: gin.Default(),
		initFuncMap:   make(map[string][]func(*gin.RouterGroup)),
	}
	if Args.WEB_SOCKET.Enabled {
		w.upgrader = websocket.Upgrader{}
	}
	return &w
}

//	注册web处理函数
//	示例如下：
// func init() {
// 	RegisterInitFunc("/if", initIndexFunc)
// }
//
// func initIndexFunc(router *gin.RouterGroup) {
// 	router.GET("/", func(c *gin.Context) {
// 		c.HTML(http.StatusOK, "index.tmpl", gin.H{
// 			"now": "" + time.Now().String(),
// 		})
// 	})
// }
//
//
// func init() {
//		websocket
// 	RegisterInitFunc("/if", true, initIndexFunc)
// }
//
// func initIndexFunc(router *gin.RouterGroup) {
// 	router.GET("/", func(c *gin.Context) {
//		//升级get请求为webSocket协议
// 		ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
// 		if err != nil {
//			return
//		}
//		defer ws.Close()
//		for {
//		//读取ws中的数据
//		mt, message, err := ws.ReadMessage()
//		if err != nil {
//			break
//		}
//		if string(message) == "ping" {
//			message = []byte("pong")
//		}
//		//写入ws数据
//		err = ws.WriteMessage(mt, message)
//		if err != nil {
//			break
//		}
// 	}
// }
//
func (w *WebginHelper) RegisterInitFunc(path string, f func(*gin.RouterGroup)) {
	w.RegisterInitFunc1(path, false, f)
}

func (w *WebginHelper) RegisterInitFunc1(path string, isWebsocket bool, f func(*gin.RouterGroup)) {
	log.Println("RegisterInitFunc ", path, f)
	var oldList, _ = w.initFuncMap[path]
	oldList = append(oldList, f)
	w.initFuncMap[path] = oldList
}

func init() {
	if strings.EqualFold(Args.WEB.TEMPLATE_NAME, "gin") {
		log.Println("init webgin module.")
	}
}

func (w *WebginHelper) InitHttp() {
	w.ZdbcWebRouter.Static("/ui", Args.HTTPS.Static_Path)
	w.loadAllTemplatesFromPath(Args.HTTPS.Template_Path)
	for path, fList := range w.initFuncMap {
		var router = w.ZdbcWebRouter.Group(path)
		for _, f := range fList {
			f(router)
		}
	}
}

func (w *WebginHelper) StartWebListen() {
	w.ZdbcWebRouter.RunTLS(Args.HTTPS.Addr, Args.HTTPS.Cert_File, Args.HTTPS.Private_Pem_File)
}

func (w *WebginHelper) loadAllTemplatesFromPath(templateBaseDir string) error {
	var f, err = os.OpenFile(templateBaseDir, os.O_RDONLY, os.ModeAppend)
	if err != nil {
		log.Println("Cann't open the file."+templateBaseDir+", err is ", err)
		return err
	}
	var fileInfo os.FileInfo = nil
	fileInfo, err = f.Stat()
	if err != nil {
		log.Println("Cann't get the file stat info of the file.", err)
		return err
	}
	if fileInfo.IsDir() {
		var subDirFileInfoList, err = f.Readdir(-1)
		if err != nil {
			log.Println("Cann't read sub dirs.", err)
			return err
		}
		for _, subdir := range subDirFileInfoList {
			var err = w.loadAllTemplatesFromPath(templateBaseDir + string(os.PathSeparator) + subdir.Name())
			if err != nil {
				return err
			}
		}
	} else if strings.HasSuffix(templateBaseDir, ".tmpl") {
		log.Println("LoadHTMLFiles ", templateBaseDir)
		w.ZdbcWebRouter.LoadHTMLFiles(templateBaseDir)
	}
	return nil
}
