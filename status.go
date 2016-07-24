package main

import(
    "fmt"
)


/** == Functions == */

/**
 * getStatus()
 */
func getStatus() string {
    fmt.Println("Status OK")
    return string(`{"status":"OK"}`)
}
