package route

import (
	"encoding/json"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type TmpFile struct {
	Id      int       `json:"id"`
	Name    string    `json:"name"`
	Size    int       `json:"size"`
	At      time.Time `json:"at"`
	content multipart.File
}

type Sizer interface {
	Size() int64
}

var files = []TmpFile{}

func GetFiles(res http.ResponseWriter, req *http.Request) {
	setHeaders(res)
	defer onError(res)
	enc := json.NewEncoder(res)
	enc.Encode(files)
}

func GetFile(res http.ResponseWriter, req *http.Request) {
	setHeaders(res)
	defer onError(res)

	i, err := strconv.Atoi(req.URL.Query()["id"][0])
	if err != nil {
		panic(err)
	}
	for idx := range files {
		if files[idx].Id != i {
			continue
		}
		ns := strings.Split(files[idx].Name, ".")
		res.Header().Set("Content-Type", mime.TypeByExtension(ns[len(ns)-1]))
		res.Header().Set("Content-Disposition", "attachment; filename*=utf-8''"+url.QueryEscape(files[idx].Name))
		http.ServeContent(res, req, files[idx].Name, time.Now(), files[idx].content)
		return
	}
	panic("File not found : ${file.Id}")
}

func HandleFile(res http.ResponseWriter, req *http.Request) {
	setHeaders(res)
	defer onError(res)

	if req.Method == "OPTIONS" {
		res.WriteHeader(http.StatusOK)
		return
	}

	if req.Method == "GET" {
		GetFile(res, req)
		return
	}

	if req.Method == "POST" {
		onPostFile(res, req)
		return
	}

	if req.Method == "DELETE" {
		onDeleteFile(res, req)
		return
	}
}

func onPostFile(res http.ResponseWriter, req *http.Request) {
	f, fileHeader, err := req.FormFile("file")
	if err != nil {
		panic(err)
	}
	var file TmpFile
	file.Name = fileHeader.Filename
	file.Id = getIdx()
	file.content = f
	file.Size = int(f.(Sizer).Size())
	files = append(files, file)
	resp := struct {
		Status  string `json:"status"`
		Id      int    `json:"id"`
		Message string `json:"message"`
	}{
		"OK",
		file.Id,
		"",
	}
	enc := json.NewEncoder(res)
	enc.Encode(resp)
}

func onDeleteFile(res http.ResponseWriter, req *http.Request) {
	var file TmpFile
	json.NewDecoder(req.Body).Decode(&file)
	resp := struct {
		Status  string `json:"status"`
		Id      int    `json:"id"`
		Message string `json:"message"`
	}{
		"OK",
		file.Id,
		"",
	}
	enc := json.NewEncoder(res)

	if file.Id == 0 {
		resp.Status = "ERROR"
		resp.Message = "Invalid Id"
		enc.Encode(resp)
		return
	}

	for idx, each := range files {
		if each.Id != file.Id {
			continue
		}
		files = append(files[:idx], files[idx+1:]...)
		enc.Encode(resp)
		return
	}
}
