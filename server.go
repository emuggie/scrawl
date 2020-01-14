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

func getIdx() int {
	if idx >= math.MaxUint32 {
		idx = 0
	}
	idx += 1
	return idx
}

func main() {
	log.Println("Start Server")
	memos := []Memo{}

	var mux = http.NewServeMux()

	mux.HandleFunc("/memos", func(res http.ResponseWriter, req *http.Request) {
		log.Println("Get Memos")
		defer func() {
			err := recover()
			if err == nil {
				return
			}
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte("500 - Something bad happened!"))
		}()

		enc := json.NewEncoder(res)
		res.Header().Set("Content-Type", "application/json")
		res.Header().Set("Access-Control-Allow-Origin", "*")
		enc.Encode(memos)
	})

	mux.HandleFunc("/memo", func(res http.ResponseWriter, req *http.Request) {
		//Catch
		defer func() {
			err := recover()
			if err == nil {
				return
			}
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte("500 - Something bad happened!"))
		}()

		res.Header().Set("Content-Type", "application/json")
		res.Header().Set("Access-Control-Allow-Origin", "*")

		log.Println("Requesting", req.Method)
		if req.Method == "OPTIONS" {
			res.WriteHeader(http.StatusOK)
			return
		}
		var memo Memo
		json.NewDecoder(req.Body).Decode(&memo)

		if req.Method == "POST" {
			fmt.Println("POST", memo.Id, memo.Text)
			if memo.Text == "" {
				fmt.Fprintf(res, "{\"status\":\"error\"}")
				return
			}
			if memo.Id == 0 {
				//new
				memo.Modified = time.Now()
				memo.Id = getIdx()
				memos = append(memos, memo)
				fmt.Fprintf(res, "{\"status\":\"ok\",\"id\":"+strconv.Itoa(memo.Id)+"}")
				return
			}
			for _, each := range memos {
				if each.Id == memo.Id {
					each.Text = memo.Text
					each.Modified = time.Now()
					break
				}
			}
			fmt.Fprintf(res, "{\"status\":\"ok\",\"id\":"+strconv.Itoa(memo.Id)+"}")
			return
		} else if req.Method == "DELETE" {
			fmt.Println("DELETE", memo.Id, memo.Text)
			if memo.Id == 0 {
				//new
				return
			}
			for idx, each := range memos {
				if each.Id == memo.Id {
					memos = append(memos[:idx], memos[idx+1:]...)
					break
				}
			}
			fmt.Fprintf(res, "{\"status\":\"ok\",\"id\":"+strconv.Itoa(memo.Id)+"}")
			return
		}
	})

	// Static resource serving
	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		filePath := strings.Trim(req.URL.Path, "/")
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
