package main

import (
	"flag"
	"fmt"
	"github.com/jtzjtz/kit/balance"
	pb "github.com/jtzjtz/kit/balance/example/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const svcName = "project"

var addr = "127.0.0.1"
var port = 50054

func main() {
	flag.StringVar(&addr, "addr", addr, "addr to lis")
	flag.IntVar(&port, "port", port, "port to lis")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%v", addr, port))
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}
	defer lis.Close()

	s := grpc.NewServer()
	defer s.GracefulStop()

	pb.RegisterGreeterServer(s, &server{})

	//注册服务 开始

	err = balance.RegisterService("testservice", "test-app", addr, port, 10, nil, balance.NacosConfig{NacosIp: "localhost", NacosPort: 8848, NamespaceId: "96023333-63a2-42aa-a456-b880587931e2"})
	//err = balance.RegisterService("testservice", "test-app", addr, port, 10, nil, balance.NacosConfig{NacosIp: "localhost", NacosPort: 8848, NamespaceId: "fffaacb2-a7f7-45d5-83aa-dec49122a098"})
	if err != nil {
		panic("服务注册失败，禁止启动")
	}
	//注册服务 结束

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <-ch
		//取消注册 开始
		balance.UnRegisterService("testservice", "test-app", addr, port)
		//取消注册 结束

		if i, ok := s.(syscall.Signal); ok {
			os.Exit(int(i))
		} else {
			os.Exit(0)
		}

	}()

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: fmt.Sprintf("Hello "+in.Name+"! From %s:%v", addr, port)}, nil
}
