package view

import (
    "fmt"
    //"reflect"
    //"errors"

    "filestorage/model/bucket"
    "filestorage/instance/log"
    //"filestorage/instance/views"
)

//
type ListCustom struct {
    Items []*Custom          `json:"items"`
}

// 
func GetListBucket() (listCustom *ListCustom, err error) {
    log.Debug("start", "view.GetListBucket:")

    listCustom = new(ListCustom)
    listBucket, err := bucket.GetListBucket()
    if err != nil {
        log.Error("error", "view.GetListBucket:", err)
        return nil, err
    }

    for _, item := range listBucket.Items {
        custom, err := NewCustomByBucket(item)
        if err != nil {
            break
        }
        listCustom.Items = append(listCustom.Items, custom)
    }

    if err != nil {
        log.Error("error", "view.GetListBucket:", err)
        return nil, err
    }

    log.Debug("success", "view.GetListBucket:", fmt.Sprintf("%+v", listCustom))
    return listCustom, nil
}

