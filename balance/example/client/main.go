package main

import (
	"fmt"

	"github.com/jtzjtz/kit/balance"
	pb "github.com/jtzjtz/kit/balance/example/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"time"
)

func main() {

	//客户端发现服务开始
	r := balance.NewResolver(balance.NacosConfig{NacosIp: "localhost", NacosPort: 8848, NamespaceId: "962222fe-63a2-42aa-a456-b880587931e2", GroupName: "test-app"})
	resolver.Register(r)
	serviceName := "testservice"
	//pick_first/round_robin
	conn, err := grpc.Dial(r.Scheme()+":///"+serviceName, grpc.WithDefaultServiceConfig(`{"LoadBalancingPolicy": "round_robin"}`), grpc.WithInsecure())
	//conn, err := grpc.Dial(r.Scheme()+"://author/"+serviceName, grpc.WithBalancerName("round_robin"), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	//客户端发现服务结束

	client := pb.NewGreeterClient(conn)

	for {
		resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "hello"}, grpc.FailFast(true))
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(resp)
		}

		<-time.After(time.Second)
	}
}
