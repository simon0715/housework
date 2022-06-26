package main

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	Register()
	mux := http.NewServeMux()
	mux.HandleFunc("/", defaultHandle)
	mux.HandleFunc("/healthz", healthz)
	mux.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func defaultHandle(w http.ResponseWriter, r *http.Request) {
	glog.V(4).Info("entering root handler")
	timer := NewTimer()
	defer timer.ObserveTotal()
	user := r.URL.Query().Get("user")
	delay := rand.Intn(2000)
	time.Sleep(time.Millisecond * time.Duration(delay))
	if user != "" {
		io.WriteString(w, fmt.Sprintf("hello [#{user}]\n"))
	}
	glog.V(4).Infof("Respond in %d ms", delay)
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
