package main

import(
    "encoding/json"
    "fmt"
    "log"
    "runtime"
    "github.com/nats-io/nats"
    "gopkg.in/mgo.v2"
)


type MainfluxMessage struct {
    Method string `json: "method"`
    Id string `json: "id"`
    Body map[string]interface{} `json: "body"`
}

/**
 * MongoDB Globals
 */
type MongoConn struct {
    session *mgo.Session
    dColl *mgo.Collection
    cColl *mgo.Collection
}

var mc MongoConn

/**
 * main()
 */
func main() {
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
    channelMongo := mgoSession.DB("test").C("channels")

    /** Set-up globals */
    mc.session = mgoSession
    mc.dColl = deviceMongo
    mc.cColl = channelMongo

    /**
     * NATS
     */
    nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
        log.Fatalf("Can't connect: %v\n", err)
	}

    // Req-Reply
    nc.Subscribe("core_in", func(msg *nats.Msg) {
        var mfMsg MainfluxMessage

        log.Println(msg.Subject, string(msg.Data))

        // Unmarshal the message recieved from NATS
        err := json.Unmarshal(msg.Data, &mfMsg)
        if err != nil {
		      fmt.Println("error:", err)
        }

        fmt.Println(mfMsg)

        var res string
        switch mfMsg.Method {
            case "getStatus":
                res = getStatus()
            case "createDevice":
                res = createDevice(mfMsg.Body)
            case "getDevices":
                res = getDevices()
            case "getDevice":
                res = getDevice(mfMsg.Id)
            case "updateDevice":
                res = updateDevice(mfMsg.Id, mfMsg.Body)
            case "deleteDevice":
                res = deleteDevice(mfMsg.Id)
            default:
                fmt.Println("error: Unknown method!")
        }

        fmt.Println(res)
        nc.Publish(msg.Reply, []byte(res))
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
