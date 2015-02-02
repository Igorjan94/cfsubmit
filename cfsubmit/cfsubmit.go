package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"

	"github.com/cnt0/cfsubmit"
)

const settingsFileName = "cfsubmit_settings.json"

var (
	contestId string
	problemId string
	langId    string
)

var (
	errNoSubmission = errors.New("Submission file not specified")
	errUnkownExt    = errors.New("Unknown extension. Ext must be in lowercase in your settings file")
)

var CFAuthData struct {
	XUser    string            `json:"X-User"`
	CSRF     string            `json:"CSRF-Token"`
	ExtId    map[string]string `json:"Ext-ID"`
	CFDomain string            `json:"CF-Domain"`
}

func init() {
	//load settings from json
	if jsonData, err := os.Open(settingsFileName); err != nil {
		log.Fatal(err)
	} else {
		if err := json.NewDecoder(jsonData).Decode(&CFAuthData); err != nil {
			log.Fatal(err)
		}
	}

	if len(os.Args) < 2 {
		log.Fatal(errNoSubmission)
	}

	submission, err := cfsubmit.New(path.Base(os.Args[1]))
	if err != nil {
		log.Fatal(err)
	}

	contestId = submission.ContestID
	problemId = submission.ProblemID
	if l, ok := CFAuthData.ExtId[submission.Extension]; !ok {
		log.Fatal(errUnkownExt)
	} else {
		langId = l
	}
}

func createMultipartForm() (io.Reader, string, error) {
	solutionText, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		return nil, "", err
	}

	//multipart form field: name - value
	formFields := [][]string{
		[]string{"csrf_token", CFAuthData.CSRF},
		[]string{"action", "submitSolutionFormSubmitted"},
		[]string{"submittedProblemIndex", problemId},
		[]string{"programTypeId", langId},
		[]string{"sourceFile", ""},
		[]string{"_tta", "222"},
		[]string{"source", string(solutionText)},
	}

	//cause bytes.Buffer implements both io.Reader and io.Writer
	var b bytes.Buffer
	formWriter := multipart.NewWriter(&b)

	for _, field := range formFields {
		if err := formWriter.WriteField(field[0], field[1]); err != nil {
			return nil, "", err
		}
	}

	if err := formWriter.Close(); err != nil {
		return nil, "", err
	}

	return &b, formWriter.Boundary(), nil
}

func main() {
	//request url
	reqUrl := "http://codeforces." + CFAuthData.CFDomain +
		"/contest/" + contestId +
		"/problem/" + problemId +
		"?csrf_token=" + CFAuthData.CSRF

	//get request body data; boundary for header
	form, boundary, err := createMultipartForm()
	if err != nil {
		log.Fatal(err)
	}

	//ok, construct request
	req, err := http.NewRequest("POST", reqUrl, form)
	if err != nil {
		log.Fatal(err)
	}

	//add required headers and cookies
	req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)
	req.AddCookie(&http.Cookie{Name: "X-User", Value: CFAuthData.XUser})

	//send request
	if _, err := http.DefaultClient.Do(req); err != nil {
		log.Fatal(err)
	}

	//maybe success
	log.Println("Solution sent.")
}
