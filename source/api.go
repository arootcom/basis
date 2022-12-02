package main

import (
    "log"
    "woodchuck/api"
)

func main() {
    log.Println("Start...")

    server := api.New()

    err := server.Run("localhost:9101")
    if err != nil {
        log.Fatalln("Failed to run service: ", err)
    }

    log.Println("Stop...")
}
