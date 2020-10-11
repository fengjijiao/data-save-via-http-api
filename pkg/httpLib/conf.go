package httpLib

var (
	listenAddress string = ":8090"
	staticFileBasePath string = "./www"
)

func SetListenAddress(newListenAdress string) {
	listenAddress = newListenAdress
}

func SetStaticFileBasePath(newStaticFileBasePath string) {
	staticFileBasePath = newStaticFileBasePath
}