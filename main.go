package main

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"text/template"

	yaml "gopkg.in/yaml.v3"
)

func main() {
	prepareHTTP()
	log.Printf("Starting server on %s", "0.0.0.0:10099")
	log.Fatal(http.ListenAndServe(":10099", nil))
}

type responseData struct {
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func writeResponse(result *responseData, w http.ResponseWriter) error {
	b, err := json.Marshal(result)
	if err != nil || result.Code == 0 {
		w.Write([]byte("{\"code\": 500, \"message\": \"internal error\"}"))
		return err
	}
	w.Write(b)
	return nil
}

func prepareHTTP() {
	http.Handle("/static", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		tem, err := template.ParseFiles("views/index.html")
		if err != nil {
			log.Printf("parse html view/index.html error %v", err)
			return
		}
		tem.Execute(writer, nil)
	})
	http.HandleFunc("/yaml2json", func(writer http.ResponseWriter, request *http.Request) {
		var ret = &responseData{}
		defer writeResponse(ret, writer)
		if http.MethodPost != request.Method {
			ret.Code = http.StatusMethodNotAllowed
			ret.Message = "yaml2json not allowed method" + request.Method + ", must be " + http.MethodPost
			return
		}

		var (
			body []byte
			data interface{}
			out  []byte
			err  error
		)
		body, err = ioutil.ReadAll(request.Body)
		if err != nil {
			log.Printf("yaml to json read body err %v\n", err)
			ret.Code = http.StatusBadRequest
			ret.Message = "read body failed, " + err.Error()
			return
		}
		err = yaml.Unmarshal(body, &data)
		if err != nil {
			log.Printf("yaml body convert to json err %v", err)
			ret.Code = http.StatusBadRequest
			ret.Message = "convert body failed, " + err.Error()
			return
		}
		out, err = json.Marshal(data)
		if err != nil {
			log.Printf("yaml data convert to json err %v", err)
			ret.Code = http.StatusBadRequest
			ret.Message = "marshal body failed, " + err.Error()
			return
		}
		ret.Code = http.StatusOK
		ret.Message = "Success"
		ret.Data = string(out)
	})
	http.HandleFunc("/json2yaml", func(writer http.ResponseWriter, request *http.Request) {
		var ret = &responseData{}
		defer writeResponse(ret, writer)
		if http.MethodPost != request.Method {
			ret = &responseData{
				Code:    http.StatusMethodNotAllowed,
				Message: "json2yaml not allowed method" + request.Method + ", must be " + http.MethodPost,
			}
			return
		}

		var (
			body []byte
			data interface{}
			out  []byte
			err  error
		)
		body, err = ioutil.ReadAll(request.Body)
		if err != nil {
			log.Printf("json to yaml read body err %v\n", err)
			ret = &responseData{
				Code:    http.StatusBadRequest,
				Message: "read body failed, " + err.Error(),
			}
			return
		}
		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Printf("json body unmarshal err %v", err)
			ret = &responseData{
				Code:    http.StatusBadRequest,
				Message: "convert body failed, " + err.Error(),
			}
			return
		}
		out, err = yaml.Marshal(data)
		if err != nil {
			log.Printf("json data convert to yaml err %v", err)
			ret = &responseData{
				Code:    http.StatusBadRequest,
				Message: "convert body failed, " + err.Error(),
			}
			return
		}
		ret.Code = http.StatusOK
		ret.Message = "Success"
		ret.Data = string(out)
	})
	http.HandleFunc("/json2string", func(writer http.ResponseWriter, request *http.Request) {
		var ret = &responseData{}
		defer writeResponse(ret, writer)
		if http.MethodPost != request.Method {
			ret.Code = http.StatusMethodNotAllowed
			ret.Message = "json2string not allowed method" + request.Method + ", must be " + http.MethodPost
			return
		}

		var (
			body     []byte
			jsonData interface{}
			out      []byte
			dst      = new(bytes.Buffer)
			err      error
		)
		body, err = ioutil.ReadAll(request.Body)
		if err != nil {
			log.Printf("json to string read body err %v\n", err)
			ret.Code = http.StatusBadRequest
			ret.Message = "read body failed, " + err.Error()
			return
		}
		err = json.Unmarshal(body, &jsonData)
		if err != nil {
			log.Printf("json body unmarshal err %v", err)
			ret.Code = http.StatusBadRequest
			ret.Message = "convert body failed, " + err.Error()
			return
		}
		out, err = json.Marshal(jsonData)
		if err != nil {
			log.Printf("json body marshal err %v", err)
			ret.Code = http.StatusBadRequest
			ret.Message = "convert body failed, " + err.Error()
			return
		}
		err = json.Compact(dst, out)
		if err != nil {
			log.Printf("json body compact err %v", err)
			ret.Code = http.StatusBadRequest
			ret.Message = "convert body failed, " + err.Error()
			return
		}
		//data = strconv.Quote(string(out))
		//data = data[1 : len(data)-1]
		ret.Code = http.StatusOK
		ret.Message = "Success"
		ret.Data = dst.String()
	})
	http.HandleFunc("/string2json", func(writer http.ResponseWriter, request *http.Request) {
		var ret = &responseData{}
		defer writeResponse(ret, writer)
		if http.MethodPost != request.Method {
			ret.Code = http.StatusMethodNotAllowed
			ret.Message = "string2json not allowed method" + request.Method + ", must be " + http.MethodPost
			return
		}

		var (
			body []byte
			data interface{}
			out  []byte
			err  error
		)
		body, err = ioutil.ReadAll(request.Body)
		if err != nil {
			log.Printf("string to json read body err %v\n", err)
			ret.Code = http.StatusBadRequest
			ret.Message = "read body failed, " + err.Error()
			return
		}
		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Printf("string body unmarshal err %v", err)
			ret.Code = http.StatusBadRequest
			ret.Message = "convert body failed, " + err.Error()
			return
		}
		out, err = json.MarshalIndent(data, "", "    ")
		if err != nil {
			log.Printf("json marshal indent err %v", err)
			ret.Code = http.StatusBadRequest
			ret.Message = "marshal indent failed, " + err.Error()
			return
		}
		ret.Code = http.StatusOK
		ret.Message = "Success"
		ret.Data = string(out)
	})
	http.HandleFunc("/base64encode", func(writer http.ResponseWriter, request *http.Request) {
		var ret = &responseData{}
		defer writeResponse(ret, writer)
		if http.MethodPost != request.Method {
			ret.Code = http.StatusMethodNotAllowed
			ret.Message = "base64 encode not allowed method" + request.Method + ", must be " + http.MethodPost
			return
		}

		var (
			body []byte
			data interface{}
			err  error
		)
		body, err = ioutil.ReadAll(request.Body)
		if err != nil {
			log.Printf("string to json read body err %v\n", err)
			ret.Code = http.StatusBadRequest
			ret.Message = "read body failed, " + err.Error()
			return
		}
		data = base64.StdEncoding.EncodeToString(body)
		ret.Code = http.StatusOK
		ret.Message = "Success"
		ret.Data = data
	})

	http.HandleFunc("/base64decode", func(writer http.ResponseWriter, request *http.Request) {
		var ret = &responseData{}
		defer writeResponse(ret, writer)
		if http.MethodPost != request.Method {
			ret.Code = http.StatusMethodNotAllowed
			ret.Message = "base64 decode not allowed method" + request.Method + ", must be " + http.MethodPost
			return
		}

		var (
			body []byte
			data []byte
			err  error
		)
		body, err = ioutil.ReadAll(request.Body)
		if err != nil {
			log.Printf("base64 decode read body err %v\n", err)
			ret.Code = http.StatusBadRequest
			ret.Message = "read body failed, " + err.Error()
			return
		}
		data, err = base64.StdEncoding.DecodeString(string(body))
		if err != nil {
			log.Printf("base64 decode err %v\n", err)
			ret.Code = http.StatusBadRequest
			ret.Message = "base64 decode failed, " + err.Error()
			return
		}
		ret.Code = http.StatusOK
		ret.Message = "Success"
		ret.Data = string(data)
	})
	http.HandleFunc("/urlencode", func(writer http.ResponseWriter, request *http.Request) {
		var ret = &responseData{}
		defer writeResponse(ret, writer)
		if http.MethodPost != request.Method {
			ret.Code = http.StatusMethodNotAllowed
			ret.Message = "url encode not allowed method" + request.Method + ", must be " + http.MethodPost
			return
		}

		var (
			body []byte
			data string
			err  error
		)
		body, err = ioutil.ReadAll(request.Body)
		if err != nil {
			log.Printf("url encode read body err %v\n", err)
			ret.Code = http.StatusBadRequest
			ret.Message = "read body failed, " + err.Error()
			return
		}
		data = url.QueryEscape(string(body))
		ret.Code = http.StatusOK
		ret.Message = "Success"
		ret.Data = data
	})

	http.HandleFunc("/urldecode", func(writer http.ResponseWriter, request *http.Request) {
		var ret = &responseData{}
		defer writeResponse(ret, writer)
		if http.MethodPost != request.Method {
			ret.Code = http.StatusMethodNotAllowed
			ret.Message = "url decode not allowed method" + request.Method + ", must be " + http.MethodPost
			return
		}

		var (
			body []byte
			data string
			err  error
		)
		body, err = ioutil.ReadAll(request.Body)
		if err != nil {
			log.Printf("url decode read body err %v\n", err)
			ret.Code = http.StatusBadRequest
			ret.Message = "read body failed, " + err.Error()
			return
		}
		data, err = url.QueryUnescape(string(body))
		if err != nil {
			log.Printf("url decode err %v\n", err)
			ret.Code = http.StatusBadRequest
			ret.Message = "url decode failed, " + err.Error()
			return
		}
		ret.Code = http.StatusOK
		ret.Message = "Success"
		ret.Data = data
	})
	http.HandleFunc("/md5", func(writer http.ResponseWriter, request *http.Request) {
		var ret = &responseData{}
		defer writeResponse(ret, writer)
		if http.MethodPost != request.Method {
			ret.Code = http.StatusBadRequest
			ret.Message = "md5 not allowed method" + request.Method + ", must be " + http.MethodPost
			return
		}

		var (
			body    []byte
			data    []byte
			md5Hash = md5.New()
			err     error
		)
		body, err = ioutil.ReadAll(request.Body)
		if err != nil {
			log.Printf("md5 read body err %v\n", err)
			ret.Code = http.StatusBadRequest
			ret.Message = "read body failed, " + err.Error()
			return
		}
		_, err = md5Hash.Write(body)
		if err != nil {
			log.Printf("md5 encode err %v\n", err)
			ret.Code = http.StatusBadRequest
			ret.Message = "md5 encode failed, " + err.Error()
			return
		}
		data = md5Hash.Sum(nil)
		ret.Code = http.StatusOK
		ret.Message = "Success"
		ret.Data = hex.EncodeToString(data)
	})
}
