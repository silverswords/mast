package mast

import (
	"errors"
	"testing"
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
func TestBuild(t *testing.T) {
	// testrpc(t)
}

// func testrpc(t *testing.T) {
// 	mast := &Mast{BuilderOptions: rpc.defaultRPCBuildOptions()}
// 	mast.BuilderOptions.rcvrs["Arith"] = new(Arith)
// 	mast.BuildRPCServer()

// 	args := &Args{7, 8}
// 	var reply int
// 	err := mast.BuildRPCClient().Call("Arith.Multiply", args, &reply)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	if args.B*args.A != reply {
// 		t.Errorf("Arith: %d*%d=%d", args.A, args.B, reply)
// 	}
// 	t.Logf("Arith: %d*%d=%d", args.A, args.B, reply)
// }
