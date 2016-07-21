package main

import(
    "encoding/json"
    "fmt"
    "log"
)

type Device struct {
    DeviceId string   `json: "deviceId"`
    DeviceName string `json: "deviceName"`
    DeviceType string `json: "deviceType"`
}

/** == Functions == */
/**
 * createDevice ()
 */
func createDevice(d Device) string {
    fmt.Println("function f parameter:", d)

    // Insert Datas
    err := mc.dColl.Insert(d)
	if err != nil {
		panic(err)
	}

    return "Created Device req.deviceId"
}

/**
 * getDevices()
 */
func getDevices(req Device) string {
    fmt.Println("function g parameter:", req)

    result := []Device{}
    err := mc.dColl.Find(nil).All(&result)
    if err != nil {
        log.Fatal(err)
    }

    b, err := json.Marshal(result)
    if err != nil {
        fmt.Println("error:", err)
    }
    return string(b)
}
