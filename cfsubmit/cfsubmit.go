package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/cnt0/cfsubmit"
)

const settingsFileName = "cfsubmit_settings.json"

const (
	TimeLimitSeconds = 60
)

var (
	contestId int
	problemId string
	langId    string
)

var (
	errNoSubmission = errors.New("Submission file not specified")
	errUnkownExt    = errors.New("Unknown extension. Ext must be in lowercase in your settings file")
)

var CFAuthData *cfsubmit.CFSettings

func init() {
	//load settings from json
	var err error
	CFAuthData, err = cfsubmit.ReadSettings()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	if len(os.Args) < 2 {
		fmt.Println(errNoSubmission)
		os.Exit(0)
	}

	submission, err := cfsubmit.NewSubmission(path.Base(os.Args[1]))
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	contestId = submission.ContestID
	problemId = submission.Problem.Index
	if l, ok := CFAuthData.ExtId[submission.ProgrammingLanguage]; !ok {
		fmt.Println(errUnkownExt)
		os.Exit(0)
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
		"/contest/" + strconv.Itoa(contestId) +
		"/problem/" + problemId +
		"?csrf_token=" + CFAuthData.CSRF

	//get request body data; boundary for header
	form, boundary, err := createMultipartForm()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	//ok, construct request
	req, err := http.NewRequest("POST", reqUrl, form)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	//add required headers and cookies
	req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)
	req.AddCookie(&http.Cookie{Name: "X-User", Value: CFAuthData.XUser})

	timeNow := time.Now()
	fmt.Printf("Now: %v\n", timeNow)

	//send request
	if _, err := http.DefaultClient.Do(req); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	fmt.Print("Solution sent.")

	//check results
	if CFAuthData.CheckResults {
		fmt.Print(" Please wait for results...")

		for i := 0; i < TimeLimitSeconds; i++ {
			select {
			case <-time.After(time.Second):

				//check out your last 5 submissions and search the first one with
				//given vontestId and problemId
				submissions, err := cfsubmit.UserStatus(CFAuthData.Handle, 5)
				if err == nil {
					for _, s := range submissions {
						if s.ContestID == contestId && s.Problem.Index == problemId {

							//this was sent before current cfsubmit run; no need to check it
							if timeNow.After(time.Time(s.CreationTime)) {
								fmt.Print(".")
								break
							}

							if s.Verdict != cfsubmit.VerTesting && s.Verdict != cfsubmit.VerMissing {

								//ok, we've got some useful verdict
								fmt.Printf("\nVerdict: %s | Tests passed: %d | Time: %s | Memory: %s\n",
									s.Verdict,
									s.PassedTestCount,
									s.TimeConsumed.String(),
									s.MemoryConsumed.String())
								os.Exit(0)
							} else {

								//still testing, wait...
								fmt.Print(".")
								break
							}
						}
					}
				} else {

					//Error like incorrect username or missing connection; no reason to wait anything
					fmt.Println(err)
					os.Exit(0)
				}
			}
		}
		fmt.Println("\nToo long testing, I'll exit now. Please check results manually\n")
	} else {
		fmt.Println()
	}
}
