package main


import (
    "fmt"
    "net/http"
	"github.com/gorilla/mux"
	"time"
)
func main() {

	r := mux.NewRouter()
	r.HandleFunc("/hrose3/status", statusHandler).Methods(http.MethodGet)
	r.PathPrefix("/").HandlerFunc(catchAllHandler)
    
    http.ListenAndServe(":8080", r)

}


func statusHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, time.Now().String())
}

func catchAllHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "404 - Page not found!")
}