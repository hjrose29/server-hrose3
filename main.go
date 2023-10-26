package main


import (
    "fmt"
    "net/http"
	"github.com/gorilla/mux"
	"time"
    "log"
    "github.com/jamespearly/loggly"
)

func main() {
    

	r := mux.NewRouter()
	r.HandleFunc("/hrose3/status", statusHandler).Methods(http.MethodGet)
	r.PathPrefix("/").HandlerFunc(catchAllHandler).Methods(http.MethodGet)
    r.Use(RequestLoggerMiddleware(r))
    http.ListenAndServe(":8080", r)

}



func statusHandler(w http.ResponseWriter, r *http.Request) {

    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, time.Now().String())
}

func catchAllHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "404 - Page not found!")
}

func RequestLoggerMiddleware(r *mux.Router) mux.MiddlewareFunc {
    client := loggly.New("Web Scraper")

    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
            defer func() {
                log.Printf(
                    "[%s] %s %s %s",
                    req.Method,
                    req.Host,
                    req.URL.Path,
                    req.URL.RawQuery,
                )
                    err := client.EchoSend("info", "Request made: [" + string(req.Method) + "]\t" + req.Host + "\t" + req.URL.Path + "\t" + req.URL.RawQuery)
                    if(err != nil){
                        fmt.Println(err)
                    }
                
            }()

            next.ServeHTTP(w, req)
        })
    }
}
