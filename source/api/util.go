package api

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "encoding/json"

    "filestorage/view"
    "filestorage/instance/log"
)

//
type ErrorBody struct {
    Code string         `json:"code"`
    Note string         `json:"note"`
}

//
type Model interface {}

//
func toJSON(res http.ResponseWriter, model Model, status int) {
    jsonb, err := json.Marshal(model)
    if err != nil {
        panic(err)
    }

    res.Header().Set("Content-Type", "application/json")
    res.WriteHeader(status)

    _, err = res.Write([]byte(jsonb))
    if err != nil {
        panic(err)
    }
    return
}

//
func toError(res http.ResponseWriter, code string, note error, status int) {
    errb := ErrorBody{
        Code: code,
        Note: fmt.Sprintf("%s", note),
    }
    log.Error("error", errb)
    toJSON(res, errb, status)
}

//
func bodyToJSON(req *http.Request, model Model) error {
    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
        log.Error("error", err)
        return err
    }
    log.Info("body", string(body))

    err = json.Unmarshal(body, model)
    if err != nil {
        log.Error("error", err)
        return err
    }

    return nil
}

//
func fromBody(req *http.Request) (*view.Custom, error) {
    log.Debug("start", "fomBody: req = ", fmt.Sprintf("%+v", req))

    custom, err := view.NewCustom()
    if err != nil {
        log.Error("error", "fomBody:", err)
        return nil, err
    }

    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
        log.Error("error", "fomBody:", err)
        return nil, err
    }
    log.Info("body", "fomBody:", string(body))

    err = json.Unmarshal(body, custom)
    if err != nil {
        log.Error("error", "fomBody:", err)
        return nil, err
    }

    log.Debug("success", "fomBody:", fmt.Sprintf("%+v", custom))
    return custom, nil
}


