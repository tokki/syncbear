package ray

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	proxymancmd "v2ray.com/core/app/proxyman/command"
	statscmd "v2ray.com/core/app/stats/command"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/vmess"
)

type ServiceClient struct {
	APIAddress  string
	APIPort     int
	statClient  statscmd.StatsServiceClient
	proxyClient proxymancmd.HandlerServiceClient
}

func New(addr string, port int) *ServiceClient {
	Conn, err := grpc.Dial(fmt.Sprintf("%s:%d", addr, port), grpc.WithInsecure())
	if err != nil {
		log.Println(err)
		return nil
	}

	client := ServiceClient{APIAddress: addr, APIPort: port,
		statClient:  statscmd.NewStatsServiceClient(Conn),
		proxyClient: proxymancmd.NewHandlerServiceClient(Conn),
	}
	return &client
}

func (h *ServiceClient) Traffic(pattern string, reset bool) map[string]int64 {
	resp, err := h.statClient.QueryStats(context.Background(), &statscmd.QueryStatsRequest{
		Pattern: pattern,
		Reset_:  reset,
	})

	result := make(map[string]int64)

	if err != nil {
		fmt.Println(err)
	} else {
		for _, res := range resp.Stat {
			result[res.Name] =res.Value
		}

	}
	return result
}


func (h *ServiceClient) AddUser(inboundTag string, email string, level uint32, uuid string, alterID uint32) {
	resp, err := h.proxyClient.AlterInbound(context.Background(), &proxymancmd.AlterInboundRequest{
		Tag: inboundTag,
		Operation: serial.ToTypedMessage(&proxymancmd.AddUserOperation{
			User: &protocol.User{
				Level: level,
				Email: email,
				Account: serial.ToTypedMessage(&vmess.Account{
					Id:               uuid,
					AlterId:          alterID,
					SecuritySettings: &protocol.SecurityConfig{Type: protocol.SecurityType_AUTO},
				}),
			},
		}),
	})

	if err != nil {
		fmt.Println( err)
	} else {
		fmt.Println( resp)
	}
}

func (h *ServiceClient) RemoveUser(inboundTag string, email string) {
	resp, err := h.proxyClient.AlterInbound(context.Background(), &proxymancmd.AlterInboundRequest{
		Tag: inboundTag,
		Operation: serial.ToTypedMessage(&proxymancmd.RemoveUserOperation{
			Email: email,
		}),
	})

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp)
	}
}
