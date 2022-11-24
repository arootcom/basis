package bucket

import(
    "testing"
)

//
type tBucket struct {
    name string
    versioning bool
    tags map[string]string
}

// 
var buckets = []tBucket{
    {"wo-versioning", false, map[string]string{"Service": "Base", "Type": "Bucket"}},
    {"versioning", true, map[string]string{"Service": "Base", "Type": "Bucket"}},
}

//
func TestCreate(t *testing.T) {
    for _, tbucket := range buckets {
        bucket, _ := NewCreateBucket()
        bucket.Name = tbucket.name
        bucket.Versioning = tbucket.versioning
        bucket.Tags = tbucket.tags

        exists := IsExistsBucketByName(tbucket.name) // test func
        if exists {
            t.Error("For:", tbucket.name, "expected exists:", false, "got:", exists)
        }

        err := bucket.CreateInStorage() // test func
        if err != nil {
            t.Error("For:", tbucket.name, "expected:", "create bucket", "got:", err)
        }
    }
}

//
func TestList(t *testing.T) {
    list, err := GetListBucket()
    if err != nil {
        panic(err)
    }

    ttotal := len(buckets)
    if list.Total != ttotal {
        t.Error("For:", "Bucket total", "expected:", ttotal, "got:", list.Total)
    }

    for _, tbucket := range buckets {
        var exists = false

        for _, bucket := range list.Items {
            if bucket.Name == tbucket.name {
                checkBucket(&tbucket, bucket, t)
                exists = true
                break
            }
        }

        if !exists {
            t.Error("For:", "Bucket", "expected:", tbucket.name, "got:", "none")
        }
    }
}

func TestDelete(t *testing.T) {
    for _, tbucket := range buckets {
        bucket, err := GetBucketByName(tbucket.name)
        if err != nil {
            t.Error("For bucket:", tbucket.name, "Expected:", "GetBucketByName", "Got:", err)
        } else {
            checkBucket(&tbucket, bucket, t)

            del, err := bucket.DeleteBucket()
            if err != nil {
                t.Error("For bucket:", tbucket.name, "Expected:", "DeleteBucket", "Got:", err)
            } else if !del {
                t.Error("For:", tbucket.name, "expected:", true, "got:", del)
            }

            exists := IsExistsBucketByName(tbucket.name) // test func
            if exists {
                t.Error("For:", tbucket.name, "expected:", false, "got:", exists)
            }
        }
    }
}

func checkBucket (tbucket *tBucket, bucket *Bucket, t *testing.T) {
    if tbucket.name != bucket.Name {
        t.Error("For bucket:", tbucket.name, "Expected name:", tbucket.name, "Got name:", bucket.Name)
    }

    if tbucket.versioning != bucket.Versioning {
        t.Error("For bucket:", tbucket.name, "expected:", tbucket.versioning, "got:", bucket.Versioning)
    }

    for ttag, tvalue := range tbucket.tags {
        value, exists := bucket.Tags[ttag]
        if !exists {
            t.Error("For:", tbucket.name, "expected tag:", ttag, "got tag:", "none")
        } else if value != tvalue {
            t.Error("For:", tbucket.name, "expected tag value:", tvalue, "got:", value)
        }
    }
}

