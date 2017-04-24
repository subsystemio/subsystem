# Subsystem

`subsystem.json` describes a subsystem.
SubSystem is used to manage connections to other SubSystems and the Manager.

```go
s := SubSystem.New()

s.Register("http://localhost:8080/v1")

s.StartHeartbeat()
``` 

The SubSystem is registered, and then finally a continous 5-second heart beat with the server keeps it alive. HTTP or RPC servers can be used to keep the server alive if the HeartBeat is run concurrently.

## Example subsystem.json
```json
{
	"name": "Hello",
	"version": "0.1",
	"repository": {
    	"url": "github.com/subsystemio/subsystem"
  	},
	"hash": {
		"source": "123",
		"build": "123"
	},
	"api": ["SayHello"],
	"requirements": [{
		"name": "World",
		"version": "0.1"
	}],
	"limit": 1
}
```

API is registered with Manager, and made available to all connected Subsystems.
