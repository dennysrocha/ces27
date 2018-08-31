package main

import (
	"fmt"
)

func Sqrt(x float64) float64 {
	y := float64(1)
	z := float64(1)
	for {
		y = z
		// a notação z -= (z*z - x)/(2*z) impede que números de ordem maior que
		// aproximadamente 154 sejam utilizados pois causa overflow quando a operação
		// z*z é realizada
		z -= (z/2) - x/(2*z)
		if z-y<1.0e-10 && y-z<1.0e-10 {
			break
		}
	}
	return z
}

func main() {
	fmt.Println(Sqrt(1.7976931348623157e+308))
}