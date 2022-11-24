package views

import(
    "fmt"
    "testing"
)

//
type tCreateBucket struct {
    vService string
    vBase string
    vType string
}

var tBuckets = []tCreateBucket{
    {"Base", "CreateBucket", "Bucket"},
    {"Base", "CreateBucket", "Bucket"},
    {"RegisterOfMedicines", "CreateBucket", "DossierChangeRequest"},
    {"RegisterOfMedicines", "Bucket", "DossierChangeRequest"},
}

//
func TestView(t *testing.T) {
    instance := GetInstance()
    for _, tview := range tBuckets {
        view, err := instance.GetViewByType(tview.vService, tview.vBase, tview.vType)
        if err != nil {
            t.Error("For:", "GetObjectByType", "expected:", tview, "got:", err)
            continue
        }

        fmt.Println("vobj:", view)
    }
}

