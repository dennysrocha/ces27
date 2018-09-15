package main

import (
	"fmt"
	"net"
	"os"
	"encoding/json"
	"time"
)

var ServConn *net.UDPConn //conexão do meu servidor (onde recebo mensagens dos outros processos)

type Message struct { //crio a estrutura para o vector timestamp (aqui chamado de VM)
	Id int
	Text string
	LogicalClock int
}

var MessageReceived Message

func CheckError(err error) {
	if err != nil {
		fmt.Println("Erro: ", err)
		os.Exit(0)
	}
}

func PrintError(err error) {
	if err != nil {
		fmt.Println("Erro: ", err)
	}
}

func main() {
	fmt.Println("Regiao crítica estabelecida!")
	fmt.Println("-----------------------------")
	Address, err := net.ResolveUDPAddr("udp", ":10001")
	CheckError(err)
	ServConn, err = net.ListenUDP("udp", Address)
	CheckError(err)
	defer ServConn.Close()
	for {
		//Loop infinito para receber mensagem e escrever todo
		buf := make([]byte, 1024)
		n,_,err := ServConn.ReadFromUDP(buf)
		CheckError(err)
		//conteúdo (processo que enviou, seu relógio e texto)
		err = json.Unmarshal(buf[:n], &MessageReceived) //interpreto por meio do json e passo pra estrutura de dados
		PrintError(err)
		
		
		//na tela
		fmt.Print("Processo de ID ", MessageReceived.Id, " entrou na CS\n")
		fmt.Print("Relógio Lógico: ", MessageReceived.LogicalClock, "\n-\n")
		fmt.Print("3...\n")
		time.Sleep(1*time.Second)
		fmt.Print("2...\n")
		time.Sleep(1*time.Second)
		fmt.Print("1...\n")
		time.Sleep(1*time.Second)
		fmt.Print("Go!\n")
		time.Sleep(1*time.Second)
		fmt.Print("Same song, different chorus: \n")
		time.Sleep(3*time.Second)
		fmt.Print("It's stupid, contagious\n")
		time.Sleep(2*time.Second)
		fmt.Print("To be broke and famous\n")
		time.Sleep(2*time.Second)
		fmt.Print("Can someone please save us from punk rock 101!\n")
		time.Sleep(4*time.Second)
		fmt.Println("-----------------------------")
		fmt.Print(MessageReceived.Text)
		time.Sleep(1*time.Second)
		fmt.Println("-----------------------------")
	}
}