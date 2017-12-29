package g

// Log配置
var (
	Log *Logger
)

var (
	// 元数据目录
	MetaDir string = "./meta/"

	// 元数据文件名
	MetaName string = "password.db"
)

func init() {
	parseFlags()

	initLog()
}

func parseFlags() {

}

func initLog() {
	Log = NewLogger()
}
