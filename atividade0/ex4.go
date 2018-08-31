package main

import "fmt"

type IPAddr [4]byte

// TODO: Add a "String() string" method to IPAddr.

func (m IPAddr) String() string {
	ip := ""
	for i:=0; i<len(m)-1; i++ {
		ip += fmt.Sprintf("%v.", m[i])
	}
	return ip+fmt.Sprintf("%v", m[len(m)-1])
}

func main() {
	hosts := map[string]IPAddr{
		"loopback":  {127, 0, 0, 1},
		"googleDNS": {8, 8, 8, 8},
	}
	for name, ip := range hosts {
		fmt.Printf("%v: %v\n", name, ip)
	}
}