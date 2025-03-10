package serversuite

import (
	"github.com/BeroKiTeer/MyGoMall/common/mtl"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
	prometheus "github.com/kitex-contrib/monitor-prometheus"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	consul "github.com/kitex-contrib/registry-consul"
)

type CommonServerSuite struct {
	CurrentServiceName string
	RegistryAddr       string
}

func (s CommonServerSuite) Options() []server.Option {
	opts := []server.Option{
		// 使用 HTTP2 协议
		server.WithMetaHandler(transmeta.ServerHTTP2Handler),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: s.CurrentServiceName,
		}),
		// 增加 Prometheus 中间件
		server.WithTracer(prometheus.NewServerTracer(
			"", // 默认监控地址
			"", // 默认监控路径
			prometheus.WithDisableServer(true),
			prometheus.WithRegistry(mtl.Registry),
		)),
		server.WithSuite(tracing.NewServerSuite()),
	}
	r, err := consul.NewConsulRegister(s.RegistryAddr) // 从配置获取注册中心地址
	if err != nil {
		klog.Fatal(err)
	}
	opts = append(opts, server.WithRegistry(r)) // 将服务注册组件加入配置选项

	return opts
}
