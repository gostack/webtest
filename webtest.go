package webtest

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type RequestTestCase struct {
	Method         string
	URL            string
	Body           io.Reader
	Header         http.Header
	ResponseStatus int
	ResponseBody   io.Reader
}

type RequestTestSuite []RequestTestCase

func (ts RequestTestSuite) Execute(c *http.Client, errFn func(err error)) {
	for i, tc := range ts {
		if tc.Method == "" || tc.URL == "" {
			panic("httpTestCase requires Method and URL to be set")
		}

		if tc.ResponseStatus == 0 {
			tc.ResponseStatus = http.StatusOK
		}

		req, err := http.NewRequest(tc.Method, tc.URL, tc.Body)
		if err != nil {
			panic(err)
		}

		if tc.Header != nil {
			req.Header = tc.Header
		}

		res, err := c.Do(req)
		if err != nil {
			panic(err)
		}

		if res.StatusCode != tc.ResponseStatus {
			errFn(fmt.Errorf("Expected test case #%d (%s %s) to respond with %d but got %d", i, tc.Method, tc.URL, tc.ResponseStatus, res.StatusCode))
		}

		if tc.ResponseBody != nil {
			eb, err := ioutil.ReadAll(tc.ResponseBody)
			if err != nil {
				panic(err)
			}

			ab, err := ioutil.ReadAll(res.Body)
			if err != nil {
				panic(err)
			}

			if strings.TrimSpace(string(ab)) != strings.TrimSpace(string(eb)) {
				errFn(fmt.Errorf("Expected test case #%d (%s %s) response to be:\n%s\nbut got:\n%s", i, tc.Method, tc.URL, eb, ab))
			}
		}
	}
}
