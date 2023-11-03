package main


import (
    "fmt"
    "net/http"
	"github.com/gorilla/mux"
	"time"
    "log"
    "github.com/jamespearly/loggly"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
    "encoding/json"
)

type Datum struct {
	Ticker string
	DateTime int64
	Price float32
}

func main() {
    

	r := mux.NewRouter()
	r.HandleFunc("/hrose3/status", statusHandler).Methods(http.MethodGet)
    r.HandleFunc("/hrose3/all", allHandler).Methods(http.MethodGet)
    r.HandleFunc("/hrose3/search", searchHandler).Methods(http.MethodGet)
	r.PathPrefix("/").HandlerFunc(catchAllHandler).Methods(http.MethodGet)
    r.Use(RequestLoggerMiddleware(r))
    http.ListenAndServe(":8080", r)

}



func statusHandler(w http.ResponseWriter, r *http.Request) {

    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, time.Now().String())
}

func allHandler(w http.ResponseWriter, r *http.Request) {
    sess := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))

    // Create DynamoDB client
    svc := dynamodb.New(sess)

    input := &dynamodb.ScanInput{
        TableName: aws.String("hrose3"),
    }

    result, err := svc.Scan(input)

    if err != nil {
        log.Fatalf("Got error calling Scan: %s", err)
    }

    // Initialize a slice to store the records

    for _, item := range result.Items {
        var record Datum
        err = dynamodbattribute.UnmarshalMap(item, &record)
        if err != nil {
            log.Fatalf("Failed to unmarshal Record: %v", err)
        }
        json.NewEncoder(w).Encode(record)
    }

    // Respond with JSON-encoded records
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
}


func searchHandler(w http.ResponseWriter, r *http.Request) {
    sess := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))
    
    ticker := r.URL.Query().Get("ticker")
    datetime := r.URL.Query().Get("datetime")
    


    // Create DynamoDB client
    svc := dynamodb.New(sess)

    key := map[string]*dynamodb.AttributeValue{
        "Ticker": {
            S: aws.String(ticker),
        },
        "DateTime": {
            N: aws.String(datetime),
        },
    }
    
    result, err := svc.GetItem(&dynamodb.GetItemInput{
        TableName: aws.String("hrose3"),
        Key: key,
    })
    
    if err != nil {
        
        fmt.Fprintf(w, "Got error calling GetItem: %s", err)
        return
    }
    found := 1
    if result.Item == nil {
        found = 0
        fmt.Fprintf(w, "Could not find the item")
    }
    
    item := Datum{}
    err = dynamodbattribute.UnmarshalMap(result.Item, &item)
    if err != nil {
        fmt.Fprintf(w, "Failed to unmarshal Record: %v", err)
        return
    }
    
    if found == 1{
        fmt.Fprintf(w, "Found item:\n")
        fmt.Fprintf(w, "\tTicker:  %s\n", item.Ticker)
        fmt.Fprintf(w, "\tDateTime: %d\n", item.DateTime)
        fmt.Fprintf(w, "\tPrice: %f\n", item.Price)
    }
    w.WriteHeader(http.StatusOK)

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