package main

import (
	"flag"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"strconv"
	"strings"

	"github.com/emuggie/scrawl/route"
)

/*static content path*/
const RESOURCE_PATH = "static"

func main() {
	var mux = http.NewServeMux()

	mux.HandleFunc("/memos", route.GetMemos)
	mux.HandleFunc("/memo", route.PostMemo)
	mux.HandleFunc("/files", route.GetFiles)
	mux.HandleFunc("/file", route.HandleFile)

	// Static resource serving
	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		filePath := strings.Trim(req.URL.Path, "/")
		// "/ for index.html"
		if filePath == "" {
			filePath = "index.html"
		}
		contents, err := ioutil.ReadFile(RESOURCE_PATH + "/" + filePath)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte("404 - File Not Found"))
			return
		}
		ns := strings.Split(filePath, ".")
		res.Header().Set("Content-Length", strconv.Itoa(len(contents)))
		res.Header().Set("Content-Type", mime.TypeByExtension(ns[len(ns)-1]))
		res.Write(contents)
		return

		// if val, ok := resourceMap[strings.Trim(req.URL.Path, "/")]; ok {
		// 	data, err := base64.StdEncoding.DecodeString(val)
		// 	if err != nil {
		// 		return
		// 	}
		// 	res.Header().Set("Content-Length", strconv.Itoa(len(data)))
		// 	res.Header().Set("Content-Type", http.DetectContentType(data))
		// 	res.Write(data)
		// 	return
		// }
	})
	port := flag.Int("port", 8080, " port")
	host := flag.String("host", "0.0.0.0", "host")
	flag.Parse()
	log.Println("Start Server", *host+":"+strconv.Itoa(*port))
	http.ListenAndServe(*host+":"+strconv.Itoa(*port), mux)
}
