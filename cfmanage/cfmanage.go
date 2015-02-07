package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/cnt0/cfsubmit"
)

const (
	defaultContest  = "" //dont create contest
	defaultCount    = 5
	defaultTemplate = "" //create empty files
)

var (
	archiveFlag  bool
	gzipFlag     bool
	contestFlag  string
	countFlag    int
	templateFlag string
)

func ArchiveSubmissions(dir string) error {
	return filepath.Walk(dir, func(path1 string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if submission, err := cfsubmit.New(info.Name()); err == nil {
				os.Mkdir(submission.ContestID, os.ModeDir|os.ModePerm)
				os.Rename(info.Name(), path.Join(submission.ContestID, info.Name()))
			}
		}
		return nil
	})
}

func ArchiveSubmissionsTGZ(dir string) error {

	buffers := make(map[string]*bytes.Buffer)
	tarWriters := make(map[string]*tar.Writer)

	err := filepath.Walk(dir, func(path1 string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if submission, err := cfsubmit.New(info.Name()); err == nil {

				body, err := ioutil.ReadFile(info.Name())
				if err != nil {
					return err
				}

				buf, ok := buffers[submission.ContestID]
				if !ok {
					buf = bytes.NewBuffer([]byte{})
					buffers[submission.ContestID] = buf
				}
				tw, ok := tarWriters[submission.ContestID]
				if !ok {
					tw = tar.NewWriter(buf)
					tarWriters[submission.ContestID] = tw
				}

				hdr := &tar.Header{
					Name: info.Name(),
					Size: info.Size(),
				}
				if err := tw.WriteHeader(hdr); err != nil {
					return err
				}
				if _, err := tw.Write(body); err != nil {
					return err
				}

				os.Remove(info.Name())
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	for _, tw := range tarWriters {
		if err := tw.Close(); err != nil {
			return err
		}
	}
	for str, buf := range buffers {
		fout, err := os.Create(str + ".tar.gz")
		if err != nil {
			return err
		}
		gw := gzip.NewWriter(fout)
		if _, err := gw.Write(buf.Bytes()); err != nil {
			return err
		}
		if err := gw.Close(); err != nil {
			return err
		}
		if err := fout.Close(); err != nil {
			return err
		}
	}
	return nil
}

func CreateTemplates() error {
	if len(contestFlag) == 0 {
		return nil
	}
	if len(templateFlag) == 0 {
		for i := 0; i < countFlag; i++ {
			fout, err := os.Create(contestFlag + string('A'+byte(i)))
			if err != nil {
				return err
			}
			if err := fout.Close(); err != nil {
				return err
			}
		}
	}

	buf, err := ioutil.ReadFile(templateFlag)
	if err != nil {
		return err
	}

	ext := path.Ext(templateFlag)
	for i := 0; i < countFlag; i++ {
		fout, err := os.Create(contestFlag + string('A'+byte(i)) + ext)
		if err != nil {
			return err
		}
		if _, err := fout.Write(buf); err != nil {
			return err
		}
		if err := fout.Close(); err != nil {
			return err
		}
	}
	return nil
}

func init() {

	flag.BoolVar(&archiveFlag, "a", false, "arhive old submissions into one folder per contest; dominates -z flag")
	flag.BoolVar(&gzipFlag, "z", false, "arhive old submissions into one gzip file per contest")
	flag.StringVar(&contestFlag, "c", defaultContest, "create empty templates for contest; existing files will be rewritten")
	flag.IntVar(&countFlag, "cnt", defaultCount, "how many templates will be created (at most 26)")
	flag.StringVar(&templateFlag, "t", defaultTemplate, "which file will be used as base template")

	flag.Parse()
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	if archiveFlag {
		if err := ArchiveSubmissions(dir); err != nil {
			fmt.Println(err)
		}
	} else if gzipFlag {
		if err := ArchiveSubmissionsTGZ(dir); err != nil {
			fmt.Println(err)
		}
	}
	if err := CreateTemplates(); err != nil {
		fmt.Println(err)
	}
}
