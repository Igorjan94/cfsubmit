package cfsubmit

import (
	"errors"
	"regexp"
	"strings"
)

var ErrUnknownFilename = errors.New("Unknown filename format. Example: 123a.cpp")

var cfSubmissionFileRegex = regexp.MustCompile(`(\d+)(\w+)\.(\w+)`)

type CFSubmissionFile struct {
	ProblemID string
	ContestID string
	Extension string
}

func New(filename string) (*CFSubmissionFile, error) {
	matches := cfSubmissionFileRegex.FindStringSubmatch(filename)
	if len(matches) < 4 {
		return nil, ErrUnknownFilename
	}
	submission := &CFSubmissionFile{
		ContestID: matches[1],
		ProblemID: strings.ToUpper(matches[2]),
		Extension: strings.ToLower(matches[3]),
	}
	return submission, nil
}
