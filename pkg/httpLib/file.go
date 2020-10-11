package httpLib

import (
	"io/ioutil"
	"fmt"
	//"os"
	//"net/http"
	"path/filepath"
	"mime"
)

func ReadFile(FilePath string) []byte {
	dat, err := ioutil.ReadFile(staticFileBasePath + "/" + FilePath)
	if(err != nil) {
		return []byte(fmt.Sprintf("ReadFile Error: %s\n", err))
	}else {
		return dat
	}
}

// func GetFileContentType(FilePath string) string {
// 	f, err := os.Open(staticFileBasePath + "/" + FilePath)
// 	if err != nil {
// 		return "error"
// 	}
// 	defer f.Close()
// 	buffer := make([]byte, 256)
// 	_, err = f.Read(buffer)
// 	if err != nil {
// 		return "error"
// 	}
// 	return http.DetectContentType(buffer)
// }

func GetFileContentType(FilePath string) string {
	extension := filepath.Ext(FilePath)
	switch extension {
		case ".css":
			return "text/css"
		case ".htm", ".html":
			return "text/html"
		case ".js":
			return "application/javascript"
		default:
			return mime.TypeByExtension(extension)
	}
}