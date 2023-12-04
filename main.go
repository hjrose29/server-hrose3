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
    "strings"
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
    r.HandleFunc("/hrose3/rangedSearch", rangedSearchHandler).Methods(http.MethodGet)
	r.PathPrefix("/").HandlerFunc(catchAllHandler).Methods(http.MethodGet)
    r.Use(RequestLoggerMiddleware(r))
    http.ListenAndServe(":8080", r)

}



func statusHandler(w http.ResponseWriter, r *http.Request) {

    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, time.Now().String())

    sess := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))
    svc := dynamodb.New(sess)
    input := &dynamodb.DescribeTableInput{
        TableName: aws.String("hrose3"),
    }
    desc, err := svc.DescribeTable(input)

    if err!= nil{
        log.Fatalf("Could not describe table: %s", err)
    }
    table := desc.Table;
    fmt.Fprintf(w, "\nNumber of items: %d", *table.ItemCount);
    
    
}


//Dumps all contents of DynamoDB table in JSON.
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



//Barebones Search Handler (Accepts ticker and UNIX datetime)
func searchHandler(w http.ResponseWriter, r *http.Request) {
    sess := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))
    

    ticker := r.URL.Query().Get("ticker")
    ticker = strings.ToUpper(ticker)
    datetime := r.URL.Query().Get("datetime")
    
    if(ticker != "TSLA" && ticker != "AAPL" && ticker != "SPY" && ticker != "AMZN"){
        fmt.Fprintf(w, "ticker not valid, please try again, i.e. \"TSLA\", \"AMZN\", \"AAPL\", \"SPY\"")
        return
    }


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

    //If it doesn't find the item, let's user know.
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


//Ranged Handler (Accepts ticker and UNIX time lower and upper bounds)
func rangedSearchHandler(w http.ResponseWriter, r *http.Request) {
    sess := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))
    

    ticker := r.URL.Query().Get("ticker")
    ticker = strings.ToUpper(ticker)
    lower := r.URL.Query().Get("lower")
    upper := r.URL.Query().Get("upper")
    

    //Data validation step to only accept certain tickers.
    if(ticker != "TSLA" && ticker != "AAPL" && ticker != "SPY" && ticker != "AMZN"){
        fmt.Fprintf(w, "ticker not valid, please try again, i.e. \"TSLA\", \"AMZN\", \"AAPL\", \"SPY\"")
        return
    }

    //Turns lower and upper into optional filter parameters.
    if(lower == ""){
        lower = "-1";
    }
    if(upper == ""){
        upper = "2551641525";
    }

    // Create DynamoDB client
    svc := dynamodb.New(sess)

    var KeyConditions = map[string]*dynamodb.Condition{
            "Ticker": {
                ComparisonOperator: aws.String("EQ"),
                AttributeValueList: []*dynamodb.AttributeValue{
                    {
                        S: aws.String(ticker),
                    },
                },
            },
            "DateTime": {
                ComparisonOperator: aws.String("BETWEEN"),
                AttributeValueList: []*dynamodb.AttributeValue{
                    {
                        N: aws.String(lower),                        
                    },
                    {
                        N: aws.String(upper),                        
                    },
                },
            },
    }
    

    result, err := svc.Query(&dynamodb.QueryInput{
        TableName:     aws.String("hrose3"),
        KeyConditions: KeyConditions, 
    })

    if err != nil {
        fmt.Fprintf(w, "Error querying DynamoDB: %v", err)
        return
    }

    var items []Datum  // Assuming Datum is the type you want to marshal
    err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &items)
    if err != nil {
        fmt.Fprintf(w, "Error unmarshaling DynamoDB query result: %v", err)
        return
    }

    // Marshal the items slice into JSON
    jsonResponse, err := json.Marshal(items)
    if err != nil {
        fmt.Fprintf(w, "Error marshaling data to JSON: %v", err)
        return
    }

    // Set content type to JSON
    w.Header().Set("Content-Type", "application/json")
    // Write the JSON response to the response writer
    w.Write(jsonResponse)

    w.WriteHeader(http.StatusOK)

}




//Any other http GET request is handled here
func catchAllHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "404 - Page not found!")
}


//Middleware allowing us to see requests with our logging tool(Solarwind's Loggly).
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