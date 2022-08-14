package examples

import (
	"bytes"
	"log"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/taudelta/plow"
)

var _ = Describe("spec", func() {
	BeforeEach(func() {
		plow.EnableScriptLogger()
		output := bytes.NewBuffer([]byte{})
		plow.RunCmd(time.Second, output, "echo test1")
		log.Println("before each run shell command 1: ", output.String())
		output.Reset()
		plow.RunCmd(time.Second, output, "date")
		log.Println("before each run shell command 2: ", output.String())
	})
	success := true
	When("run test", func() {
		It("success", func() {
			Expect(success).To(BeTrue())
		})
	})
})
