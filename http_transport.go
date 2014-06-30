package wray

import (
  "net/url"
  "net/http"
  "encoding/json"
  "bytes"
  "io/ioutil"
  "errors"
)

type HttpTransport struct {
  url string
}

func(self HttpTransport) isUsable(clientUrl string) bool {
  parsedUrl, err := url.Parse(clientUrl)
  if err != nil {
    return false
  }
  if parsedUrl.Scheme == "http" || parsedUrl.Scheme == "https" {
    return true
  }
  return false
}

func(self HttpTransport) connectionType() string {
  return "long-polling"
}

func(self HttpTransport) send(data map[string]interface{}, callback func(Response)) {
  dataBytes, _ := json.Marshal(data)
  buffer := bytes.NewBuffer(dataBytes)
  responseData, err := http.Post(self.url, "application/json", buffer)
  if err != nil {
    callback(Response{successful: false, error: err})
    return
  }
  if responseData.StatusCode != 200 {
    callback(Response{successful: false, error: errors.New(responseData.Status)})
    return
  }
  readData, _ := ioutil.ReadAll(responseData.Body)
  responseData.Body.Close()
  var jsonData []interface{}
  json.Unmarshal(readData, &jsonData)
  response := newResponse(jsonData)
  callback(response)
}

func(self *HttpTransport) setUrl(url string) {
  self.url = url
}
