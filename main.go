package mast

func main() {
	var bopts BuilderOptions
	//grpc := builder.grpc(bopts)
	c := bopts.rpcclient()
	s := bopts.rpcServer()
	// gc := grpc.BuildClient()
	// gs := grpc.BuildServer()

	// gs.Start()
	// gc.Call()

	s.Register()
	s.Accept()
	c.Call()
}
