package api

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "encoding/json"

    "woodchuck/view"
    "woodchuck/instance/log"
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
func bodyToCustom(req *http.Request) (*view.Custom, error) {
    log.Debug("start", "bodyToCustom: req = ", fmt.Sprintf("%+v", req))

    custom, err := view.NewCustom()
    if err != nil {
        log.Error("error", "bodyToCustom:", err)
        return nil, err
    }

    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
        log.Error("error", "bodyToCustom:", err)
        return nil, err
    }
    log.Info("body", "bodyToCustom:", string(body))

    err = json.Unmarshal(body, custom)
    if err != nil {
        log.Error("error", "bodyToCustom:", err)
        return nil, err
    }

    log.Debug("success", "bodyToCustom:", fmt.Sprintf("%+v", custom))
    return custom, nil
}

//
func stringToCustom(str string) (*view.Custom, error) {
    log.Debug("start", "stringToCustom: req = ", str)

    custom, err := view.NewCustom()
    if err != nil {
        log.Error("error", "stringToCustom:", err)
        return nil, err
    }

    err = json.Unmarshal([]byte(str), custom)
    if err != nil {
        log.Error("error", "stringToCustom:", err)
        return nil, err
    }

    log.Debug("success", "stringToCustom: custom =", fmt.Sprintf("%+v", custom))
    return custom, nil
}


