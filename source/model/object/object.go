package object

import (
    "io"
    "fmt"
    "regexp"
    "context"
    "io/ioutil"
    "github.com/minio/minio-go/v7"
    "github.com/minio/minio-go/v7/pkg/tags"

    "filestorage/instance/storage"
    "filestorage/instance/log"
)

//
type Object struct {
    Bucket string                   `json:"bucket"`
    Prefix string                   `json:"prefix"`
    Name string                     `json:"name"`
    Key string                      `json:"key"`
    VersionId string                `json:"versionId"`
    ContentType string              `json:"contentType"`
    Size int64                      `json:"size"`
    Utime string                    `json:"utime"`
    Metadata map[string]string      `json:"metadata"`
    Tags map[string]string          `json:tags`
}

// 
func IsExistsObjectByKey(bucket string, key string) bool {
    client := storage.GetInstance()
    ctx := context.Background()

    _, err := client.StatObject(ctx, bucket, key, minio.StatObjectOptions{})
    if err != nil {
        return false
    }

    return true
}

// Get object from storage
func NewObject(bucket string, key string) (*Object, error) {
    client := storage.GetInstance()
    ctx := context.Background()

    stat, err := client.StatObject(ctx, bucket, key, minio.StatObjectOptions{})
    if err != nil {
        log.Error("error", "NewObject(", bucket, ", ", key, "), error:", err)
        return nil, err
    }

    tags, err := client.GetObjectTagging(ctx, bucket, key, minio.GetObjectTaggingOptions{})
    if err != nil {
        log.Error("error", "NewObject(", bucket, ", ", key, "), error:", err)
        return nil, err
    }

    reName := regexp.MustCompile(`[^/]+$`)
    rePrefix := regexp.MustCompile(`.*/`)

    object := Object{
        Bucket: bucket,
        Prefix: rePrefix.FindString(stat.Key),
        Name: reName.FindString(stat.Key),
        Key: stat.Key,
        VersionId: stat.VersionID,
        ContentType: stat.ContentType,
        Size: stat.Size,
        Utime: fmt.Sprintf("%s", stat.LastModified),
        Metadata: stat.UserMetadata,
        Tags: tags.ToMap(),
    }

    log.Debug("success", "NewObject:", object.GetObjectPath(), object)
    return &object, nil
}

// Get object version from storage
func NewObjectByVersion(bucket string, key string, versionId string) (*Object, error) {
    log.Debug("start", "NewObjectByVersion: bucket =", bucket, "key =", key, "versionId =", versionId)

    client := storage.GetInstance()
    ctx := context.Background()

    stat, err := client.StatObject(ctx, bucket, key, minio.StatObjectOptions{VersionID : versionId})
    if err != nil {
        return nil, err
    }

    tags, err := client.GetObjectTagging(ctx, bucket, key, minio.GetObjectTaggingOptions{VersionID : versionId})
    if err != nil {
        log.Error("error", fmt.Sprintf("NewObject(%s, %s, %s)", bucket, key, versionId), "error:", err)
        return nil, err
    }

    reName := regexp.MustCompile(`[^/]+$`)
    rePrefix := regexp.MustCompile(`.*/`)

    object := Object{
        Bucket: bucket,
        Prefix: rePrefix.FindString(stat.Key),
        Name: reName.FindString(stat.Key),
        Key: stat.Key,
        VersionId: stat.VersionID,
        ContentType: stat.ContentType,
        Size: stat.Size,
        Utime: fmt.Sprintf("%s", stat.LastModified),
        Metadata: stat.UserMetadata,
        Tags: tags.ToMap(),
    }

    log.Debug("success", "NewObjectByVersion:", object.GetObjectPath(), object)
    return &object, nil
}

// Get object's data
// TODO: Оптимизировать, что бы данные файла сразу передавались в сторону браузера без буфера
//       Иначе при чтении больших файлов будет расходываться память
func (o *Object) DataObject() ([]uint8, error) {
    log.Debug("start", "DataObject:", o.GetObjectPath())

    client := storage.GetInstance()
    ctx := context.Background()

    reader, err := client.GetObject(ctx, o.Bucket, o.Key, minio.GetObjectOptions{VersionID : o.VersionId})
    if err != nil {
        log.Error("error", "DataObject:", o.GetObjectPath(), "error:", err)
        return nil, err
    }
    defer reader.Close()

    buf, err := ioutil.ReadAll(reader)
    if err != nil {
        log.Error("error", "DataObject:", o.GetObjectPath(), "error:", err)
        return nil, err
    }

    log.Debug("success", "DataObject:", o.GetObjectPath())
    return buf, nil
}

// Get path object
func (o *Object) GetObjectPath() string {
    if o.VersionId == "" {
        return fmt.Sprintf("%s/%s", o.Bucket, o.Key)
    }
    return fmt.Sprintf("%s/%s?versionId=%s", o.Bucket, o.Key, o.VersionId)
}

// Update file object
func (o *Object) UpdateObjectInStorage(reader io.Reader, size int64) error {
    client := storage.GetInstance()
    ctx := context.Background()

    _, err := client.PutObject(ctx, o.Bucket, o.Key, reader, size,
        minio.PutObjectOptions{
            ContentType: o.ContentType,
            UserMetadata: o.Metadata,
        },
    )
    if err != nil {
        log.Error("msg", "UpdateObjectInStorage:", o.GetObjectPath(), ", error:", err)
        return err
    }
    log.Debug("msg", "UpdateObjectInStorage:", o.GetObjectPath())

    if o.Tags != nil {
        tags, err := tags.NewTags(o.Tags, true)
        if err != nil {
            log.Error("error", "UpdateObjectInStorage:", o.GetObjectPath(), ", error:", err)
            return err
        }

        err = client.PutObjectTagging(ctx, o.Bucket, o.Key,  tags,  minio.PutObjectTaggingOptions{})
        if err != nil {
            log.Error("error", "UpdateObjectInStorage:", o.GetObjectPath(), ", error:", err)
            return err
        }
        log.Debug("msg", "UpdateObjectInStorage:", o.GetObjectPath(), "PutObjectTagging:", tags)
    }

    return nil
}

// Update object tags
func (o *Object) UpdateObjectTags(updateTags map[string]string) error {
    client := storage.GetInstance()
    ctx := context.Background()

    log.Debug("msg", "UpdateObjectTags:", o.GetObjectPath(), "update tags:", updateTags)

    for key, value := range updateTags {
        o.Tags[key] = value
    }

    tags, err := tags.NewTags(o.Tags, true)
    if err != nil {
        log.Error("error", "UpdateObjectTags:", o.GetObjectPath(), ", error:", err)
        return err
    }

    err = client.PutObjectTagging(ctx, o.Bucket, o.Key,  tags,  minio.PutObjectTaggingOptions{})
    if err != nil {
        log.Error("error", "UpdateObjectTags:", o.GetObjectPath(), ", error:", err)
        return err
    }

    log.Debug("msg", "UpdateObjectTags:", o.GetObjectPath(), "tags:", o.Tags)
    return nil
}


// Delete object
func (o *Object) DeleteObject() (bool, error) {
    client := storage.GetInstance()
    ctx := context.Background()

    err := client.RemoveObject(ctx, o.Bucket, o.Key, minio.RemoveObjectOptions{
		GovernanceBypass: true,
		VersionID: o.VersionId,
	})
    if err != nil {
        return false, err
    }

    log.Debug("msg", "DeleteObject:", o.GetObjectPath())
    return true, nil
}

