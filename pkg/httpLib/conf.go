package httpLib

var (
	listenAddress = ":8090"
	parseMultipartFormMaxMemory int64 = 64
)

func SetListenAddress(newListenAdress string) {
	listenAddress = newListenAdress
}