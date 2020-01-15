package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Memo struct {
	Id       int    `json:"id"`
	Text     string `json:"text"`
	Modified time.Time
}

var idx = 1
var memos = []Memo{}

func getIdx() int {
	if idx >= math.MaxUint32 {
		idx = 0
	}
	idx += 1
	return idx
}

func setHeaders(res http.ResponseWriter) {
	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "*")
}

func onError(res http.ResponseWriter) {
	err := recover()
	if err == nil {
		return
	}
	fmt.Println(err)
	res.WriteHeader(http.StatusInternalServerError)
	res.Write([]byte("500"))
	return
}

func onPost(res http.ResponseWriter, req *http.Request) {
	var memo Memo
	json.NewDecoder(req.Body).Decode(&memo)

	resp := struct {
		Status  string `json:"status"`
		Id      int    `json:"id"`
		Message string `json:"message"`
	}{
		"OK",
		memo.Id,
		"",
	}
	enc := json.NewEncoder(res)

	if memo.Text == "" {
		resp.Status = "ERROR"
		resp.Message = "Emtpy text"
		enc.Encode(resp)
		return
	}
	// Insert
	if memo.Id == 0 {
		memo.Modified = time.Now()
		memo.Id = getIdx()
		memos = append(memos, memo)
		resp.Id = memo.Id
		enc.Encode(resp)
		return
	}
	// Update
	for _, each := range memos {
		if each.Id != memo.Id {
			continue
		}
		each.Text = memo.Text
		each.Modified = time.Now()
		enc.Encode(resp)
		return
	}
}

// Delete
func onDelete(res http.ResponseWriter, req *http.Request) {
	var memo Memo
	json.NewDecoder(req.Body).Decode(&memo)
	fmt.Println("DELETE", memo.Id, memo.Text)
	resp := struct {
		Status  string `json:"status"`
		Id      int    `json:"id"`
		Message string `json:"message"`
	}{
		"OK",
		memo.Id,
		"",
	}
	enc := json.NewEncoder(res)

	if memo.Id == 0 {
		resp.Status = "ERROR"
		resp.Message = "Invalid Id"
		enc.Encode(resp)
		return
	}

	for idx, each := range memos {
		if each.Id != memo.Id {
			continue
		}
		memos = append(memos[:idx], memos[idx+1:]...)
		enc.Encode(resp)
		return
	}
}

func main() {
	log.Println("Start Server")
	var mux = http.NewServeMux()

	mux.HandleFunc("/memos", func(res http.ResponseWriter, req *http.Request) {
		setHeaders(res)
		defer onError(res)
		enc := json.NewEncoder(res)
		enc.Encode(memos)
	})

	mux.HandleFunc("/memo", func(res http.ResponseWriter, req *http.Request) {
		setHeaders(res)
		defer onError(res)

		if req.Method == "OPTIONS" {
			res.WriteHeader(http.StatusOK)
			return
		}

		if req.Method == "POST" {
			onPost(res, req)
			return
		}

		if req.Method == "DELETE" {
			onDelete(res, req)
			return
		}
	})

	// Static resource serving
	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		filePath := strings.Trim(req.URL.Path, "/")
		// "/ for index.html"
		if filePath == "" {
			filePath = "index.html"
		}
		contents, err := ioutil.ReadFile(filePath)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte("404 - File Not Found"))
			return
		}
		res.Header().Set("Content-Length", strconv.Itoa(len(contents)))
		res.Header().Set("Content-Type", http.DetectContentType(contents))
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

	http.ListenAndServe(":8081", mux)
}
