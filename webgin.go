//	从初始化gin的web环境，加载static资源和template资源
package owrvsutils

import (
	// "flag"
	"log"
	// "net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var ZdbcWebRouter *gin.Engine = gin.Default()

var initFuncMap = make(map[string][]func(*gin.RouterGroup))

//	注册web处理函数
//	示例如下：
// func init() {
// 	RegisterInitFunc("/if", initIndexFunc)
// }

// func initIndexFunc(router *gin.RouterGroup) {
// 	router.GET("/", func(c *gin.Context) {
// 		c.HTML(http.StatusOK, "index.tmpl", gin.H{
// 			"now": "" + time.Now().String(),
// 		})
// 	})
// }
func RegisterInitFunc(path string, f func(*gin.RouterGroup)) {
	log.Println("RegisterInitFunc ", path, f)
	var oldList, _ = initFuncMap[path]
	oldList = append(oldList, f)
	initFuncMap[path] = oldList
}

func init() {
	if strings.EqualFold(Args.WEB.TEMPLATE_NAME, "gin") {
		log.Println("init gin web.")
	}
}

func InitWeb() {
	ZdbcWebRouter.Static("/ui", Args.HTTPS.Static_Path)
	loadAllTemplatesFromPath(Args.HTTPS.Template_Path)
	for path, fList := range initFuncMap {
		var router = ZdbcWebRouter.Group(path)
		for _, f := range fList {
			f(router)
		}
	}
}

func StartWebListen() {
	ZdbcWebRouter.RunTLS(Args.HTTPS.Addr, Args.HTTPS.Cert_File, Args.HTTPS.Private_Pem_File)
}

func loadAllTemplatesFromPath(templateBaseDir string) error {
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
			var err = loadAllTemplatesFromPath(templateBaseDir + string(os.PathSeparator) + subdir.Name())
			if err != nil {
				return err
			}
		}
	} else if strings.HasSuffix(templateBaseDir, ".tmpl") {
		log.Println("LoadHTMLFiles ", templateBaseDir)
		ZdbcWebRouter.LoadHTMLFiles(templateBaseDir)
	}
	return nil
}
