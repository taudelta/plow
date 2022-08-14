package examples

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/taudelta/plow"
)

type Todo struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Datetime string `json:"datetime"`
}

var todoList = []Todo{
	{ID: 1, Title: "Breakfast", Datetime: "2000-01-01 08:00:00"},
	{ID: 2, Title: "Go to work", Datetime: "2000-01-01 09:00:00"},
}

var newTodo = Todo{
	ID:       3,
	Title:    "Dinner",
	Datetime: "2000-01-01 13:00:00",
}

func RunTestServer() *httptest.Server {
	server := http.NewServeMux()

	server.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			json.NewEncoder(w).Encode(&todoList)
		} else if r.Method == "POST" {
			json.NewEncoder(w).Encode(&newTodo)
		}
	})

	return httptest.NewServer(server)
}

var _ = BeforeSuite(func() {
	output := bytes.NewBuffer([]byte{})
	plow.RunCmd(time.Second, output, "echo before")
})

var _ = AfterSuite(func() {
	output := bytes.NewBuffer([]byte{})
	plow.RunCmd(time.Second, output, "echo after")
	output.Reset()
})

var _ = Describe("spec", func() {
	success := true
	When("run test", func() {
		It("success", func() {
			Expect(success).To(BeTrue())
		})
	})
	When("send GET request to the http microservice", func() {
		var testServer *httptest.Server
		var resp *http.Response
		var err error
		JustBeforeEach(func() {
			testServer = RunTestServer()
			resp, err = plow.SendGETRequest(testServer.URL+"/todos", plow.RequestOptions{})
		})
		It("must send a response", func() {
			var todos []Todo
			plow.ExpectJSON(resp, err, &todoList, &todos)
		})
		AfterEach(func() {
			testServer.Close()
		})
	})

	When("send POST request to the http microservice", func() {
		var todo Todo
		plow.CheckJsonPOST("/todos", RunTestServer, plow.CheckJSONResponse{
			RequestFile: "./request.json",
			Expected:    &newTodo,
			Target:      &todo,
		})
	})
	When("run test with postgres fixtures", func() {
		BeforeEach(func() {
			plow.UsePostgresDB(&plow.DbConfig{
				DSN: "postgres://todos:1@localhost:5432/todos",
			})
			plow.LoadPostgresFixtures("./fixtures", []string{
				"empty_todos.yaml",
				"todos.yaml",
			})
		})
		It("success", func() {
			Expect(success).To(BeTrue())
		})
	})
})
