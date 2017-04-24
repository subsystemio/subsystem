# subsystem.io
Tool and libs for easy deployment of robot subsytems.

## Highly work in progress and still in research phase for best functionality.

# Subsystem

`subsystem.json` describes a subsystem.
Most important are the 'requirements': currently running subsystems; and 'api' the RPC usable functions.

The subsystem 'Hello' for example will provide the RPC function, 'Hello.SayHello'.

Alternatively, when building a subsystem that requires another, the Subsystem's API (Hello/api) can be required by Golang, so as to avoid guessing of correct inputs and outputs.

## Example subsystem.json
```json
{
	"name": "Hello",
	"version": "0.1",
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
