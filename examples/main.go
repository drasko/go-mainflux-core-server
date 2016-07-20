package main

import (
    "fmt"
    mfcore "github.com/drasko/go-mainflux-core-server"
)

func main() {
    fmt.Println("Hello world!")
    mfcore.ServerStart()
}
