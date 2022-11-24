package object

import(
    "os"
    "fmt"
    "testing"
    "io/ioutil"
    "crypto/md5"
)

//
type tObject struct {
    bucket string
    prefix string
    name string
    key string
    contentType string
    metadata map[string]string
    tags map[string]string
    versioning bool
    file string
    update tUpdate
}

type tUpdate struct {
    tags map[string]string
    file string
}

//
var tObjects = []tObject{
    {
        "wo-versioning", "", "request.xml", "request.xml", "application/xml",
        map[string]string{"Service":"Base", "Type":"Object"},
        map[string]string{"Status": "Edit"},
        false,
        "../../test/request_v1.xml",
        tUpdate{
            map[string]string{"Status": "Archive"},
            "../../test/request_v2.xml",
        },
    },
    {
        "wo-versioning", "signed/", "request.xml", "signed/request.xml", "application/xml",
        map[string]string{"Service":"Base", "Type":"Object"},
        map[string]string{"Antivirus":"NotVerified", "Signature": "NotVerified"},
        false,
        "../../test/request_v1.xml",
        tUpdate{
            map[string]string{"Antivirus":"Verified", "Signature": "Verified"},
            "../../test/request_v2.xml",
        },
    },
    {
        "versioning", "", "request.xml", "request.xml", "application/xml",
        map[string]string{"Service":"Base", "Type":"Object"},
        map[string]string{"Status": "Edit"},
        true,
        "../../test/request_v1.xml",
        tUpdate{
            map[string]string{"Status": "Archive"},
            "../../test/request_v2.xml",
        },
    },
    {
        "versioning", "signed/", "request.xml", "signed/request.xml", "application/xml",
        map[string]string{"Service":"Base", "Type":"Object"},
        map[string]string{"Antivirus":"NotVerified", "Signature": "NotVerified"},
        true,
        "../../test/request_v1.xml",
        tUpdate{
            map[string]string{"Antivirus":"Verified", "Signature": "Verified"},
            "../../test/request_v2.xml",
        },
    },
}

//
func testObject(t *testing.T, object *Object, tobject *tObject, update *tUpdate) {
    if object.Bucket != tobject.bucket {
        t.Error("For:", object.GetObjectPath(), "expected bucket:", tobject.bucket, "got:", object.Bucket)
    }

    if object.Prefix != tobject.prefix {
        t.Error("For:", object.GetObjectPath(), "expected prefix:", tobject.prefix, "got:", object.Prefix)
    }

    if object.Name != tobject.name {
        t.Error("For:", object.GetObjectPath(), "expected name:", tobject.name, "got:", object.Name)
    }

    if object.Key != tobject.key {
        t.Error("For:", object.GetObjectPath(), "expected key:", tobject.key, "got:", object.Key)
    }

    if object.Key != tobject.key {
        t.Error("For:", object.GetObjectPath(), "expected key:", tobject.key, "got:", object.Key)
    }

    for tkey, tvalue := range tobject.metadata {
        value, exists := object.Metadata[tkey]
        if !exists {
            t.Error("For:", object.GetObjectPath(), "expected metadata:", tkey, "got:", "none")
        } else if ( tvalue != value ) {
            t.Error("For:", object.GetObjectPath(), "expected metadata:", fmt.Sprintf("%s => %s", tkey, tvalue), "got:", value)
        }
    }

    var tags map[string]string
    if update != nil {
        tags = update.tags
    } else {
        tags = tobject.tags
    }

    for tkey, tvalue := range tags {
        value, exists := object.Tags[tkey]
        if !exists {
            t.Error("For:", object.GetObjectPath(), "expected tags:", tkey, "got:", "none")
        } else if value != tvalue {
            t.Error("For:", object.GetObjectPath(), "expected tags:", fmt.Sprintf("%s => %s", tkey, tvalue), "got:", value)
        }
    }

    data, err := object.DataObject()
    if err != nil {
        t.Error("For:", "DataObject", "expected:", object.GetObjectPath(), "got:", err)
        return
    }

    var file string
    if update != nil {
        file = update.file
    } else {
        file = tobject.file
    }

    reader, err := os.Open(file)
    defer reader.Close()
    if err != nil {
        t.Error("For:", "Open file", "expected:", file, "got:", err)
        return
    }

    tdata, err := ioutil.ReadAll(reader)
    if err != nil {
        t.Error("For:", "Test file", "expected:", file, "got:", err)
        return
    }

    dataMd5 := md5.Sum(data)
    tdataMd5 := md5.Sum(tdata)
    if dataMd5 != tdataMd5 {
        t.Error("For:", object.GetObjectPath(), "expected:", fmt.Sprintf("%x", tdataMd5), "got:", fmt.Sprintf("%x", dataMd5))
        return
    }

    return
}

//
func TestCreate(t *testing.T) {
    for _, tcreate := range tObjects {
        create, _ := NewCreateObject()
        create.Bucket = tcreate.bucket
        create.Prefix = tcreate.prefix
        create.Name = tcreate.name
        create.ContentType = tcreate.contentType
        create.Metadata = tcreate.metadata
        create.Tags = tcreate.tags

        err := create.ValidationCreateObject()
        if err != nil {
            t.Error("For:", "CreateObject", "expected:", create.Name, "/", create.GetCreateObjectKey(), "got:", err)
            continue
        }

        key := create.GetCreateObjectKey()
        if key != tcreate.key {
            t.Error("For:", "GetCreateObjectKey", "expected:", tcreate.key, "got:", key)
            continue
        }

        exists := IsExistsObjectByKey(create.Bucket, create.GetCreateObjectKey())
        if exists {
            t.Error("For:", "Upload", "expected:", fmt.Sprintf("upload %s", create.Name), "got:", "already exists")
            continue
        }

        reader, err := os.Open(tcreate.file)
        defer reader.Close()
        if err != nil {
            t.Error("For:", "Upload file", "expected:", tcreate.file, "got:", err)
            continue
        }

        stat, err := reader.Stat()
        if err != nil {
            t.Error("For:", "Stat file", "expected:", tcreate.file, "got:", err)
            continue
        }

        err = create.CreateObjectInStorage(reader, stat.Size())
        if err != nil {
            t.Error("For:", "Stat file", "expected:", tcreate.file, "got:", err)
            continue
        }
    }
}

//
func TestGetFirst(t *testing.T) {
    for _, tobject := range tObjects {
        object, err := NewObject(tobject.bucket, tobject.key)
        if err != nil {
            t.Error("For:", "NewObject", "expected:", fmt.Sprintf("%s/%s", tobject.bucket, tobject.key), "got:", err)
            continue
        }

        testObject(t, object, &tobject, nil)
    }
}

//
func TestUpdate(t *testing.T) {
    for _, tobject := range tObjects {
        object, err := NewObject(tobject.bucket, tobject.key)
        if err != nil {
            t.Error("For:", "NewObject", "expected:", fmt.Sprintf("%s/%s", tobject.bucket, tobject.key), "got:", err)
            continue
        }

        reader, err := os.Open(tobject.update.file)
        defer reader.Close()
        if err != nil {
            t.Error("For:", "Upload update file", "expected:", tobject.update.file, "got:", err)
            continue
        }

        stat, err := reader.Stat()
        if err != nil {
            t.Error("For:", "Stat update file", "expected:", tobject.update.file, "got:", err)
            continue
        }

        err = object.UpdateObjectInStorage(reader, stat.Size())
        if err != nil {
            t.Error("For:", "Update  file", "expected:", tobject.update.file, "got:", err)
            continue
        }

        err = object.UpdateObjectTags(tobject.update.tags)
        if err != nil {
            t.Error("For:", "Update tags", "expected:", tobject.update.tags, "got:", err)
            continue
        }
    }
}

//
func TestGetSecond(t *testing.T) {
    for _, tobject := range tObjects {
        object, err := NewObject(tobject.bucket, tobject.key)
        if err != nil {
            t.Error("For:", "NewObject", "expected:", fmt.Sprintf("%s/%s", tobject.bucket, tobject.key), "got:", err)
            continue
        }

        testObject(t, object, &tobject, &tobject.update)
    }
}

//
func TestVersioning(t *testing.T) {
    for _, tobject := range tObjects {
        objects, err := GetListObjectVersion(tobject.bucket, tobject.key)
        if err != nil {
            t.Error("For:", "GetListObjectVersion", "expected:", fmt.Sprintf("%s/%s", tobject.bucket, tobject.key), "got:", err)
            continue
        }

        if !tobject.versioning && objects.Total != 1 {
            t.Error("For:", "tListObject.Total", "expected:", "1", "got:", objects.Total)
            continue
        } else if tobject.versioning && objects.Total != 2 {
            t.Error("For:", "tListObject.Total", "expected:", "2", "got:", objects.Total)
            continue
        }

        if tobject.versioning {
            count := len(objects.Items) - 1
            for i := count; i >= 0 ; i-- {
                object := objects.Items[i]
                if i == count {
                    testObject(t, object, &tobject, nil)
                } else {
                    testObject(t, object, &tobject, &tobject.update)
                }
            }
        }
    }
}

//
func TestDelete(t *testing.T) {
    for _, tobject := range tObjects {
        exists := IsExistsObjectByKey(tobject.bucket, tobject.key)
        if !exists {
            t.Error("For:", "IsExistsObjectByKey", "expected:", fmt.Sprintf("exists %s", tobject.name), "got:", "not exists")
            continue
        }

        objects, err := GetListObjectVersion(tobject.bucket, tobject.key)
        if err != nil {
            t.Error("For:", "Object versions", "expected:", tobject.key, "got:", err)
            continue
        }

        objects.DeleteAllObject()
    }
}

