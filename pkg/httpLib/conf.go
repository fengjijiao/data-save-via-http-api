package httpLib

var (
	listenAddress string = ":8090"
	staticFileBasePath string = "./www"
	parseMultipartFormMaxMemory int64 = 64
)

func SetListenAddress(newListenAdress string) {
	listenAddress = newListenAdress
}

func SetStaticFileBasePath(newStaticFileBasePath string) {
	staticFileBasePath = newStaticFileBasePath
}