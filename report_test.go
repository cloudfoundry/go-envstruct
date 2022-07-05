package envstruct_test

import (
	"bytes"
	"os"

	envstruct "code.cloudfoundry.org/go-envstruct"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Report", func() {
	var (
		ts         SmallTestStruct
		outputText string
	)

	Describe("Report()", func() {
		BeforeEach(func() {
			for k, v := range baseEnvVars {
				os.Setenv(k, v)
			}

			err := envstruct.Load(&ts)
			Expect(err).ToNot(HaveOccurred())

			outputBuffer := bytes.NewBuffer(nil)
			envstruct.ReportWriter = outputBuffer

			err = envstruct.WriteReport(&ts)
			Expect(err).ToNot(HaveOccurred())

			outputText = outputBuffer.String()
		})

		It("prints a report of the given envstruct struct", func() {
			Expect(outputText).To(Equal(expectedReportOutput))
		})
	})
})

const (
	expectedReportOutput = `FIELD NAME:                         TYPE:       ENV:                  REQUIRED:  VALUE:
SmallTestStruct.HiddenThing         string      HIDDEN_THING          false      (OMITTED)
SmallTestStruct.StringThing         string      STRING_THING          false      stringy thingy
SmallTestStruct.BoolThing           bool        BOOL_THING            false      true
SmallTestStruct.IntThing            int         INT_THING             false      100
SmallTestStruct.FloatThing          float64     FLOAT_THING           false      3.14159
SmallTestStruct.ComplexThing        complex128  COMPLEX_THING         false      (3+14159i)
SmallTestStruct.URLThing            *url.URL    URL_THING             false      http://github.com/some/path
SmallTestStruct.StringSliceThing    []string    STRING_SLICE_THING    false      [one two three]
SmallTestStruct.CaseSensitiveThing  string      CASE_SENSITIVE_THING  false      case sensitive
SmallTestSubStruct.SecretThing      string      SECRET_THING          false      (OMITTED)
SmallTestSubStruct.SecretThing      string      SECRET_THING          false      (OMITTED)
`
)
