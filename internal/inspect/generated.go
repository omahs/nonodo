// Package inspect provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.0.0 DO NOT EDIT.
package inspect

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
)

// Defines values for CompletionStatus.
const (
	Accepted           CompletionStatus = "Accepted"
	CycleLimitExceeded CompletionStatus = "CycleLimitExceeded"
	Exception          CompletionStatus = "Exception"
	MachineHalted      CompletionStatus = "MachineHalted"
	Rejected           CompletionStatus = "Rejected"
	TimeLimitExceeded  CompletionStatus = "TimeLimitExceeded"
)

// CompletionStatus Whether inspection completed or not (and why not)
type CompletionStatus string

// Error Detailed error message.
type Error = string

// InspectResult defines model for InspectResult.
type InspectResult struct {
	// ExceptionPayload Payload in the Ethereum hex binary format.
	// The first two characters are '0x' followed by pairs of hexadecimal numbers that correspond to one byte.
	// For instance, '0xdeadbeef' corresponds to a payload with length 4 and bytes 222, 173, 190, 175.
	// An empty payload is represented by the string '0x'.
	ExceptionPayload Payload `json:"exception_payload"`

	// ProcessedInputCount Number of processed inputs since genesis
	ProcessedInputCount int      `json:"processed_input_count"`
	Reports             []Report `json:"reports"`

	// Status Whether inspection completed or not (and why not)
	Status CompletionStatus `json:"status"`
}

// Payload Payload in the Ethereum hex binary format.
// The first two characters are '0x' followed by pairs of hexadecimal numbers that correspond to one byte.
// For instance, '0xdeadbeef' corresponds to a payload with length 4 and bytes 222, 173, 190, 175.
// An empty payload is represented by the string '0x'.
type Payload = string

// Report defines model for Report.
type Report struct {
	// Payload Payload in the Ethereum hex binary format.
	// The first two characters are '0x' followed by pairs of hexadecimal numbers that correspond to one byte.
	// For instance, '0xdeadbeef' corresponds to a payload with length 4 and bytes 222, 173, 190, 175.
	// An empty payload is represented by the string '0x'.
	Payload Payload `json:"payload"`
}

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// InspectPostWithBody request with any body
	InspectPostWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	// Inspect request
	Inspect(ctx context.Context, payload string, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) InspectPostWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewInspectPostRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) Inspect(ctx context.Context, payload string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewInspectRequest(c.Server, payload)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewInspectPostRequestWithBody generates requests for InspectPost with any type of body
func NewInspectPostRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewInspectRequest generates requests for Inspect
func NewInspectRequest(server string, payload string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "payload", runtime.ParamLocationPath, payload)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// InspectPostWithBodyWithResponse request with any body
	InspectPostWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*InspectPostResponse, error)

	// InspectWithResponse request
	InspectWithResponse(ctx context.Context, payload string, reqEditors ...RequestEditorFn) (*InspectResponse, error)
}

type InspectPostResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *InspectResult
}

// Status returns HTTPResponse.Status
func (r InspectPostResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r InspectPostResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type InspectResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *InspectResult
}

// Status returns HTTPResponse.Status
func (r InspectResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r InspectResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// InspectPostWithBodyWithResponse request with arbitrary body returning *InspectPostResponse
func (c *ClientWithResponses) InspectPostWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*InspectPostResponse, error) {
	rsp, err := c.InspectPostWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseInspectPostResponse(rsp)
}

// InspectWithResponse request returning *InspectResponse
func (c *ClientWithResponses) InspectWithResponse(ctx context.Context, payload string, reqEditors ...RequestEditorFn) (*InspectResponse, error) {
	rsp, err := c.Inspect(ctx, payload, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseInspectResponse(rsp)
}

// ParseInspectPostResponse parses an HTTP response from a InspectPostWithResponse call
func ParseInspectPostResponse(rsp *http.Response) (*InspectPostResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &InspectPostResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest InspectResult
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseInspectResponse parses an HTTP response from a InspectWithResponse call
func ParseInspectResponse(rsp *http.Response) (*InspectResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &InspectResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest InspectResult
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Inspect DApp state via POST
	// (POST /)
	InspectPost(ctx echo.Context) error
	// Inspect DApp state via GET
	// (GET /{payload})
	Inspect(ctx echo.Context, payload string) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// InspectPost converts echo context to params.
func (w *ServerInterfaceWrapper) InspectPost(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.InspectPost(ctx)
	return err
}

// Inspect converts echo context to params.
func (w *ServerInterfaceWrapper) Inspect(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "payload" -------------
	var payload string

	err = runtime.BindStyledParameterWithLocation("simple", false, "payload", runtime.ParamLocationPath, ctx.Param("payload"), &payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter payload: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.Inspect(ctx, payload)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST(baseURL+"/", wrapper.InspectPost)
	router.GET(baseURL+"/:payload", wrapper.Inspect)

}
