package socket_manager

import (
	"net"

	"crypto/tls"
	"mylib/src/public"
	"time"
	"context"
	"github.com/quic-go/quic-go"
	//"github.com/quic-go/quic-go/http3"
)

type userCallbackBigPayload	func(params string, big_payload	string)(interface{}, bool)
type userCallback			func(params string)(interface{}, bool)

type routeType	int8
const(
	TYPE_TCP	routeType = 1
	TYPE_UDP	routeType = 2
	TYPE_QUIC	routeType = 3
)

type Route_Socket_Manager struct{
	routes		[]*Route_Socket_Unit
}

type Route_Socket_Unit struct{
	route_type				routeType
	callback				userCallbackBigPayload
	callback_big_payload	userCallback
}


func New()*Route_Socket_Manager{
	return &Route_Socket_Manager{}
}

func new_route(callback interface{}, route_type routeType)*Route_Socket_Unit{}{
	ret				:= &Route_Socket_Unit{}
	ret.route_type	= route_type
	switch callback.(type){
		case userCallbackBigPayload:
			ret.callback_big_payload	= callback.(userCallbackBigPayload)
		case userCallback:
			ret.callback				= callback.(userCallback)
		default:
			public.DBG_ERR("callback[", callback, "] no support.")
	}

	return ret
}

func (rsm *Route_Socket_Manager) Route_TCP(callback interface{})*Route_Socket_Unit{
	ret := new_route(callback, TYPE_TCP)
	
	rsm.routes = append(rsm.routes, ret)

	return ret
}

func (rsm *Route_Socket_Manager) Route_UDP(callback interface{})*Route_Socket_Unit{
	ret := new_route(callback, TYPE_UDP)
	
	rsm.routes = append(rsm.routes, ret)
	
	return ret
}

func (rsm *Route_Socket_Manager) Route_QUIC(callback interface{})*Route_Socket_Unit{
	ret := new_route(callback, TYPE_QUIC)
	
	rsm.routes = append(rsm.routes, ret)
	
	return ret
}



