package envstruct_test

import (
	"crypto/tls"
	"net/url"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

var (
	baseEnvVars = map[string]string{
		"STRING_THING":            "stringy thingy",
		"REQUIRED_THING_A":        "im so required",
		"REQUIRED_THING_B":        "im so required",
		"BOOL_THING":              "true",
		"INT_THING":               "100",
		"INT8_THING":              "20",
		"INT16_THING":             "2000",
		"INT32_THING":             "200000",
		"INT64_THING":             "200000000",
		"UINT_THING":              "100",
		"UINT8_THING":             "20",
		"UINT16_THING":            "2000",
		"UINT32_THING":            "200000",
		"UINT64_THING":            "200000000",
		"FLOAT_THING":             "3.14159",
		"FLOAT32_THING":           "1.2345",
		"FLOAT64_THING":           "9.8765",
		"COMPLEX_THING":           "(3+14159i)",
		"COMPLEX64_THING":         "(1+2345i)",
		"COMPLEX128_THING":        "(9+8765i)",
		"POINTER_TO_STRING":       "pointy stringy thingy",
		"POINTER_TO_BOOL":         "true",
		"POINTER_TO_INT":          "20",
		"POINTER_TO_UINT":         "20",
		"STRING_SLICE_THING":      "one,two,three",
		"INT_SLICE_THING":         "1,2,3",
		"MAP_STRING_STRING_THING": "key_one:value_one,key_two:value_two:with_colon",
		"MAP_INT_STRING_THING":    "1:value_one,2:value_two:with_colon",
		"DURATION_THING":          "2s",
		"URL_THING":               "http://github.com/some/path",
		"UNMARSHALLER_POINTER":    "pointer",
		"UNMARSHALLER_VALUE":      "value",
		"SUB_THING_A":             "sub-string-a",
		"SUB_THING_B":             "200",
		"CaSe_SeNsItIvE_ThInG":    "case sensitive",
	}
)

type LargeTestStruct struct {
	NonEnvThing        string
	DefaultThing       string `env:"DEFAULT_THING"`
	StringThing        string `env:"STRING_THING"`
	RequiredThingA     string `env:"REQUIRED_THING_A,noreport,required"`
	RequiredThingB     string `env:"REQUIRED_THING_B,noreport,required"`
	CaseSensitiveThing string `env:"CaSe_SeNsItIvE_ThInG"`

	BoolThing bool `env:"BOOL_THING"`

	IntThing    int    `env:"INT_THING"`
	Int8Thing   int8   `env:"INT8_THING"`
	Int16Thing  int16  `env:"INT16_THING"`
	Int32Thing  int32  `env:"INT32_THING"`
	Int64Thing  int64  `env:"INT64_THING"`
	UintThing   uint   `env:"UINT_THING"`
	Uint8Thing  uint8  `env:"UINT8_THING"`
	Uint16Thing uint16 `env:"UINT16_THING"`
	Uint32Thing uint32 `env:"UINT32_THING"`
	Uint64Thing uint64 `env:"UINT64_THING"`

	Float32Thing float32 `env:"FLOAT32_THING"`
	Float64Thing float64 `env:"FLOAT64_THING"`

	Complex64Thing  complex64  `env:"COMPLEX64_THING"`
	Complex128Thing complex128 `env:"COMPLEX128_THING"`

	PtrToString *string `env:"POINTER_TO_STRING"`
	PtrToBool   *bool   `env:"POINTER_TO_BOOL"`
	PtrToInt    *int    `env:"POINTER_TO_INT"`
	PtrToUint   *uint   `env:"POINTER_TO_UINT"`

	StringSliceThing []string `env:"STRING_SLICE_THING"`
	IntSliceThing    []int    `env:"INT_SLICE_THING"`

	MapStringStringThing map[string]string `env:"MAP_STRING_STRING_THING"`
	MapIntStringThing    map[int]string    `env:"MAP_INT_STRING_THING"`

	DurationThing time.Duration `env:"DURATION_THING"`
	URLThing      *url.URL      `env:"URL_THING"`

	SubStruct        SubTestStruct
	SubPointerStruct *SubTestStruct

	UnmarshallerPointer *spyUnmarshaller `env:"UNMARSHALLER_POINTER"`
	UnmarshallerValue   spyUnmarshaller  `env:"UNMARSHALLER_VALUE"`

	TLSConfigThing *tls.Config // Has unexported fields and no env struct tag
}

type SmallTestStruct struct {
	HiddenThing           string     `env:"HIDDEN_THING"`
	StringThing           string     `env:"STRING_THING,report"`
	BoolThing             bool       `env:"BOOL_THING,report"`
	IntThing              int        `env:"INT_THING,report"`
	FloatThing            float64    `env:"FLOAT_THING,report"`
	ComplexThing          complex128 `env:"COMPLEX_THING,report"`
	URLThing              *url.URL   `env:"URL_THING,report"`
	StringSliceThing      []string   `env:"STRING_SLICE_THING,report"`
	CaseSensitiveThing    string     `env:"CaSe_SeNsItIvE_ThInG,report"`
	SmallTestSubStruct    SmallTestSubStruct
	PtrSmallTestSubStruct *SmallTestSubStruct
	NotReported           SmallTestStructWithNoEnv
}

type SmallTestSubStruct struct {
	SecretThing string `env:"SECRET_THING"`
}

type SmallTestStructWithNoEnv struct {
	FieldThing string
}

type SmallTestStructWithSubStructWithoutMarshaller struct {
	FieldThing noMarshaller `env:"FIELD_THING"`
}

type ToEnvTestStruct struct {
	HiddenThing        string   `env:"HIDDEN_THING,noreport"`
	StringThing        string   `env:"STRING_THING"`
	BoolThing          bool     `env:"BOOL_THING"`
	IntThing           int      `env:"INT_THING"`
	URLThing           *url.URL `env:"URL_THING"`
	StringSliceThing   []string `env:"STRING_SLICE_THING"`
	CaseSensitiveThing string   `env:"CaSe_SeNsItIvE_ThInG"`
	SubStruct          SubTestStruct
	SubPointerStruct   *SubTestStruct
}

type ToEnvMapTestStruct struct {
	MapStringStringThing map[string]string `env:"MAP_STRING_STRING_THING"`
}

type SubTestStruct struct {
	SubThingA string `env:"SUB_THING_A"`
	SubThingB int    `env:"SUB_THING_B,required"`
}

type UnsupportedType uintptr

type WithUnsupportedTypeStruct struct {
	Unsupported UnsupportedType `env:"UNSUPPORTED"`
}

type spyUnmarshaller struct {
	UnmarshalEnvInput  string
	UnmarshalEnvOutput error
}

func (s *spyUnmarshaller) UnmarshalEnv(v string) error {
	s.UnmarshalEnvInput = v
	return s.UnmarshalEnvOutput
}

type noMarshaller struct {
}

func TestEnvstruct(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Envstruct Suite")
}
