package main

import (
	"golang.org/x/tour/tree"
	"fmt"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	var walk func(t *tree.Tree)
	walk = func(t *tree.Tree) { //caminhar pela arvore e ir jogando as folhas no canal
		if(t.Left!=nil) {
			walk(t.Left)
		}
		ch<-t.Value
		if(t.Right!=nil) {
			walk(t.Right)
		}
	}
	walk(t)
	close(ch)
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1, ch2 := make(chan int), make(chan int)
	go Walk(t1, ch1)
	go Walk(t2, ch2)
	for n:=range ch1 { //compara item por item dos canais, ordenadamente
		if n!=<-ch2 {
			return false
		}
	}
	return true
}


func main() {
	ch0 := make(chan int)
	go Walk(tree.New(1), ch0)
	for i:=0; i<10; i++ { //printar os itens obtidos no canal pela arvore
		n:=<-ch0
		fmt.Print(n," ")
	}
	fmt.Println()
	fmt.Println(Same(tree.New(1), tree.New(1))) //comparacao das arvores; caso OK
	fmt.Println(Same(tree.New(1), tree.New(2))) //comparacao das arvores; caso FALSO
}
