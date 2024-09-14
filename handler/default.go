package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type LambdaHandler struct {
	baseEvent      events.LambdaFunctionURLRequest
	lambdaEndpoint string
}

func NewLambdaHandler(eventJson events.LambdaFunctionURLRequest, lambdaEndpoint string) *LambdaHandler {
	return &LambdaHandler{
		baseEvent:      eventJson,
		lambdaEndpoint: lambdaEndpoint,
	}
}

func (h *LambdaHandler) RequestToEvent(r *http.Request) (*events.LambdaFunctionURLRequest, error) {
	// will copy the base event
	out := h.baseEvent

	if r.URL.RawPath != "" {
		out.RawPath = r.URL.RawPath
	}

	if r.URL.RawQuery != "" {
		out.RawQueryString = r.URL.RawQuery
	}

	cookies := r.Cookies()
	for _, cookie := range cookies {
		out.Cookies = append(out.Cookies, cookie.String())
	}

	for k, v := range r.Header {
		out.Headers[strings.ToLower(k)] = v[0]
	}

	query := r.URL.Query()
	for k, v := range query {
		out.QueryStringParameters[k] = v[0]
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if len(body) > 0 {
		out.Body = string(body)
		out.IsBase64Encoded = false
	}

	return &out, nil
}

func (h *LambdaHandler) HandlerFunc(w http.ResponseWriter, r *http.Request) {
	event, err := h.RequestToEvent(r)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	proxyRequest, err := json.Marshal(event)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	resp, err := http.Post(h.lambdaEndpoint, "application/json", bytes.NewBuffer(proxyRequest))
	if err != nil {
		w.WriteHeader(resp.StatusCode)
		w.Write([]byte(err.Error()))
		return
	}

	respBody, err := io.ReadAll(resp.Body)
	var proxyResponse events.LambdaFunctionURLResponse
	err = json.Unmarshal(respBody, &proxyResponse)
	if err != nil {
		w.WriteHeader(resp.StatusCode)
		w.Write([]byte(err.Error()))
		return
	}
	for _, v := range proxyResponse.Cookies {
		w.Header().Add("Set-Cookie", v)
	}
	for k, v := range proxyResponse.Headers {
		w.Header().Add(k, v)
	}
	w.WriteHeader(proxyResponse.StatusCode)
	w.Write([]byte(proxyResponse.Body))
}
