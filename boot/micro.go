package boot

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/995933447/gonetutil"
	"github.com/gzjjyz/logger"
	"github.com/gzjjyz/micro"
	"github.com/gzjjyz/micro/discovery"
	"github.com/gzjjyz/micro/factory"
	"github.com/gzjjyz/srvlib/handler/grpchandler"
	"github.com/gzjjyz/srvlib/pb3/health"
	"github.com/gzjjyz/srvlib/utils"
	"google.golang.org/grpc"
	"net"
	"net/http"
	_ "net/http/pprof"
)

func InitMicroSuiteWithGrpc(ctx context.Context, metaFilePath string) error {
	if err := micro.InitSuitWithGrpc(ctx, metaFilePath); err != nil {
		return err
	}
	return nil
}

func ServeGrpc(ctx context.Context, srvName string, ipVar string, port, pprofPort int, registerCustomServiceServer func(*grpc.Server)) error {
	go func() {
		err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", pprofPort), nil)
		if err != nil {
			logger.Errorf(err.Error())
		}
	}()

	ip, err := utils.EvalVarToParseIp(ipVar)
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}

	nodeExtra, err := json.Marshal(&health.NodeHealthDesc{})
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	node := discovery.NewNode(ip, port)
	node.Extra = string(nodeExtra)

	grpcServer := grpc.NewServer()
	registerCustomServiceServer(grpcServer)
	health.RegisterHealthServer(grpcServer, &grpchandler.Health{})

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return err
	}

	discover, err := factory.GetOrMakeDiscovery()
	if err != nil {
		return err
	}

	defer func() {
		err = discover.Unregister(ctx, srvName, node, true)
		if err != nil {
			logger.Errorf(err.Error())
		}
	}()
	err = discover.Register(ctx, srvName, node)
	if err != nil {
		return err
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		return err
	}

	return nil
}

type ServeGrpcReq struct {
	RegDiscoverKeyPrefix            string
	SrvName                         string
	IpVar                           string
	Port                            int
	PProfIpVar                      string
	PProfPort                       int
	RegisterCustomServiceServerFunc func(*grpc.Server) error
	BeforeRegDiscover               func(discovery.Discovery, *discovery.Node) error
	AfterRegDiscover                func(discovery.Discovery, *discovery.Node) error
	EnabledHealth                   bool
	SrvOpts                         []grpc.ServerOption
}

func ServeGrpcV2(ctx context.Context, req *ServeGrpcReq) error {
	if req.PProfIpVar != "" && req.PProfPort > 0 {
		go func() {
			ip, err := gonetutil.EvalVarToParseIp(req.PProfIpVar)
			if err != nil {
				logger.Errorf("%v", err)
				return
			}

			err = http.ListenAndServe(fmt.Sprintf("%s:%d", ip, req.PProfPort), nil)
			if err != nil {
				logger.Errorf("%v", err)
			}
		}()
	}

	ip, err := gonetutil.EvalVarToParseIp(req.IpVar)
	if err != nil {
		return err
	}

	node := discovery.NewNode(ip, req.Port)
	grpcServer := grpc.NewServer(req.SrvOpts...)
	if req.RegisterCustomServiceServerFunc != nil {
		if err = req.RegisterCustomServiceServerFunc(grpcServer); err != nil {
			return err
		}
	}

	if req.EnabledHealth {
		health.RegisterHealthServer(grpcServer, &grpchandler.Health{})
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, req.Port))
	if err != nil {
		return err
	}

	discover, err := factory.GetOrMakeDiscovery()
	if err != nil {
		return err
	}

	if req.BeforeRegDiscover != nil {
		if err = req.BeforeRegDiscover(discover, node); err != nil {
			return err
		}
	}

	err = discover.Register(ctx, req.SrvName, node)
	if err != nil {
		return err
	}

	if req.AfterRegDiscover != nil {
		if err = req.AfterRegDiscover(discover, node); err != nil {
			return err
		}
	}

	defer func() {
		err = discover.Unregister(ctx, req.SrvName, node, true)
		if err != nil {
			logger.Errorf("%v", err)
		}
	}()

	err = grpcServer.Serve(listener)
	if err != nil {
		return err
	}

	return nil
}
