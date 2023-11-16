package dispatcher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"

	bgcaction "bosh-google-cpi/action"
	bgcapi "bosh-google-cpi/api"
	"bosh-google-cpi/constant"
	"bosh-google-cpi/redactor"
)

const (
	jsonLogTag = "json"

	jsonCloudErrorType          = "Bosh::Clouds::CloudError"
	jsonCpiErrorType            = "Bosh::Clouds::CpiError"
	jsonNotImplementedErrorType = "Bosh::Clouds::NotImplemented"
)

type Request struct {
	Method    string        `json:"method"`
	Arguments []interface{} `json:"arguments"`

	Context    map[string]interface{} `json:"context"`
	ApiVersion int                    `json:"api_version,omitempty"`
}

type Response struct {
	Result interface{}    `json:"result"`
	Error  *ResponseError `json:"error"`

	Log string `json:"log"`
}

type ResponseError struct {
	Type    string `json:"type"`
	Message string `json:"message"`

	CanRetry bool `json:"ok_to_retry"`
}

func (r ResponseError) Error() string {
	return r.Message
}

type JSON struct {
	actionFactory bgcaction.Factory
	caller        Caller
	logger        bgcapi.MultiLogger
}

func NewJSON(
	actionFactory bgcaction.Factory,
	caller Caller,
	logger bgcapi.MultiLogger,
) JSON {
	return JSON{
		actionFactory: actionFactory,
		caller:        caller,
		logger:        logger,
	}
}

func (c JSON) Dispatch(reqBytes []byte) []byte {
	var req Request

	c.logger.DebugWithDetails(jsonLogTag, "Request bytes", redactor.RedactSecrets(string(reqBytes)))

	decoder := json.NewDecoder(bytes.NewReader(reqBytes))
	decoder.UseNumber()
	if err := decoder.Decode(&req); err != nil {
		return c.buildCpiError("Must provide valid JSON payload")
	}

	if req.Method == "" {
		return c.buildCpiError("Must provide method key")
	}

	if req.Arguments == nil {
		return c.buildCpiError("Must provide arguments key")
	}

	apiVersion := int(math.Max(1, float64(req.ApiVersion)))

	if apiVersion > constant.MaxSupportedAPIVersion {
		return c.buildCpiError(fmt.Sprintf("API version %d not supported", apiVersion))
	}

	action, err := c.actionFactory.Create(req.Method, req.Context, apiVersion)
	if err != nil {
		return c.buildNotImplementedError()
	}

	result, err := c.caller.Call(action, req.Arguments)
	if err != nil {
		return c.buildCloudError(err)
	}

	resp := Response{
		Result: result,
		Log:    c.logger.LogBuff.String(),
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		return c.buildCpiError("Failed to serialize result")
	}

	c.logger.DebugWithDetails(jsonLogTag, "Response bytes", redactor.RedactSecrets(string(respBytes)))

	return respBytes
}

func (c JSON) buildCloudError(err error) []byte {
	respErr := Response{
		Log:   c.logger.LogBuff.String(),
		Error: &ResponseError{},
	}

	respErr.Log = c.logger.LogBuff.String()

	if typedErr, ok := err.(bgcapi.CloudError); ok {
		respErr.Error.Type = typedErr.Type()
	} else {
		respErr.Error.Type = jsonCloudErrorType
	}

	respErr.Error.Message = err.Error()

	if typedErr, ok := err.(bgcapi.RetryableError); ok {
		respErr.Error.CanRetry = typedErr.CanRetry()
	}

	respErrBytes, err := json.Marshal(respErr)
	if err != nil {
		panic(err)
	}

	c.logger.DebugWithDetails(jsonLogTag, "CloudError response bytes", redactor.RedactSecrets(string(respErrBytes)))

	return respErrBytes
}

func (c JSON) buildCpiError(message string) []byte {
	respErr := Response{
		Log: c.logger.LogBuff.String(),
		Error: &ResponseError{
			Type:    jsonCpiErrorType,
			Message: message,
		},
	}

	respErrBytes, err := json.Marshal(respErr)
	if err != nil {
		panic(err)
	}

	c.logger.DebugWithDetails(jsonLogTag, "CpiError response bytes", redactor.RedactSecrets(string(respErrBytes)))

	return respErrBytes
}

func (c JSON) buildNotImplementedError() []byte {
	respErr := Response{
		Log: c.logger.LogBuff.String(),
		Error: &ResponseError{
			Type:    jsonNotImplementedErrorType,
			Message: "Must call implemented method",
		},
	}

	respErrBytes, err := json.Marshal(respErr)
	if err != nil {
		panic(err)
	}

	c.logger.DebugWithDetails(jsonLogTag, "NotImplementedError response bytes", redactor.RedactSecrets(string(respErrBytes)))

	return respErrBytes
}
