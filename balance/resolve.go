package balance

import (
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/util"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"google.golang.org/grpc/resolver"
	"log"
	"strconv"
)

const schema = "sq"

type NacosResolver struct {
	nacosConf NacosConfig
	cc        resolver.ClientConn
}

// NewResolver initialize an etcd client
func NewResolver(nacosConf NacosConfig) resolver.Builder {
	return &NacosResolver{nacosConf: nacosConf}
}

func (r *NacosResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	var err error

	if cli == nil {
		cli, err = createNacosClient(r.nacosConf)
		if err != nil {
			return nil, err
		}
	}
	r.cc = cc

	go r.watch(target.Endpoint, target.Scheme)
	return r, nil
}

func (r NacosResolver) Scheme() string {
	return r.nacosConf.GroupName
}

func (r NacosResolver) ResolveNow(rn resolver.ResolveNowOptions) {

	log.Println(r.nacosConf.GroupName + "ResolveNow")
}

// Close closes the resolver.
func (r NacosResolver) Close() {
	log.Println("Close")
}

func (r *NacosResolver) watch(serviceName string, group string) {
	var addrList []resolver.Address

	getResp, err := cli.SelectAllInstances(vo.SelectAllInstancesParam{
		ServiceName: serviceName,
		GroupName:   group,
	})
	if err != nil {
		log.Println(err)
	} else {
		for _, instance := range getResp {
			addr := instance.Ip + ":" + strconv.FormatUint(instance.Port, 10)
			if instance.Enable {
				addrList = append(addrList, resolver.Address{Addr: addr, ServerName: instance.ServiceName})
			}
		}

	}
	r.cc.UpdateState(resolver.State{Addresses: addrList})
	subScribeParam := &vo.SubscribeParam{
		ServiceName: serviceName,
		//Clusters:    []string{"cluster-b"},
		GroupName: group,
		SubscribeCallback: func(services []model.SubscribeService, err error) {
			log.Printf("订阅 return services:%s \n\n", util.ToJsonString(services))
			var addrListNew []resolver.Address
			for _, eventService := range services {
				addr := eventService.Ip + ":" + strconv.FormatUint(eventService.Port, 10)
				if eventService.Enable {
					addrListNew = append(addrListNew, resolver.Address{Addr: addr, ServerName: eventService.ServiceName})
				}
			}
			r.cc.UpdateState(resolver.State{Addresses: addrListNew})
			log.Printf("新的address list:%s \n\n", util.ToJsonString(addrListNew))

		},
	}
	cli.Subscribe(subScribeParam)

}
