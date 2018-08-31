package main

import (
	"fmt"
)

type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string {
	return fmt.Sprintf("cannot Sqrt negative number: %.0f", e) // Usei o fmt.Sprint pra concatenar uma string com um número, ao inves de ficar convertendo o número pra string
}

func Sqrt(x float64) (float64, error) {
	if x<0 {
		err := ErrNegativeSqrt(x)
		return x, err
	}
	
	y := float64(1)
	z := float64(1)
	for {
		y = z
		z -= (z/2) - x/(2*z)
		if z-y<1.0e-10 && y-z<1.0e-10 { //definindo um intervalo de erro confiavel
			break
		}
	}
	return z, nil
}

func main() {
	fmt.Println(Sqrt(2))
	fmt.Println(Sqrt(-2))
}
