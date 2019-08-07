# mast

[![Go Report Card](https://goreportcard.com/badge/github.com/silverswords/mast)](https://goreportcard.com/report/github.com/silverswords/mast)

mast is a builder for rpc client and server.By use same options to reduce lines of code.

## Usage

## RPC
### Server
```go
package main

import (
    "errors"
    "github.com/silverswords/mast"
)

type Args struct { 
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}


func main() {
	m := &mast.Mast{BuilderOptions: defaultRPCBuildOptions()}
	m.BuilderOptions.rcvrs["Arith"] = new(Arith)
	m.BuildServer("RPC")
}
```

### Client.go
```go
package main

import (
    "log"
    "github.com/silverswords/mast"
)

type Args struct { 
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

func main(){
	m := &mast.Mast{BuilderOptions: defaultRPCBuildOptions()}
	args := &Args{7, 8}
	var reply int
    err := mast.BuildClient("RPC").(*rpc.Client).Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal(err)
    }
}
```