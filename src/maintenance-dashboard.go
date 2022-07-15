package main

import (
        "fmt"
        "io"
        "log"
        "net/http"
        // "encoding/json"
        "os"

        // "github.com/prometheus/client_golang/prometheus"
        // "github.com/prometheus/client_golang/prometheus/promhttp"
)

func makeApiCall() {
        NewReleasesApiKey := os.Args[1]
        client := http.Client{}
        url := "https://api.newreleases.io/v1/projects"
        req , err := http.NewRequest("GET", url, nil)
        if err != nil {
                log.Fatalln(err)
        }
        // req.Header.Set("X-Key", NewReleasesApiKey)

        req.Header = http.Header{
                "Content-Type": {"application/json"},
                "X-Key": {NewReleasesApiKey},
        }

        res, err := client.Do(req)
        if err != nil {
                log.Fatalln(err)
        }
        body, err := io.ReadAll(res.Body)
        res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", body)
}

func main() {
        makeApiCall()
        // http.Handle("/metrics", promhttp.Handler())
        // log.Fatal(http.ListenAndServe(":2112", nil))
}
