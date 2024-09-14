package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/k0kubun/pp/v3"
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
		out.Headers[k] = v[0]
	}

	query := r.URL.Query()
	for k, v := range query {
		out.QueryStringParameters[k] = v[0]
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	out.Body = string(body)

	pp.Println(out)

	return &out, nil
}

func (h *LambdaHandler) HandlerFunc(w http.ResponseWriter, r *http.Request) {
	event, err := h.RequestToEvent(r)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	json, err := json.Marshal(event)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	resp, err := http.Post(h.lambdaEndpoint, "application/json", bytes.NewBuffer(json))
	if err != nil {
		w.WriteHeader(resp.StatusCode)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(resp.StatusCode)
	respBody, err := io.ReadAll(resp.Body)

	w.Write(respBody)
}
