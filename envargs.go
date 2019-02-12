//	从环境变量中初始化运行参数
package cogutils

import (
	"log"

	"github.com/timest/env"
)

type StarterArgs struct {
	HTTPS struct {
		Addr             string `default:"0.0.0.0:4000"`
		Static_Path      string `default:"./static"`
		Template_Path    string `default:"./templates"`
		Private_Pem_File string `default:"./private.pem"`
		Cert_File        string `default:"./file.crt"`
	}
	WEB_SOCKET struct {
		Enabled bool   `default:"true"`
		Path    string `default:"/ws"`
	}
	WEB struct {
		TEMPLATE_NAME string `default:"gin"`
	}
}

var Args StarterArgs

func init() {
	log.Println("init envargs module.")
	env.Fill(&Args)
	log.Println("EnvArgs is ", Args)
}
