package main

import (
    "fmt"
    mfcore "https://github.com/drasko/go-mainflux-core-server"
)

func main() {
    fmt.Println("Hello world!")
    mfcore.ServerStart()
}
