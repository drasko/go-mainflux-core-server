package main

import(
    "encoding/json"
    "fmt"
    "log"
    "runtime"
    "github.com/nats-io/nats"
    "gopkg.in/mgo.v2"
    //"gopkg.in/mgo.v2/bson"
)


type MainfluxMessage struct {
    Method string `json: "method"`
    Body Device `json: "body"`
}

/**
 * MongoDB Globals
 */
type MongoConn struct {
    session *mgo.Session
    dColl *mgo.Collection
    sColl *mgo.Collection
}

var mc MongoConn

type Device struct {
    DeviceId string   `json: "deviceId"`
    DeviceName string `json: "deviceName"`
    DeviceType string `json: "deviceType"`
}

type Stream struct {

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

/**
 * main()
 */
func main() {

    /** Callback map */
    fncMap := map[string]func(Device) string {
        "createDevice": createDevice,
        "getDevices": getDevices,
    }

    /**
     * MongoDB
     */
    mgoSession, err := mgo.Dial("localhost")
    if err != nil {
            panic(err)
    }
    //defer mgoSession.Close()

    // Optional. Switch the session to a monotonic behavior.
    mgoSession.SetMode(mgo.Monotonic, true)

    deviceMongo := mgoSession.DB("test").C("devices")
    streamMongo := mgoSession.DB("test").C("streams")

    /** Set-up globals */
    mc.session = mgoSession
    mc.dColl = deviceMongo
    mc.sColl = streamMongo

    /**
     * NATS
     */
    nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Can't connect: %v\n", err)
	}

    // Replying
    nc.Subscribe("core_in", func(msg *nats.Msg) {
        var mfMsg MainfluxMessage

        log.Println(msg.Subject, string(msg.Data))

        err := json.Unmarshal(msg.Data, &mfMsg)
        if err != nil {
		    fmt.Println("error:", err)
	    }

        fmt.Println(mfMsg)
        fmt.Printf("%+v", mfMsg)

        f := fncMap[mfMsg.Method]
        res := f(mfMsg.Body)
        fmt.Println(res)
        //nc.Publish(msg.Reply, []byte(res))
    })

	log.Println("Listening on 'core_in'")

    fmt.Println(banner)

    /** Keep mainf() runnig */
	runtime.Goexit()
}

var banner = `
_|      _|            _|                _|_|  _|                      
_|_|  _|_|    _|_|_|      _|_|_|      _|      _|  _|    _|  _|    _|  
_|  _|  _|  _|    _|  _|  _|    _|  _|_|_|_|  _|  _|    _|    _|_|    
_|      _|  _|    _|  _|  _|    _|    _|      _|  _|    _|  _|    _|  
_|      _|    _|_|_|  _|  _|    _|    _|      _|    _|_|_|  _|    _|  
                                                                      
    
                == Industrial IoT System ==
       
                Made with <3 by Mainflux Team
[w] http://mainflux.io
[t] @mainflux
`
