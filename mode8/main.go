package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	http.HandleFunc("/", defaultHandle)
	http.HandleFunc("/healthz", healthz)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func defaultHandle(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello simon !")
}

func healthz(w http.ResponseWriter, r *http.Request) {
	requestHeader := r.Header
	for key, value := range requestHeader {
		w.Header().Set(key, value[0])
	}
	w.Header().Set("VERSION", os.Getenv("VERSION"))
	var params []string
	if len(strings.Split(r.URL.RequestURI(), "?")) > 1 {
		params = strings.Split(strings.Split(r.URL.RequestURI(), "?")[1], "&")
	}
	responseObject := &ResponseObject{
		Host:    r.Host,
		Uri:     r.URL.RequestURI(),
		Param:   params,
		Version: os.Getenv("VERSION"),
		Env:     os.Getenv("ENV"),
	}
	response_str := "ok\n"
	response_str = response_str + "version:" + responseObject.Version + "\n"
	response_str = response_str + "env:" + responseObject.Env + "\n"
	response_str = response_str + "ip:" + responseObject.Host + "\n"
	response_str = response_str + "uri:" + responseObject.Uri + "\n"
	response_str = response_str + "params:" + strings.Join(responseObject.Param, ",") + "\n"
	io.WriteString(w, response_str)
	logStr := "host=" + r.Host + "|response status=" + strconv.Itoa(http.StatusOK)
	log.Println(logStr)
}

type ResponseObject struct {
	Host    string
	Uri     string
	Param   []string
	Version string
	Env     string
}
