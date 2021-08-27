package socket

import (
	"fmt"
	"net"
)

type IP []byte

func Ip() {
	fmt.Println(net.IPv4(255,255,0,0))
	fmt.Println(net.IPv4(0,0,0,0))

	newIp := net.ParseIP("0:0:1:1:0:0:1:1")
	fmt.Println(newIp.String())

	address := net.IPv4(60,60,0,0)

	mask := address.DefaultMask()
	network := address.Mask(mask)
	_, n := mask.Size()
	fmt.Printf("address : %v, mask : %v, network: %v, size : %v\n", address, mask, network, n)

	resolve := "www.promedic.com.ng"
	resolve2 := "www.salespro.ng"
	resolvedAddr,err := net.ResolveIPAddr("ip", resolve)
	resolvedAddr2,err := net.LookupIP(resolve2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v : %v\n", resolve, resolvedAddr)
	fmt.Printf("%v : %v\n", resolve2, resolvedAddr2)

}