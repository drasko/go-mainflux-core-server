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

func validateJsonSchema(b map[string]interface{}) bool {
    schemaLoader := gojsonschema.NewReferenceLoader("file:///home/drasko/mainflux/go-mainflux-core-server/schema/deviceSchema.json")
    bodyLoader := gojsonschema.NewGoLoader(b)
    result, err := gojsonschema.Validate(schemaLoader, bodyLoader)
    if err != nil {
        log.Print(err.Error())
    }

    if result.Valid() {
        fmt.Printf("The document is valid\n")
        return true
    } else {
        fmt.Printf("The document is not valid. see errors :\n")
        for _, desc := range result.Errors() {
            fmt.Printf("- %s\n", desc)
        }
        return false
    }
}

/** == Functions == */
/**
 * createDevice ()
 */
func createDevice(b map[string]interface{}) string {
    if validateJsonSchema(b) != true {
        println("Invalid schema")
    }

    // Turn map into a JSON to put it in the Device struct later
    j, err := json.Marshal(&b)
    if err != nil {
        fmt.Println(err)
    }

    // Set up defaults and pick up new values from user-provided JSON
    d := Device{Id: "Some Id", Name: "Some Name"}
    json.Unmarshal(j, &d)

    // Creating UUID Version 4
    uuid := uuid.NewV4()
    fmt.Println(uuid.String())

    d.Id = uuid.String()

    // Insert Device
    erri := mc.dColl.Insert(d)
	if erri != nil {
        println("CANNOT INSERT")
		panic(erri)
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
        log.Print(err)
    }

    r, err := json.Marshal(results)
    if err != nil {
        fmt.Println("error:", err)
    }
    return string(r)
}

/**
 * getDevice()
 */
func getDevice(id string) string {
    result := Device{}
    err := mc.dColl.Find(bson.M{"Id": id}).One(&result)
    if err != nil {
        log.Print(err)
    }

    r, err := json.Marshal(result)
    if err != nil {
        fmt.Println("error:", err)
    }
    return string(r)
}

/**
 * updateDevice()
 */
func updateDevice(id string, b map[string]interface{}) string {
    // Validate JSON schema user provided
    if validateJsonSchema(b) != true {
        println("Invalid schema")
    }

    // Check if someone is trying to change "id" key
    // and protect us from this
    if _, ok := b["id"]; ok {
        println("Error: can not change device ID")
    }

    colQuerier := bson.M{"id": id}
	change := bson.M{"$set": b}
    err := mc.dColl.Update(colQuerier, change)
    if err != nil {
        log.Print(err)
    }

    return string(`{"status":"updated"}`)
}

/**
 * deleteDevice()
 */
func deleteDevice(id string) string {
    err := mc.dColl.Remove(bson.M{"id": id})
    if err != nil {
        log.Print(err)
    }

    return string(`{"status":"deleted"}`)
}
