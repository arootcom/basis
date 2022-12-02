package bucket

import (
    "fmt"
    "errors"
    "strings"
    "context"
    "github.com/minio/minio-go/v7"
    "github.com/minio/minio-go/v7/pkg/tags"

    "woodchuck/instance/storage"
    "woodchuck/instance/log"
)

//
type CreateBucket struct {
    Name string                     `json:"name"`
    Versioning bool                 `json:"versioning"`
    Tags map[string]string          `json:tags`
}

//
func NewCreateBucket() (*CreateBucket, error) {
    return &CreateBucket{
        Tags: make(map[string]string),
    }, nil
}

//
func (b *CreateBucket) SetService(value string) error {
    if value == "" {
        return errors.New("Service value empty")
    }

    b.Tags["Service"] = value
    return nil
}

//
func (b *CreateBucket) SetType(value string) error {
    if value == "" {
        return errors.New("Type value empty")
    }

    b.Tags["Type"] = value
    return nil
}

//
func (b *CreateBucket) SetTag(key string, value string) error {
    if key == "" {
        return errors.New("Tag key empty")
    }

    if value == "" {
        return errors.New("Tag value empty")
    }

    key = strings.Title(key)

    b.Tags[key] = value
    return nil
}

//
func (b *CreateBucket) GetLocation () string {
    return fmt.Sprintf("/%s", b.Name)
}

//
func (b *CreateBucket) ValidationCreateBucket() error {
    if b.Name == "" {
        err := errors.New("The bucket name is not defined")
        log.Error("error", "ValidationCreateBucket: ", b.Name, ", error:", err)
        return err
    }
    return nil
}

//
func (b *CreateBucket) CreateInStorage() error {
    log.Debug("start", "CreateInStorage:", fmt.Sprintf("%+v", b))

    client := storage.GetInstance()
    ctx := context.Background()

    err := client.MakeBucket(ctx, b.Name, minio.MakeBucketOptions{Region: "us-east-1"})
    if err != nil {
        log.Error("error", "CreateInStorage:", b.Name, ", error:", err)
        return err
    }
    log.Debug("create", "CreateInStorage:", b.Name)

    if b.Versioning {
        err = client.EnableVersioning(ctx, b.Name)
        if err != nil {
            log.Error("error", "CreateInStorage:", b.Name, ", error:", err)
            return err
        }
        log.Debug("enable", "CreateInStorage:", b.Name, "EnableVersioning")
    }

    if b.Tags != nil {
        tags, err := tags.NewTags(b.Tags, false)
        if err != nil {
            log.Error("error", "CreateInStorage:", b.Name, ", error:", err)
            return err
        }

        err = client.SetBucketTagging(ctx, b.Name, tags)
        if err != nil {
            log.Error("error", "CreateInStorage:", b.Name, ", error:", err)
            return err
        }
        log.Debug("set", "CreateInStorage:", b.Name, "SetBucketTagging:", tags)
    }

    log.Debug("success", "CreateInStorage:", b.Name, ", versioning:", b.Versioning, ", tags:", b.Tags)
    return nil
}

