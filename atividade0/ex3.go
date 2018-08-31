package main

import "fmt"

// fibonacci is a function that returns
// a function that returns an int.
func fibonacci() func() int {
	n:=0
	m:=1
	return func() int {
		aux:=n
		n+=m
		m=aux
		return aux
	}
}

func main() {
	f := fibonacci()
	for i := 0; i < 47; i++ { // é importante notar que o overflow ocorre quanto i=47
		fmt.Println(f())
	}
}