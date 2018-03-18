package gokvadmin

import (
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"html/template"
	"log"
	"strconv"
	"errors"
)

var errMarshalError = errors.New("Error marshalling response")
var errBadLogin = errors.New("Wrong login or password")

func (gkv *GoKVAdmin) HomeHandler(w http.ResponseWriter, r *http.Request) {
	t := filepath.Join(assetsPath, "index.html")
	tmpl, err := template.ParseFiles(t)
	if err != nil {
		log.Printf("Error pasing template: %v", err.Error())
	}

	tmpl.ExecuteTemplate(w, "index.html", nil)
}

func (gkv *GoKVAdmin) Auth(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Login    string
		Password string
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errResponse(w, fmt.Sprintf("Can not read request body. %v", err.Error()))
		return
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		errResponse(w, fmt.Sprintf("Can not unmarshal request body. %v", err.Error()))
		return
	}

	if (gkv.Config.Auth == nil || (req.Login == gkv.Config.Auth.Login && req.Password == gkv.Config.Auth.Password)) {
		token, err := gkv.Config.Auth.GenerateToken()

		if err != nil {
			errResponse(w, err.Error())
			return
		}

		response, _ := json.Marshal(map[string]interface{}{
			"token" : token,
		})
		sendResponse(w, response)
		return
	}

	errResponse(w, errBadLogin.Error())
}

func (gkv *GoKVAdmin) Init(w http.ResponseWriter, r *http.Request) {
	response, _ := json.Marshal(map[string]interface{}{
		"engine" : gkv.Engine.GetName(),
	})

	sendResponse(w, response)
}

func (gkv *GoKVAdmin) SetKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	bucket := getBucket(r)

	var req struct {
		Value string `json:"value"`
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errResponse(w, fmt.Sprintf("Can not read request body. %v", err.Error()))
		return
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		errResponse(w, fmt.Sprintf("Can not unmarshal request body. %v", err.Error()))
		return
	}

	err = gkv.Engine.Set([]byte(key), []byte(req.Value), bucket)
	if err != nil {
		errResponse(w, fmt.Sprintf("Error setting the key. %v", err.Error()))
		return
	}

	response, err := json.Marshal(map[string]interface{}{
		"key" : key,
		"bucket" : string(bucket),
		"value" : req.Value,
		"exists": true,
	})

	if err != nil {
		errResponse(w, errMarshalError.Error());
		return
	}

	sendResponse(w, response)
}

func (gkv *GoKVAdmin) GetKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	bucket := getBucket(r)

	value, err := gkv.Engine.Get([]byte(key), bucket)
	if err != nil {
		errResponse(w, err.Error());
		return
	}

	response, err := json.Marshal(map[string]interface{}{
		"key" : key,
		"bucket" : string(bucket),
		"value" : string(value),
		"exists": value != nil,
	})

	if err != nil {
		errResponse(w, errMarshalError.Error());
		return
	}

	sendResponse(w, response)
}

func (gkv *GoKVAdmin) GetKeyByPrefix(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	bucket := getBucket(r)
	limit, offset := getGetKeyByPrefixParams(r)

	records, err := gkv.Engine.GetByPrefix([]byte(key), bucket, limit, offset)
	if err != nil {
		errResponse(w, err.Error());
		return
	}

	var results []map[string]interface{}
	for _, record := range records {
		results = append(results, map[string]interface{}{
			"key" : string(record.Key),
			"bucket" : string(bucket),
			"value" : string(record.Value),
			"exists": record.Value != nil,
		})
	}

	response, err := json.Marshal(results)
	if err != nil {
		errResponse(w, errMarshalError.Error());
		return
	}

	sendResponse(w, response)
}

func getBucket(r *http.Request) []byte {
	buckets, ok := r.URL.Query()["bucket"]
	if ok && len(buckets) > 0 {
		return []byte(buckets[0])
	}

	return []byte("")
}

func getGetKeyByPrefixParams(r *http.Request) (int, int) {
	var defaultLimit, defaultOffset = 100, 0
	var limit, offset string
	query := r.URL.Query()
	limits, ok := query["limit"]
	if !ok || len(limits) < 1 {
		limit = "100"
	} else {
		limit = string(limits[0])
	}

	offsets, ok := query["offset"]
	if !ok || len(offsets) < 1 {
		offset = "0"
	} else {
		offset = string(offsets[0])
	}

	iLimit, err := strconv.Atoi(limit)
	if err != nil {
		iLimit = defaultLimit
	}

	iOffset, err := strconv.Atoi(offset)
	if err != nil {
		iOffset = defaultOffset
	}

	return iLimit, iOffset
}

func (gkv *GoKVAdmin) DeleteKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	bucket := getBucket(r)

	err := gkv.Engine.Delete([]byte(key), []byte(bucket))

	if err != nil {
		errResponse(w, err.Error());
		return
	}

	response, _ := json.Marshal(map[string]interface{}{
		"success" : true,

	})

	sendResponse(w, response)
}

func sendResponse(w http.ResponseWriter, message []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(message)
}

func errResponse(w http.ResponseWriter, error string) {
	response, _ := json.Marshal(map[string]interface{}{
		"error" : error,
	})
	sendResponse(w, response)
}