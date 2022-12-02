package log

import(
    "os"
    "fmt"
    "testing"
)

var tLevels = []string{"All", "Debug", "Error", "Info", "None"}

//
func TestLevel(t *testing.T) {
    isLevel := false
    for _, tLevel := range tLevels {
        if tLevel == os.Getenv("WOODCHUCK_LOG_LEVEL") {
            isLevel = true
            break
        }
    }

    if !isLevel {
        t.Error("For:", "WOODCHUCK_LOG_LEVEL", "expected:", tLevels, "got:", os.Getenv("WOODCHUCK_LOG_LEVEL"))
    } else {
        fmt.Println("WOODCHUCK_LOG_LEVEL =", os.Getenv("WOODCHUCK_LOG_LEVEL"))
    }

    Debug("msg", "show", "debug", "level")
    Info("msg", "show", "info", "level")
    Warn("msg", "show", "warn", "level")
    Error("msg", "show", "error", "level")
}
