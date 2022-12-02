package object

import (
    "fmt"
    "errors"
    "context"
    "github.com/minio/minio-go/v7"

    "woodchuck/instance/storage"
    "woodchuck/instance/log"
)

//
type ListObject struct {
    Items []*Object      `json:"items"`
    Total int64         `json:"total"`
}

//
// TODO: Доработать постраничный вывод c limit и after
func GetListObjectByBucket(bucket string, prefix string) (*ListObject, error) {
    client := storage.GetInstance()
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    listCh := client.ListObjects(ctx, bucket, minio.ListObjectsOptions{
        Prefix: prefix,
        Recursive: true,
        WithMetadata: true,
    })

    var objects = new(ListObject)
    objects.Total = 0

    for item := range listCh {
        if item.Err != nil {
            return nil, item.Err
        }

        object, err := NewObject(bucket, item.Key)
        if err != nil {
            return nil, err
        }

        objects.Items = append(objects.Items, object)
        objects.Total++
    }

    return objects, nil
}

// List of object versions
// TODO: Доработать постраничный вывод c limit и after
func GetListObjectVersion(bucket string, key string) (*ListObject, error) {
    log.Debug("start", "GetListObjectVersion: bucket =", bucket, "key =", key)

    client := storage.GetInstance()
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    listCh := client.ListObjects(ctx, bucket, minio.ListObjectsOptions{
        Prefix: key,
        Recursive: true,
        WithVersions: true,
        WithMetadata: true,
    })

    objects := new(ListObject)
    objects.Total = 0

    for item := range listCh {
        if item.Err != nil {
            log.Error("error", "GetListObjectVersion(", bucket, ", ", key, "), error:", item.Err)
            return nil, item.Err
        }

        object, err := NewObjectByVersion(bucket, item.Key, item.VersionID)
        if err != nil {
            log.Error("error", "GetListObjectVersion(", bucket, ", ", key, "), error:", err)
            return nil, err
        }

        objects.Items = append(objects.Items, object)
        objects.Total++
    }

    log.Debug("success", "GetListObjectVersion:", fmt.Sprintf("%+v", objects))
    return objects, nil
}

//
func (o *ListObject) DeleteAllObject() (bool, error) {
    log.Debug("start", "DeleteAllObject:")

    for _, object := range o.Items {
        deleted, err := object.DeleteObject()
        if err != nil {
            log.Error("error", "DeleteAllObject:", object.GetObjectPath(), ", ", err)
            return false, err
        } else if !deleted {
            err := errors.New(fmt.Sprintf(object.GetObjectPath(), ", not deleted"))
            log.Error("error", "DeleteAllObject:", err)
            return false, err
        }
    }

    log.Debug("success", "DeleteAllObject:")
    return true, nil
}
