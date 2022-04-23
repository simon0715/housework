package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	http.HandleFunc("/healthz", healthz)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func healthz(w http.ResponseWriter, r *http.Request) {
	requestHeader := r.Header
	for key, value := range requestHeader {
		w.Header().Set(key, value[0])
	}
	w.Header().Set("VERSION", os.Getenv("VERSION"))
	logStr := "host=" + r.Host + "|response status=" + strconv.Itoa(http.StatusOK)
	log.Println(logStr)
	io.WriteString(w, "ok")
}
