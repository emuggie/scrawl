package route

import (
	"encoding/json"
	"math"
	"net/http"
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

func GetMemos(res http.ResponseWriter, req *http.Request) {
	setHeaders(res)
	defer onError(res)
	enc := json.NewEncoder(res)
	enc.Encode(memos)
}

func PostMemo(res http.ResponseWriter, req *http.Request) {
	setHeaders(res)
	defer onError(res)

	if req.Method == "OPTIONS" {
		res.WriteHeader(http.StatusOK)
		return
	}

	if req.Method == "POST" {
		onPostMemo(res, req)
		return
	}

	if req.Method == "DELETE" {
		onDeleteMemo(res, req)
		return
	}
}

func onPostMemo(res http.ResponseWriter, req *http.Request) {
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
	for idx := range memos {
		if memos[idx].Id != memo.Id {
			continue
		}
		memos[idx].Text = memo.Text
		memos[idx].Modified = time.Now()
		enc.Encode(resp)
		return
	}
	resp.Status = "ERROR"
	resp.Message = "Not Found"
	enc.Encode(resp)
	return
}

// Delete
func onDeleteMemo(res http.ResponseWriter, req *http.Request) {
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
