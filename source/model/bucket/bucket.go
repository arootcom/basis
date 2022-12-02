package bucket

import (
    "fmt"
    "errors"
    "context"
    "strings"
    "github.com/minio/minio-go/v7"
    //"github.com/minio/minio-go/v7/pkg/tags"

    "woodchuck/instance/storage"
    "woodchuck/instance/log"
)

//
type Bucket struct {
    Name string                     `json:"name"`
    Versioning bool                 `json:"versioning"`
    Tags map[string]string          `json:tags`
}

//
func newByBacketInfo(b *minio.BucketInfo) (*Bucket, error) {
    client := storage.GetInstance()
    ctx := context.Background()

    versioning, err := getVersioningByBucketName(b.Name)
    if err != nil {
        log.Error("error", "newByBacketInfo:", err)
        return nil, err
    }

    tags, err := client.GetBucketTagging(ctx, b.Name)
    if err != nil {
        log.Error("error", "newByBacketInfo:", err)
        return nil, err
    }

    return &Bucket{
        Name: b.Name,
        Versioning: versioning,
        Tags: tags.ToMap(),
    }, nil
}

//
func getVersioningByBucketName(name string) (bool, error) {
    client := storage.GetInstance()
    ctx := context.Background()

    config, err := client.GetBucketVersioning(ctx, name)
    if err != nil {
        log.Error("error", "getVersioningByBucketName(", name, "), error:", err)
        return false, err
    }

    if config.Status == "Enabled" {
        return true, nil
    }
    return false, nil
}

//
func GetBucketByName(name string) (*Bucket, error) {
    client := storage.GetInstance()
    ctx := context.Background()

    versioning, err := getVersioningByBucketName(name)
    if err != nil {
        log.Error("error", "GetBucketByName(", name, "), error:", err)
        return nil, err
    }

    tags, err := client.GetBucketTagging(ctx, name)
    if err != nil {
        log.Error("error", "GetBucketByName(", name, "), error:", err)
        return nil, err
    }

    return &Bucket{
        Name: name,
        Versioning: versioning,
        Tags: tags.ToMap(),
    }, nil
}

// 
func IsExistsBucketByName(name string) bool {
    client := storage.GetInstance()
    ctx := context.Background()

    exists, err := client.BucketExists(ctx, name)
    if err != nil {
        log.Error("error", "IsExistsBucketByName(", name, "), error:", err)
        return false
    } else if !exists {
        return false
    }
    return true
}

//
func (b *Bucket) GetService() (service string, err error) {
    service, exists := b.Tags["Service"]
    if !exists {
        err = errors.New("Service not defined")
        return "", err
    }
    return service, nil
}

//
func (b *Bucket) GetType() (t string, err error) {
    t, exists := b.Tags["Type"]
    if !exists {
        err = errors.New("Type not defined")
        return "", err
    }
    return t, nil
}

//
func (b *Bucket) GetTagValueByKey(key string) (string, error) {
    key = strings.Title(key)

    value, exists := b.Tags[key]
    if !exists {
        err := errors.New(fmt.Sprintf("%s not defined", key))
        return "", err
    }

    return value, nil
}

//
func (b *Bucket) DeleteBucket() (bool, error) {
    client := storage.GetInstance()
    ctx := context.Background()

    err := client.RemoveBucket(ctx, b.Name)
    if err != nil {
        log.Error("error", "DeleteBucket(", b.Name, "), error:", err)
        return false, err
    }

    log.Info("msg", "DeleteBucket:", b.Name, ", versioning:", b.Versioning, ", tags:", b.Tags)
    return true, nil
}

