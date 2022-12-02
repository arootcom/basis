package object

import (
    "io"
    "fmt"
    "regexp"
    "errors"
    "context"
    "github.com/minio/minio-go/v7"
    "github.com/minio/minio-go/v7/pkg/tags"

    "woodchuck/instance/storage"
    "woodchuck/instance/log"
)

type CreateObject struct {
    Bucket string                   `json:"bucket"`
    Prefix string                   `json:"prefix"`
    Name string                     `json:"name"`
    ContentType string              `json:"contentType"`
    Metadata map[string]string      `json:"metadata"`
    Tags map[string]string          `json:tags`
}

// New create object
func NewCreateObject() (*CreateObject, error) {
    return &CreateObject{}, nil
}

// Get create object key
func (o *CreateObject) GetCreateObjectKey() string {
    return fmt.Sprintf("%s%s", o.Prefix, o.Name)
}

// Get path object
func (o *CreateObject) GetCreateObjectPath() string {
    if o.Prefix == "" {
        return fmt.Sprintf("%s/%s", o.Bucket, o.Name)
    }
    return fmt.Sprintf("%s/%s%s", o.Bucket, o.Prefix, o.Name)
}

// Validation data for create object
func (o *CreateObject) ValidationCreateObject() error {
    if o.Bucket == "" {
        return errors.New("The bucket name is not defined")
    }

    if o.Prefix != "" {
        prefixRe := regexp.MustCompile(`^/`)
        if prefixRe.MatchString(o.Prefix)  {
            return errors.New("The prefix cannot start with \"/\"")
        }

        prefixRe = regexp.MustCompile(`/$`)
        if !prefixRe.MatchString(o.Prefix)  {
            return errors.New("The prefix must end with \"/\"")
        }
    }

    if o.Name == "" {
        return errors.New("The name object is not defined")
    }

    if o.ContentType == "" {
        return errors.New("The contentType object is not defined")
    }

    return nil
}

// Create object
func (o *CreateObject) CreateObjectInStorage(reader io.Reader, size int64) error {
    client := storage.GetInstance()
    ctx := context.Background()

    object, err := client.PutObject(ctx, o.Bucket, o.GetCreateObjectKey(), reader, size,
        minio.PutObjectOptions{
            ContentType: o.ContentType,
            UserMetadata: o.Metadata,
        },
    )
    if err != nil {
        log.Error("msg", "CreateObjectInStorage:", o.GetCreateObjectPath(), ", error:", err)
        return err
    }
    log.Debug("success", "CreateObjectInStorage:", o.GetCreateObjectPath(), "PutObject:", object)

    if o.Tags != nil {
        tags, err := tags.NewTags(o.Tags, true)
        if err != nil {
            log.Error("error", "CreateObjectInStorage:", o.GetCreateObjectPath(), ", error:", err)
            return err
        }

        err = client.PutObjectTagging(ctx, o.Bucket, o.GetCreateObjectKey(),  tags,  minio.PutObjectTaggingOptions{})
        if err != nil {
            log.Error("error", "CreateObjectInStorage:", o.GetCreateObjectPath(), ", error:", err)
            return err
        }
        log.Debug("success", "CreateObjectInStorage:", o.GetCreateObjectPath(), "PutObjectTagging:", tags)
    }

    return nil
}


