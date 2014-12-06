package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"
	"unicode"
)

var (
	contestId string
	problemId string
	langId    string
)

var CFAuthData struct {
	XUser    string            `json:"X-User"`
	CSRF     string            `json:"CSRF-Token"`
	ExtId    map[string]string `json:"Ext-ID"`
	CFDomain string            `json:"CF-Domain"`
}

func init() {
	//load settings from json
	if jsonData, err := os.Open("cfsubmit_settings.json"); err != nil {
		log.Fatal(err)
	} else {
		if err := json.NewDecoder(jsonData).Decode(&CFAuthData); err != nil {
			log.Fatal(err)
		}
	}

	if len(os.Args) < 2 {
		log.Fatal("Submission file not specified")
	}

	//parse lang id
	if ext := path.Ext(os.Args[1]); len(ext) == 0 {
		log.Fatal("Unknown extension")
	} else {
		id, ok := CFAuthData.ExtId[strings.ToLower(ext[1:])]
		if !ok {
			log.Fatal("Unknown extension")
		}
		langId = id
	}

	//parse contest id & problem id
	cId, pId := []rune{}, []rune{}

	var idx int
	var c rune
	for idx, c = range path.Base(os.Args[1]) {
		if unicode.IsDigit(c) {
			cId = append(cId, c)
		} else {
			break
		}
	}
	for _, c = range path.Base(os.Args[1])[idx:] {
		if c != '.' {
			pId = append(pId, c)
		} else {
			break
		}
	}

	problemId = strings.ToUpper(string(pId))
	contestId = strings.ToUpper(string(cId))

	if len(problemId) == 0 || len(contestId) == 0 {
		log.Fatal("Unknown submission filename format")
	}
}

func main() {

	//prepare multipart form for http request
	var b bytes.Buffer
	mw := multipart.NewWriter(&b)

	if err := mw.WriteField("csrf_token", CFAuthData.CSRF); err != nil {
		log.Fatal(err)
	}
	if err := mw.WriteField("action", "submitSolutionFormSubmitted"); err != nil {
		log.Fatal(err)
	}
	if err := mw.WriteField("submittedProblemIndex", problemId); err != nil {
		log.Fatal(err)
	}
	if err := mw.WriteField("programTypeId", langId); err != nil {
		log.Fatal(err)
	}
	if err := mw.WriteField("sourceFile", ""); err != nil {
		log.Fatal(err)
	}
	if err := mw.WriteField("_tta", "222"); err != nil {
		log.Fatal(err)
	}

	if sol, err := ioutil.ReadFile(os.Args[1]); err != nil {
		log.Fatal(err)
	} else {
		if err := mw.WriteField("source", string(sol)); err != nil {
			log.Fatal(err)
		}
	}
	if err := mw.Close(); err != nil {
		log.Fatal(err)
	}

	//request url
	reqUrl := "http://codeforces." + CFAuthData.CFDomain +
		"/contest/" + contestId +
		"/problem/" + problemId +
		"?csrf_token=" + CFAuthData.CSRF

	req, err := http.NewRequest("POST", reqUrl, &b)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "multipart/form-data; boundary="+mw.Boundary())
	req.AddCookie(&http.Cookie{Name: "X-User", Value: CFAuthData.XUser})

	//send request
	if _, err := http.DefaultClient.Do(req); err != nil {
		log.Fatal(err)
	}

	//maube success
	log.Println("Solution sent. Check result in CF website")
}
