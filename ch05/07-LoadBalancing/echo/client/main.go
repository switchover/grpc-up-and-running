package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	ecpb "google.golang.org/grpc/examples/features/proto/echo"
	"google.golang.org/grpc/resolver"
)

//---------------------------------------------------------
// 코드 5-15 부분
//---------------------------------------------------------
const (
	exampleScheme      = "example"
	exampleServiceName = "lb.example.grpc.io"
)

var addrs = []string{"localhost:50051", "localhost:50052"}

type exampleResolverBuilder struct{}

func (*exampleResolverBuilder) Build(target resolver.Target,
	cc resolver.ClientConn,
	opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &exampleResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			exampleServiceName: addrs,
		},
	}
	r.start()
	return r, nil
}

func (*exampleResolverBuilder) Scheme() string { return exampleScheme }

type exampleResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (r *exampleResolver) start() {
	addrStrs := r.addrsStore[r.target.Endpoint]
	addrs := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrs[i] = resolver.Address{Addr: s}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}

func (*exampleResolver) ResolveNow(o resolver.ResolveNowOptions) {}

func (*exampleResolver) Close() {}

func init() {
	resolver.Register(&exampleResolverBuilder{})
}

//---------------------------------------------------------

func callUnaryEcho(c ecpb.EchoClient, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.UnaryEcho(ctx, &ecpb.EchoRequest{Message: message})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	fmt.Println(r.Message)
}

func makeRPCs(cc *grpc.ClientConn, n int) {
	hwc := ecpb.NewEchoClient(cc)
	for i := 0; i < n; i++ {
		callUnaryEcho(hwc, "this is examples/load_balancing")
	}
}

func main() {
	//---------------------------------------------------------
	// 코드 5-16 부분
	//---------------------------------------------------------
	pickfirstConn, err := grpc.Dial(
		fmt.Sprintf("%s:///%s",
			// exampleScheme = "example"
			// exampleServiceName = "lb.example.grpc.io"
			exampleScheme, exampleServiceName),
		// "pick_first"가 기본값이다.
		grpc.WithBalancerName("pick_first"),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer pickfirstConn.Close()

	log.Println("==== Calling helloworld.Greeter/SayHello " +
		"with pick_first ====")
	makeRPCs(pickfirstConn, 10)

	// round_robin 정책으로 다른 ClientConn을 만든다.
	roundrobinConn, err := grpc.Dial(
		fmt.Sprintf("%s:///%s", exampleScheme, exampleServiceName),
		// "example:///lb.example.grpc.io"
		grpc.WithBalancerName("round_robin"),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer roundrobinConn.Close()

	log.Println("==== Calling helloworld.Greeter/SayHello " +
		"with round_robin ====")
	makeRPCs(roundrobinConn, 10)
	//---------------------------------------------------------
}
