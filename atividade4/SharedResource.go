package main

import (
	"fmt"
	"net"
	"os"
	"encoding/json"
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
	Address, err := net.ResolveUDPAddr("udp", ":10001")
	CheckError(err)
	Connection, err := net.ListenUDP("udp", Address)
	CheckError(err)
	defer Connection.Close()
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
		fmt.Print("Relógio Lógico: ", MessageReceived.LogicalClock, "\n")
		fmt.Print("Texto: ", MessageReceived.Text, "\n")
	}
}