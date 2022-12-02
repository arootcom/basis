package views

import (
    "os"
    "fmt"
    "sync"
    "strings"
    "io/ioutil"
    "encoding/json"

    "woodchuck/instance/log"
)

var once sync.Once
var instance *Services

func GetInstance() *Services {
    once.Do(func() {
        dir := os.Getenv("WOODCHUCK_VIEWS_DIR")
        log.Info("start", "GetInstance:", dir)

        files, err := ioutil.ReadDir(dir)
        if err != nil {
            log.Error("error", "GetInstance:", err)
            panic(err)
        }

        instance = new(Services)
        instance.Items = make(map[string]*Service)

        for _, file := range files {
            if file.IsDir() || !strings.HasSuffix(file.Name(), "json"){
                continue
            }

            filename := fmt.Sprintf("%s%s", dir, file.Name())
            data, err := ioutil.ReadFile(filename)
            if err != nil {
                log.Error("error", "GetInstance:", err)
                panic(err)
            }
            log.Info("read", "GetInstance:", filename)

            srv := new(Service)
            err = json.Unmarshal(data, srv)
            if err != nil {
                log.Error("error", "GetInstance:", err)
                panic(err)
            }
            instance.Items[srv.Service] = srv
        }
        log.Info("success", "GetInstance:", fmt.Sprintf("%+v", instance))
    })
    return instance
}

