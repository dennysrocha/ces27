package main

import (
	//"fmt"
	"golang.org/x/tour/wc"
	"strings"
)

func WordCount(s string) map[string]int {
	n := strings.Fields(s)			// inicializo um array de strings
	m := make(map[string]int)		// inicializo um mapa para cada string do array
	for i:=0; i<len(n); i++ {		// inicializo o mapeamento de todas as strings para 0
		m[n[i]]=0
	}
	for i:=0; i<len(n); i++ {		// incremento o valor do mapeamento a cada vez que a string aparece
		m[n[i]]+=1
	}
	return m
}

func main() {
	wc.Test(WordCount)
}