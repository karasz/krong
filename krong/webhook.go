package krong

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Header struct {
	Name  string
	Value string
}
type Headers []Header

type WebHook struct {
	URL     string
	Method  string
	Headers Headers
	Payload []byte
	Timeout time.Duration
}

func (h *Headers) String() string {
	if len(*h) == 0 {
		return ""
	}

	result := make([]string, len(*h))

	for i, header := range *h {
		result[i] = fmt.Sprintf("%s=%s", header.Name, header.Value)
	}

	return strings.Join(result, ", ")
}

func (h *Headers) Set(value string) error {
	splitResult := strings.SplitN(value, "=", 2)

	if len(splitResult) != 2 {
		return fmt.Errorf("header flag must be in name=value format")
	}

	*h = append(*h, Header{Name: splitResult[0], Value: splitResult[1]})
	return nil
}

func (w *WebHook) Run() ([]byte, string, error) {
	c := http.Client{Timeout: w.Timeout}
	req, err := http.NewRequest(w.Method, w.URL, bytes.NewBuffer(w.Payload))
	if err != nil {
		return []byte(""), "", err
	}

	for _, h := range w.Headers {
		req.Header.Add(h.Name, h.Value)
	}

	resp, err := c.Do(req)
	if err != nil {
		return []byte(""), "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), "", err
	}
	return body, resp.Status, nil
}
func (w *WebHook) String() string {
	return fmt.Sprintf("URL: %s, Method: %s, Payload: %s, Timeout: %s", w.URL, w.Method, w.Payload, w.Timeout)
}
