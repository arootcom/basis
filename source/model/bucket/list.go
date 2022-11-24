package bucket

import (
    "fmt"
    "context"

    "filestorage/instance/storage"
    "filestorage/instance/log"
)

//
type ListBucket struct {
    Items []*Bucket          `json:"items"`
    Total int                `json:total`
}

// TODO: Постраничный вывод с лимитом записей на странице
// TODO: Фильтрация
func GetListBucket() (*ListBucket, error) {
    log.Debug("start", "bucket.GetListBucket:")

    client := storage.GetInstance()
    ctx := context.Background()

    buckets, err := client.ListBuckets(ctx)
    if err != nil {
        log.Error("error", "bucket.GetListBucket:", err)
        return nil, err
    }

    list := new(ListBucket)
    list.Total = 0

    for _, item := range buckets {
        b, err := newByBacketInfo(&item)
        if err != nil {
            log.Error("error", "bucket.GetListBucket:", err)
            return nil, err
        }
        list.Items = append(list.Items, b)
        log.Debug("append", "bucket.GetListBucket:", fmt.Sprintf("%+v", b))
        list.Total++
    }

    log.Debug("success", "bucket.GetListBucket:", fmt.Sprintf("%+v", list))
    return list, nil
}

