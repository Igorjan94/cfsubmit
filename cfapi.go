package cfsubmit

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	ErrNoSuchUser      = errors.New("No such user. Please check your settings")
	ErrUnknownFilename = errors.New("Unknown filename format. Example: 123a.cpp")
)

var (
	cfSubmissionFileRegex = regexp.MustCompile(`^(\d+)(\w+)\.(\w+)`)
)

const (
	ProblemTypeProgramming = "PROGRAMMING"
	ProblemTypeQuestion    = "QUESTION"
)

const (
	PartTypeContestant       = "CONTESTANT"
	PartTypePractice         = "PRACTICE"
	PartTypeVirtual          = "VIRTUAL"
	PartTypeManager          = "MANAGER"
	PartTypeOutOfCompetition = "OUT_OF_COMPETITION"
)

const (
	VerFailed                  = "FAILED"
	VerOK                      = "OK"
	VerPartial                 = "PARTIAL"
	VerCE                      = "COMPILATION_ERROR"
	VerRE                      = "RUNTIME_ERROR"
	VerWA                      = "WRONG_ANSWER"
	VerPE                      = "PRESENTATION_ERROR"
	VerTLE                     = "TIME_LIMIT_EXCEEDED"
	VerMLE                     = "MEMORY_LIMIT_EXCEEDED"
	VerILE                     = "IDLENESS_LIMIT_EXCEEDED"
	VerSecurityViolated        = "SECURITY_VIOLATED"
	VerCrashed                 = "CRASHED"
	VerInputPreparationCrashed = "INPUT_PREPARATION_CRASHED"
	VerChallenged              = "CHALLENGED"
	VerSkipped                 = "SKIPPED"
	VerTesting                 = "TESTING"
	VerRejected                = "REJECTED"
	VerMissing                 = ""
)

const (
	TestSetSamples    = "SAMPLES"
	TestSetPretests   = "PRETESTS"
	TestSetTests      = "TESTS"
	TestSetChallenges = "CHALLENGES"
	TestSetTests1     = "TESTS1"
	TestSetTests2     = "TESTS2"
	TestSetTests3     = "TESTS3"
	TestSetTests4     = "TESTS4"
	TestSetTests5     = "TESTS5"
	TestSetTests6     = "TESTS6"
	TestSetTests7     = "TESTS7"
	TestSetTests8     = "TESTS8"
	TestSetTests9     = "TESTS9"
	TestSetTests10    = "TESTS10"
)

type MemoryAmount int64

func (m MemoryAmount) String() string {
	if m < 1024 {
		return strconv.FormatInt(int64(m), 10) + "b"
	} else {
		x := float64(m)
		x /= 1024.0
		if x < 1024.0 {
			return strconv.FormatFloat(x, 'f', 3, 64) + "kB"
		}
		x /= 1024.0
		return strconv.FormatFloat(x, 'f', 3, 64) + "MB"
	}
}

type Problem struct {
	ContestID int      `json:"contestId"`
	Index     string   `json:"index"`
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	Points    float32  `json:"points"`
	Tags      []string `json:"tags"`
}

type Member struct {
	Handle string `json:"handle"`
}

type Party struct {
	ContestID       int       `json:"contestId"`
	Members         []Member  `json:"members"`
	ParticipantType string    `json:"participantType"`
	TeamName        string    `json:"teamName"`
	Ghost           bool      `json:"ghost"`
	Room            int       `json:"room"`
	StartTime       EpochTime `json:"startTimeSeconds"`
}

type Submission struct {
	ID                  int           `json:"id"`
	ContestID           int           `json:"contestId"`
	CreationTime        EpochTime     `json:"creationTimeSeconds"`
	RelativeTime        time.Duration `json:"relativeTimeSeconds"`
	Problem             Problem       `json:"problem"`
	Author              Party         `json:"author"`
	ProgrammingLanguage string        `json:"programmingLanguage"`
	Verdict             string        `json:"verdict"`
	TestSet             string        `json:"testset"`
	PassedTestCount     int           `json:"passedTestCount"`
	TimeConsumed        time.Duration `json:"timeConsumedMillis"`
	MemoryConsumed      MemoryAmount  `json:"memoryConsumedBytes"`
}

func NewSubmission(filename string) (*Submission, error) {
	matches := cfSubmissionFileRegex.FindStringSubmatch(filename)
	if len(matches) < 4 {
		return nil, ErrUnknownFilename
	}

	var s Submission
	var err error

	s.ContestID, err = strconv.Atoi(matches[1])
	if err != nil {
		return nil, ErrUnknownFilename
	}
	s.Problem.Index = strings.ToUpper(matches[2])
	s.ProgrammingLanguage = strings.ToLower(matches[3])
	return &s, nil
}

func UserStatus(handle string, count int) ([]Submission, error) {
	var Submissions struct {
		Status string       `json:"status"`
		Result []Submission `json:"result"`
	}
	urlValues := url.Values{
		"handle": []string{handle},
		"from":   []string{"1"},
		"count":  []string{strconv.Itoa(count)},
	}
	resp, err := http.Get(`http://codeforces.com/api/user.status?` + urlValues.Encode())
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(resp.Body).Decode(&Submissions); err != nil {
		return nil, err
	}
	if Submissions.Status != "OK" {
		return nil, ErrNoSuchUser
	}
	return Submissions.Result, nil
}
