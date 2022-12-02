package api

import(
    //"io"
    //"fmt"
    "strings"
    "testing"
    "net/http"
)

//
type tRequest struct {
    Service string
    Method string
    Path string
    Body string
    File string
}

//
var tCreateBuckets = []tRequest{
    {"Base", "POST", "", `{"type":"Bucket","attributes":{"name":"wo-versioning","versioning":false}}`, ""},
    {"Base", "POST", "", `{"type":"Bucket","attributes":{"name":"versioning","versioning":true}}`, ""},
    {"RegisterOfMedicines", "POST", "", `{"type":"DossierChangeRequest"}`, ""},
}

//
func TestCreateBucket(t *testing.T) {
    client := &http.Client{}

    for _, tReq := range tCreateBuckets {
        reader := strings.NewReader(tReq.Body)
        req, err := http.NewRequest("POST", "http://localhost:9101/", reader)
        if err != nil {
            t.Error("For:", "NewRequest", "expected:", "error", "got:", err)
            continue
        }

        req.Header.Add("X-Woodchuck-Service", tReq.Service)
        req.Header.Add("Content-Type", "application/json")

        res, err := client.Do(req)
        if err != nil {
            t.Error("For:", "Do", "expected:", "error", "got:", err)
            continue
        }
        defer res.Body.Close()

        if res.StatusCode != 201 {
            t.Error("For:", "response status code", "expected:", "201", "got:", res.StatusCode)
            continue
        }

        //body, err := io.ReadAll(res.Body)
    }
}

//
var TestCreateObjects = []tRequest{
    {"Base", "POST", "", "", ""},
}

/*
func TestCreateObjects(t *testing.T) {
    client := &http.Client{}

    for _, tReq := range tCreateObjects {
    }
}*/

//
//func TestListBucket(t *testing.t) {
//}

//
///func TestDeleteBucket(t *testing.T) {
//}
