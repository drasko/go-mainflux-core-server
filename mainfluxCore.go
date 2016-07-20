package mainfluxcore

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
    cColl *mgo.Collection
}

var mc MongoConn


/**
 * main()
 */
func ServerStart() {

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
