package rpc

import (
	"github.com/BeroKiTeer/MyGoMall/common/kitex_gen/cart/cartservice"
	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"
	"order/conf"
)

var (
	CartClient cartservice.Client
)

func initCartClient() {
	var opts []client.Option
	r, err := consul.NewConsulResolver(conf.GetConf().Registry.RegistryAddress[0])
	if err != nil {
		panic(err)
	}
	opts = append(opts, client.WithResolver(r))
	CartClient, err = cartservice.NewClient("cart", opts...)
	if err != nil {
		panic(err)
	}
}
