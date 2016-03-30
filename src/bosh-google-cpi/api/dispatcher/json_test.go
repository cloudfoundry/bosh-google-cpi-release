package dispatcher_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "bosh-google-cpi/api/dispatcher"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	fakeaction "bosh-google-cpi/action/fakes"
	fakedisp "bosh-google-cpi/api/dispatcher/fakes"
	fakeapi "bosh-google-cpi/api/fakes"
)

var _ = Describe("JSON", func() {
	var (
		actionFactory *fakeaction.FakeFactory
		caller        *fakedisp.FakeCaller
		logger        boshlog.Logger
		dispatcher    JSON
	)

	BeforeEach(func() {
		actionFactory = fakeaction.NewFakeFactory()
		caller = &fakedisp.FakeCaller{}
		logger = boshlog.NewLogger(boshlog.LevelNone)
		dispatcher = NewJSON(actionFactory, caller, logger)
	})

	Describe("Dispatch", func() {
		Context("when method is known", func() {
			var (
				action *fakeaction.FakeAction
			)

			BeforeEach(func() {
				action = &fakeaction.FakeAction{}
				actionFactory.RegisterAction("fake-action", action)
			})

			It("runs action with provided arguments", func() {
				dispatcher.Dispatch([]byte(`{"method":"fake-action","arguments":["fake-arg"]}`))
				Expect(caller.CallAction).To(Equal(action))
				Expect(caller.CallArgs).To(Equal([]interface{}{"fake-arg"}))

				dispatcher.Dispatch([]byte(`{
          "method":"fake-action",
          "arguments":[
            123,
            "fake-arg",
            [123, "fake-arg"],
            {"fake-arg2-key":"fake-arg2-value"}
          ]
        }`))
				Expect(caller.CallAction).To(Equal(action))
				Expect(caller.CallArgs).To(Equal([]interface{}{
					float64(123),
					"fake-arg",
					[]interface{}{float64(123), "fake-arg"},
					map[string]interface{}{"fake-arg2-key": "fake-arg2-value"},
				}))
			})

			Context("when running action succeeds", func() {
				Context("when result can be serialized", func() {
					BeforeEach(func() {
						caller.CallResult = "fake-result"
					})

					It("returns serialized result without including error", func() {
						response := dispatcher.Dispatch([]byte(`{"method":"fake-action","arguments":["fake-arg"]}`))
						Expect(response).To(MatchJSON(`{
							"result": "fake-result",
							"error": null,
							"log": ""
						}`))
					})
				})

				Context("when result cannot be serialized", func() {
					BeforeEach(func() {
						caller.CallResult = func() {} // funcs do not serialize
					})

					It("returns Bosh::Clouds::CpiError", func() {
						response := dispatcher.Dispatch([]byte(`{"method":"fake-action","arguments":["fake-arg"]}`))
						Expect(response).To(MatchJSON(`{
							"result": null,
              "error": {
                "type":"Bosh::Clouds::CpiError",
                "message":"Failed to serialize result",
                "ok_to_retry": false
              },
              "log": ""
            }`))
					})
				})
			})

			Context("when running action fails", func() {
				Context("when action error is a CloudError", func() {
					BeforeEach(func() {
						caller.CallErr = fakeapi.NewFakeCloudError("fake-type", "fake-message")
					})

					It("returns error without result", func() {
						response := dispatcher.Dispatch([]byte(`{"method":"fake-action","arguments":["fake-arg"]}`))
						Expect(response).To(MatchJSON(`{
							"result": null,
              "error": {
                "type":"fake-type",
                "message":"fake-message",
                "ok_to_retry": false
              },
              "log": ""
            }`))
					})
				})

				Context("when action error is a RetryableError and it can be retried", func() {
					BeforeEach(func() {
						caller.CallErr = fakeapi.NewFakeRetryableError("fake-error", true)
					})

					It("returns error with ok_to_retry set to true", func() {
						response := dispatcher.Dispatch([]byte(`{"method":"fake-action","arguments":["fake-arg"]}`))
						Expect(response).To(MatchJSON(`{
							"result": null,
              "error": {
                "type":"Bosh::Clouds::CloudError",
                "message":"fake-error",
                "ok_to_retry": true
              },
              "log": ""
            }`))
					})
				})

				Context("when action error is a RetryableError and it cannot be retried", func() {
					BeforeEach(func() {
						caller.CallErr = fakeapi.NewFakeRetryableError("fake-error", false)
					})

					It("returns error with ok_to_retry set to false", func() {
						response := dispatcher.Dispatch([]byte(`{"method":"fake-action","arguments":["fake-arg"]}`))
						Expect(response).To(MatchJSON(`{
							"result": null,
              "error": {
                "type":"Bosh::Clouds::CloudError",
                "message":"fake-error",
                "ok_to_retry": false
              },
              "log": ""
            }`))
					})
				})

				Context("when action error is neither CloudError or RetryableError", func() {
					BeforeEach(func() {
						caller.CallErr = errors.New("fake-run-err")
					})

					It("returns error without result", func() {
						response := dispatcher.Dispatch([]byte(`{"method":"fake-action","arguments":["fake-arg"]}`))
						Expect(response).To(MatchJSON(`{
							"result": null,
              "error": {
                "type":"Bosh::Clouds::CloudError",
                "message":"fake-run-err",
                "ok_to_retry": false
              },
              "log": ""
            }`))
					})
				})
			})
		})

		Context("when method is unknown", func() {
			It("responds with Bosh::Clouds::NotImplemented error", func() {
				response := dispatcher.Dispatch([]byte(`{"method":"fake-action","arguments":[]}`))
				Expect(response).To(MatchJSON(`{
					"result": null,
          "error": {
            "type":"Bosh::Clouds::NotImplemented",
            "message":"Must call implemented method",
            "ok_to_retry": false
          },
          "log": ""
        }`))
			})
		})

		Context("when method key is missing", func() {
			It("responds with Bosh::Clouds::CpiError error", func() {
				response := dispatcher.Dispatch([]byte(`{}`))
				Expect(response).To(MatchJSON(`{
					"result": null,
          "error": {
            "type":"Bosh::Clouds::CpiError",
            "message":"Must provide method key",
            "ok_to_retry": false
          },
          "log": ""
        }`))
			})
		})

		Context("when arguments key is missing", func() {
			It("responds with Bosh::Clouds::CpiError error", func() {
				response := dispatcher.Dispatch([]byte(`{"method":"fake-action"}`))
				Expect(response).To(MatchJSON(`{
					"result": null,
          "error": {
            "type":"Bosh::Clouds::CpiError",
            "message":"Must provide arguments key",
            "ok_to_retry": false
          },
          "log": ""
        }`))
			})
		})

		Context("when payload cannot be deserialized", func() {
			It("responds with Bosh::Clouds::CpiError error", func() {
				response := dispatcher.Dispatch([]byte(`{-}`))
				Expect(response).To(MatchJSON(`{
					"result": null,
          "error": {
            "type":"Bosh::Clouds::CpiError",
            "message":"Must provide valid JSON payload",
            "ok_to_retry": false
          },
          "log": ""
        }`))
			})
		})
	})
})
