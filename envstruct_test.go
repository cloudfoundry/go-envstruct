package envstruct_test

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"time"

	envstruct "code.cloudfoundry.org/go-envstruct"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("envstruct", func() {
	Describe("Load()", func() {
		var (
			ts        LargeTestStruct
			loadError error
			envVars   map[string]string
		)

		BeforeEach(func() {
			ts = LargeTestStruct{}
			ts.UnmarshallerPointer = &spyUnmarshaller{}
			ts.UnmarshallerValue = spyUnmarshaller{}

			envVars = make(map[string]string)
			for k, v := range baseEnvVars {
				envVars[k] = v
			}
		})

		JustBeforeEach(func() {
			for k, v := range envVars {
				os.Setenv(k, v)
			}
		})

		Context("when load is successful", func() {
			JustBeforeEach(func() {
				loadError = envstruct.Load(&ts)
			})

			AfterEach(func() {
				for k := range envVars {
					os.Setenv(k, "")
				}
			})

			It("does not return an error", func() {
				Expect(loadError).ToNot(HaveOccurred())
			})

			Context("with unmarshallers", func() {
				It("passes the value to the pointer field", func() {
					Expect(ts.UnmarshallerPointer.UnmarshalEnvInput).To(Equal("pointer"))
				})

				It("passes the value to the value field's address", func() {
					Expect(ts.UnmarshallerValue.UnmarshalEnvInput).To(Equal("value"))
				})
			})

			Context("with strings", func() {
				It("populates the string thing", func() {
					Expect(ts.StringThing).To(Equal("stringy thingy"))
				})
			})

			Describe("case sensitiveity", func() {
				It("populates the case sensitive thing", func() {
					Expect(ts.CaseSensitiveThing).To(Equal("case sensitive"))
				})
			})

			Context("with bools", func() {
				Context("with 'true'", func() {
					It("is true", func() {
						Expect(ts.BoolThing).To(BeTrue())
					})
				})

				Context("with 'false'", func() {
					BeforeEach(func() {
						envVars["BOOL_THING"] = "false"
					})

					It("is true", func() {
						Expect(ts.BoolThing).To(BeFalse())
					})
				})

				Context("with '1'", func() {
					BeforeEach(func() {
						envVars["BOOL_THING"] = "1"
					})

					It("is true", func() {
						Expect(ts.BoolThing).To(BeTrue())
					})
				})

				Context("with '0'", func() {
					BeforeEach(func() {
						envVars["BOOL_THING"] = "0"
					})

					It("is false", func() {
						Expect(ts.BoolThing).To(BeFalse())
					})
				})
			})

			Context("with ints", func() {
				It("populates the int thing", func() {
					Expect(ts.IntThing).To(Equal(100))
				})

				It("populates the int 8 thing", func() {
					Expect(ts.Int8Thing).To(Equal(int8(20)))
				})

				It("populates the int 16 thing", func() {
					Expect(ts.Int16Thing).To(Equal(int16(2000)))
				})

				It("populates the int 32 thing", func() {
					Expect(ts.Int32Thing).To(Equal(int32(200000)))
				})

				It("populates the int 64 thing", func() {
					Expect(ts.Int64Thing).To(Equal(int64(200000000)))
				})
			})

			Context("with uints", func() {
				It("populates the uint thing", func() {
					Expect(ts.UintThing).To(Equal(uint(100)))
				})

				It("populates the uint 8 thing", func() {
					Expect(ts.Uint8Thing).To(Equal(uint8(20)))
				})

				It("populates the uint 16 thing", func() {
					Expect(ts.Uint16Thing).To(Equal(uint16(2000)))
				})

				It("populates the uint 32 thing", func() {
					Expect(ts.Uint32Thing).To(Equal(uint32(200000)))
				})

				It("populates the uint 64 thing", func() {
					Expect(ts.Uint64Thing).To(Equal(uint64(200000000)))
				})
			})

			Context("with comma separated strings", func() {
				Context("slice of strings", func() {
					It("populates a slice of strings", func() {
						Expect(ts.StringSliceThing).To(Equal([]string{"one", "two", "three"}))
					})

					Context("with leading and trailing spaces", func() {
						BeforeEach(func() {
							envVars["STRING_SLICE_THING"] = "one , two , three"
						})

						It("populates a slice of strings", func() {
							Expect(ts.StringSliceThing).To(Equal([]string{"one", "two", "three"}))
						})
					})
				})

				Context("slice of ints", func() {
					It("populates a slice of ints", func() {
						Expect(ts.IntSliceThing).To(Equal([]int{1, 2, 3}))
					})
				})
			})

			Context("with map[string]string", func() {
				It("takes a comma separated list of key:value", func() {
					Expect(ts.MapStringStringThing).To(Equal(map[string]string{
						"key_one": "value_one",
						"key_two": "value_two:with_colon",
					}))
				})

				Context("when no value is given", func() {
					BeforeEach(func() {
						envVars["MAP_STRING_STRING_THING"] = "key"
					})

					It("returns an error if value is missing", func() {
						Expect(loadError).To(MatchError("map[string]string key 'key' is missing a value"))
					})
				})
			})

			Context("with a sub struct contains env tags", func() {
				It("populates the values of the substruct", func() {
					Expect(ts.SubStruct.SubThingA).To(Equal("sub-string-a"))
					Expect(ts.SubStruct.SubThingB).To(Equal(200))
				})

				It("populates the values of the pointer to substruct", func() {
					Expect(ts.SubPointerStruct).ToNot(BeNil())
					Expect(ts.SubPointerStruct.SubThingA).To(Equal("sub-string-a"))
					Expect(ts.SubPointerStruct.SubThingB).To(Equal(200))
				})

				Describe("with default values", func() {
					BeforeEach(func() {
						ts.SubStruct = SubTestStruct{
							SubThingA: "default-sub-a",
						}
						ts.SubPointerStruct = &SubTestStruct{
							SubThingA: "default-sub-pointer-a",
						}

						envVars["SUB_THING_A"] = ""
					})

					It("maintains the default values", func() {
						Expect(ts.SubStruct.SubThingA).To(Equal("default-sub-a"))
						Expect(ts.SubPointerStruct.SubThingA).To(Equal("default-sub-pointer-a"))
					})
				})
			})

			Context("with duration struct", func() {
				It("parses the duration string", func() {
					Expect(ts.DurationThing).To(Equal(2 * time.Second))
				})
			})

			Context("with url struct", func() {
				It("parses the url string", func() {
					Expect(ts.URLThing.Scheme).To(Equal("http"))
					Expect(ts.URLThing.Host).To(Equal("github.com"))
					Expect(ts.URLThing.Path).To(Equal("/some/path"))
				})
			})
		})

		Context("with defaults", func() {
			It("honors default values if env var is empty", func() {
				ts.DefaultThing = "Default Value"

				Expect(envstruct.Load(&ts)).To(Succeed())
				Expect(ts.DefaultThing).To(Equal("Default Value"))
			})
		})

		Context("when load is unsuccessful", func() {
			Context("when a required environment variable is not given", func() {
				BeforeEach(func() {
					envVars["REQUIRED_THING_A"] = ""
					envVars["REQUIRED_THING_B"] = ""
				})

				It("includes all required environment variables in the error", func() {
					loadError = envstruct.Load(&ts)

					Expect(loadError).To(MatchError(fmt.Errorf("missing required environment variables: REQUIRED_THING_A, REQUIRED_THING_B")))
				})
			})

			Context("when a required environment variable for substruct is not given", func() {
				BeforeEach(func() {
					envVars["SUB_THING_B"] = ""
				})

				It("returns a validation error", func() {
					loadError = envstruct.Load(&ts)

					Expect(loadError).To(MatchError(fmt.Errorf("missing required environment variables: SUB_THING_B")))
				})
			})

			Context("when top level and substruct are missing required arguments", func() {
				BeforeEach(func() {
					envVars["REQUIRED_THING_A"] = ""
					envVars["SUB_THING_B"] = ""
				})

				It("returns an error with both environment variables", func() {
					loadError = envstruct.Load(&ts)

					Expect(loadError).To(MatchError(fmt.Errorf("missing required environment variables: REQUIRED_THING_A, SUB_THING_B")))
				})
			})

			Context("with an invalid int", func() {
				BeforeEach(func() {
					envVars["INT_THING"] = "Hello!"
				})

				It("returns an error", func() {
					Expect(envstruct.Load(&ts)).ToNot(Succeed())
				})
			})

			Context("with an invalid uint", func() {
				BeforeEach(func() {
					envVars["UINT_THING"] = "Hello!"
				})

				It("returns an error", func() {
					Expect(envstruct.Load(&ts)).ToNot(Succeed())
				})
			})

			Context("with a failing unmarshaller pointer", func() {
				BeforeEach(func() {
					ts.UnmarshallerPointer.UnmarshalEnvOutput = errors.New("failed to unmarshal")
				})

				It("returns an error", func() {
					Expect(envstruct.Load(&ts)).ToNot(Succeed())
				})
			})

			Context("with a failing unmarshaller value", func() {
				BeforeEach(func() {
					ts.UnmarshallerValue.UnmarshalEnvOutput = errors.New("failed to unmarshal")
				})

				It("returns an error", func() {
					Expect(envstruct.Load(&ts)).ToNot(Succeed())
				})
			})

			Context("with a missing unmarshaller on struct with an env tag", func() {
				var withoutMarhsaller SmallTestStructWithSubStructWithoutMarshaller
				BeforeEach(func() {
					withoutMarhsaller = SmallTestStructWithSubStructWithoutMarshaller{}
				})

				It("returns an error", func() {
					Expect(envstruct.Load(&withoutMarhsaller)).ToNot(Succeed())
				})
			})
		})
	})

	Describe("ToEnv", func() {
		It("returns a slice of strings formatted as KEY=value", func() {
			url, err := url.Parse("https://example.com")
			Expect(err).ToNot(HaveOccurred())

			ts := ToEnvTestStruct{
				HiddenThing:        "hidden-thing",
				StringThing:        "string-thing",
				BoolThing:          true,
				IntThing:           200,
				URLThing:           url,
				StringSliceThing:   []string{"thing-1", "thing-2", "thing-3"},
				CaseSensitiveThing: "case-sensitive-thing",
				SubStruct: SubTestStruct{
					SubThingA: "sub-string-a",
					SubThingB: 300,
				},
				SubPointerStruct: &SubTestStruct{
					SubThingA: "sub-pointer-thing-a",
					SubThingB: 500,
				},
			}

			for k, v := range baseEnvVars {
				os.Setenv(k, v)
			}
			ret := envstruct.ToEnv(&ts)

			Expect(ret).To(ConsistOf(
				"HIDDEN_THING=hidden-thing",
				"STRING_THING=string-thing",
				"BOOL_THING=true",
				"INT_THING=200",
				"URL_THING=https://example.com",
				"STRING_SLICE_THING=thing-1,thing-2,thing-3",
				"CaSe_SeNsItIvE_ThInG=case-sensitive-thing",
				"SUB_THING_A=sub-string-a",
				"SUB_THING_B=300",
				"SUB_THING_A=sub-pointer-thing-a",
				"SUB_THING_B=500",
			))
		})

		Context("with a map", func() {
			It("returns a slice with a formatted map for environment variable", func() {
				ts := ToEnvMapTestStruct{
					MapStringStringThing: map[string]string{
						"key_one": "value_one",
						"key_two": "value_two",
					},
				}
				for k, v := range baseEnvVars {
					os.Setenv(k, v)
				}
				ret := envstruct.ToEnv(&ts)

				Expect(ret[0]).To(ContainSubstring("MAP_STRING_STRING_THING="))
				Expect(ret[0]).To(ContainSubstring("key_one:value_one"))
				Expect(ret[0]).To(ContainSubstring(","))
				Expect(ret[0]).To(ContainSubstring("key_two:value_two"))
			})
		})
	})
})
