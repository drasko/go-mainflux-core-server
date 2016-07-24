package main

import(
    "encoding/json"
    "fmt"
    "log"

    "github.com/satori/go.uuid"
    "gopkg.in/mgo.v2/bson"
        "github.com/xeipuuv/gojsonschema"
)

type Device struct {
    Id string   `json: "id"`
    Name string `json: "name"`
}

/** == Functions == */
/**
 * createDevice ()
 */
func createDevice(b map[string]interface{}) string {
    // Validate JSON schema
    schemaLoader := gojsonschema.NewReferenceLoader("file:///home/drasko/mainflux/go-mainflux-core-server/schema.json")
    bodyLoader := gojsonschema.NewGoLoader(b)
    result, err := gojsonschema.Validate(schemaLoader, bodyLoader)
    if err != nil {
        panic(err.Error())
    }

    if result.Valid() {
        fmt.Printf("The document is valid\n")
    } else {
        fmt.Printf("The document is not valid. see errors :\n")
        for _, desc := range result.Errors() {
            fmt.Printf("- %s\n", desc)
        }
    }

    j, err := json.Marshal(&b)
    if err != nil {
        fmt.Println(err)
    }

    d := Device{Id: "Some Id", Name: "Some Name"}
    json.Unmarshal(j, &d)

    // Creating UUID Version 4
    uuid := uuid.NewV4()
    fmt.Println(uuid.String())

    //d["id"] = uuid.String()

    // Insert Device
    err := mc.dColl.Insert(d)
	if err != nil {
		panic(err)
	}

    return "Created Device req.deviceId"
}

/**
 * getDevices()
 */
func getDevices() string {
    results := []Device{}
    err := mc.dColl.Find(nil).All(&results)
    if err != nil {
        log.Fatal(err)
    }

    b, err := json.Marshal(results)
    if err != nil {
        fmt.Println("error:", err)
    }
    return string(b)
}

/**
 * getDevice()
 */
func getDevice(id string) string {
    result := Device{}
    err := mc.dColl.Find(bson.M{"Id": id}).One(&result)
    if err != nil {
        log.Fatal(err)
    }

    b, err := json.Marshal(result)
    if err != nil {
        fmt.Println("error:", err)
    }
    return string(b)
}

/**
 * updateDevice()
 */
func updateDevice(id string, b map[string]interface{}) string {
    colQuerier := bson.M{"Id": id}
	change := bson.M{"$set": b}
    err := mc.dColl.Update(colQuerier, change)
    if err != nil {
        log.Fatal(err)
    }

    b, err := json.Marshal(err.Error())
    if err != nil {
        fmt.Println("error:", err)
    }
    return string(b)
}

/**
 * deleteDevice()
 */
func deleteDevice(id string) string {
    err := mc.dColl.Remove(bson.M{"Id": id})
    if err != nil {
        log.Fatal(err)
    }

    b, err := json.Marshal(err.Error())
    if err != nil {
        fmt.Println("error:", err)
    }
    return string(b)
}
