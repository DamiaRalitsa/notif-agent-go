package helpers

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"

	"gitlab.playcourt.id/notif-agent-go/log"
	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
)

var contextHttpClient = "HTTPClient"

type FormData struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ToolsAPI struct{}

func CreateToolsAPI() APIInterface {
	return &ToolsAPI{}
}

type APIInterface interface {
	CallAPI(ctx context.Context, url, method string, payload interface{}, header []Header) (body []byte, err error)
	CallAPIWithBufferString(ctx context.Context, url, method string, payload string, header []Header) (body []byte, err error)
	CallAPIFormData(ctx context.Context, url, method string, formData []FormData, headers []Header) (body []byte, err error)
	CallAPIWithoutContext(url, method string, payload interface{}, header []Header) (body []byte, err error)
	CallAPIFormDataWithoutContext(url, method string, formData []FormData, headers []Header) (body []byte, err error)
	SendToTelegram(url, method, tokenBOT, chatID, text string, IsContainArstik bool)
}

// CallAPI is
func (t *ToolsAPI) CallAPI(ctx context.Context, url, method string, payload interface{}, header []Header) (body []byte, err error) {
	// var res *http.Response
	body, err = json.Marshal(payload)
	if err != nil {
		return
	}

	var req *http.Request
	// var w http.ResponseWriter
	// client := &http.Client{}

	req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return
	}

	req.Header.Add("content-type", "application/json")
	for _, e := range header {
		req.Header.Add(e.Key, e.Value)
	}

	client := httptrace.WrapClient(http.DefaultClient)
	res, err := client.Do(req.WithContext(ctx))
	if err != nil {
		// apm.CaptureError(ctx, err).Send()
		// http.Error(w, "failed to query backend", 500)
		return
	}
	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.GetLogger().Error(contextHttpClient, fmt.Sprintf("%v", err.Error()), "CallAPI", "")
		return
	}

	return
}

func (t *ToolsAPI) CallAPIWithBufferString(ctx context.Context, url, method string, payload string, header []Header) (body []byte, err error) {
	// var res *http.Response

	var req *http.Request
	// var w http.ResponseWriter
	// client := &http.Client{}

	req, err = http.NewRequest(method, url, bytes.NewBufferString(payload))
	if err != nil {
		return
	}

	req.Header.Add("content-type", "application/json")
	for _, e := range header {
		req.Header.Add(e.Key, e.Value)
	}

	client := httptrace.WrapClient(http.DefaultClient)
	res, err := client.Do(req.WithContext(ctx))
	if err != nil {
		// apm.CaptureError(ctx, err).Send()
		// http.Error(w, "failed to query backend", 500)
		return
	}
	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.GetLogger().Error(contextHttpClient, fmt.Sprintf("%v", err.Error()), "CallAPIWithBufferString", "")
		return
	}

	return
}

func (t *ToolsAPI) CallAPIFormData(ctx context.Context, url, method string, formData []FormData, headers []Header) (body []byte, err error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	for _, element := range formData {
		_ = writer.WriteField(element.Key, element.Value)
	}
	err = writer.Close()
	if err != nil {
		log.GetLogger().Error(contextHttpClient, fmt.Sprintf("%v", err.Error()), "CallAPIFormData", "writer.Close")
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		log.GetLogger().Error(contextHttpClient, fmt.Sprintf("%v", err.Error()), "CallAPIFormData", "http.NewRequest")
		return
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())
	for _, e := range headers {
		req.Header.Add(e.Key, e.Value)
	}
	res, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.GetLogger().Error(contextHttpClient, fmt.Sprintf("%v", err.Error()), "CallAPIFormData", "ioutil.ReadAll")
		return
	}

	if res.StatusCode > 399 {
		err = fmt.Errorf(string(body))
		log.GetLogger().Error(contextHttpClient, fmt.Sprintf("%v", err.Error()), "CallAPIFormData", "res.StatusCode > 399")
		return
	}

	return
}

func (t *ToolsAPI) SendToTelegram(url, method, tokenBOT, chatID, text string, IsContainArstik bool) {
	client := &http.Client{}
	payload := strings.NewReader("chat_id=" + chatID + "&text=" + text + "&parse_mode=Markdown")
	if IsContainArstik {
		payload = strings.NewReader("chat_id=" + chatID + "&text=" + text)
	}
	req, err := http.NewRequest(method, url+tokenBOT+"/sendMessage", payload)
	if err != nil {
		log.GetLogger().Error(contextHttpClient, fmt.Sprintf("%v", err.Error()), "SendToTelegram", "http.NewRequest")
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		log.GetLogger().Error(contextHttpClient, fmt.Sprintf("%v", err.Error()), "SendToTelegram", "client.Do")
		return
	}
	defer res.Body.Close()
}

func (t *ToolsAPI) CallAPIWithoutContext(url, method string, payload interface{}, header []Header) (body []byte, err error) {
	// var res *http.Response
	body, err = json.Marshal(payload)
	if err != nil {
		return
	}

	var req *http.Request
	// var w http.ResponseWriter
	// client := &http.Client{}

	req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return
	}

	req.Header.Add("content-type", "application/json")
	for _, e := range header {
		req.Header.Add(e.Key, e.Value)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	res, err := client.Do(req)
	if err != nil {
		// apm.CaptureError(ctx, err).Send()
		// http.Error(w, "failed to query backend", 500)
		return
	}
	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.GetLogger().Error(contextHttpClient, fmt.Sprintf("%v", err.Error()), "CallAPIWithoutContext", "ioutil.ReadAll")
		return
	}

	return
}

func (t *ToolsAPI) CallAPIFormDataWithoutContext(url, method string, formData []FormData, headers []Header) (body []byte, err error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	for _, element := range formData {
		_ = writer.WriteField(element.Key, element.Value)
	}
	err = writer.Close()
	if err != nil {
		log.GetLogger().Error(contextHttpClient, fmt.Sprintf("%v", err.Error()), "CallAPIFormDataWithoutContext", "writer.Close")
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		log.GetLogger().Error(contextHttpClient, fmt.Sprintf("%v", err.Error()), "CallAPIFormDataWithoutContext", "http.NewRequest")
		return
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())
	for _, e := range headers {
		req.Header.Add(e.Key, e.Value)
	}
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.GetLogger().Error(contextHttpClient, fmt.Sprintf("%v", err.Error()), "CallAPIFormDataWithoutContext", "ioutil.ReadAll")
		return
	}

	return
}
