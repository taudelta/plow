package plow

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/onsi/ginkgo/v2"
)

type CheckJSONResponse struct {
	ItTitle     string
	RequestFile string
	Expected    interface{}
	Target      interface{}
}

func CheckJsonPOST(path string, runServer func() *httptest.Server, opts CheckJSONResponse) {
	var testServer *httptest.Server
	var resp *http.Response
	var err error
	ginkgo.JustBeforeEach(func() {
		testServer = runServer()
		request, reqErr := ioutil.ReadFile(opts.RequestFile)
		if reqErr != nil {
			log.Fatal(err)
		}
		resp, err = SendPOSTRequest(testServer.URL+path, RequestOptions{
			Headers: http.Header{
				"Content-Type": {"application/json"},
			},
			Body: bytes.NewBuffer(request),
		})
	})
	ginkgo.It(opts.ItTitle, func() {
		ExpectJSON(resp, err, &opts.Expected, &opts.Target)
	})
	ginkgo.AfterEach(func() {
		testServer.Close()
	})
}
