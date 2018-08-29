package main

import (
	"fmt"
	"net"
	"os"
	"time"
	"bufio"
)

//Variáveis globais interessantes para o processo
var err string
var myPort string //porta do meu servidor
var nServers int //qtde de outros processo
var CliConn []*net.UDPConn //vetor com conexões para os servidores dos outros processos
var ServConn *net.UDPConn //conexão do meu servidor (onde recebo mensagens dos outros processos)

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

func doServerJob() {
	//Ler (uma vez somente) da conexão UDP a mensagem
	buf := make([]byte, 1024)
	n,addr,err := ServConn.ReadFromUDP(buf)

	//Escreve na tela a msg recebida
	fmt.Println("Received ",string(buf[0:n]), " from ",addr)
	PrintError(err)
}

func readInput(ch chan string) {
	// Non-blocking async routine to listen for terminal input
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _, _ := reader.ReadLine()
		ch <- string(text)
	}
}

func initConnections() {
	myPort = os.Args[1]
	nServers = len(os.Args) - 2
	/*Esse 2 tira o nome (no caso Process) e tira a primeira porta (que é a minha). As demais portas são dos outros processos*/

	//Outros códigos para deixar ok a conexão do meu servidor
	ServAddr,err := net.ResolveUDPAddr("udp","127.0.0.1"+myPort)
	CheckError(err)
	ServConn, err = net.ListenUDP("udp", ServAddr)
	CheckError(err)
	CliConn = make([]*net.UDPConn, nServers)

	//Outros códigos para deixar ok as conexões com os servidores dos outros processos
	for i:=0; i<nServers; i++ {
		ServAddr,err = net.ResolveUDPAddr("udp","127.0.0.1"+os.Args[2+i])
		LocalAddr, err := net.ResolveUDPAddr("udp","127.0.0.1:0")
		CheckError(err)
		CliConn[i], err = net.DialUDP("udp", LocalAddr, ServAddr)
		CheckError(err)	
	}
}

func main() {
	initConnections()
	//O fechamento de conexões devem ficar aqui, assim só fecha conexão quando a main morrer
	defer ServConn.Close()
	for i := 0; i < nServers; i++ {
		defer CliConn[i].Close()
	}
	//Todo Process fará a mesma coisa: ouvir msg e mandar infinitos i’s para os outros processos
	ch := make(chan string)
	for {
		//Server
		go doServerJob()
		// When there is a request (from stdin). Do it!
		go readInput(ch)
		select {
			case x, valid := <-ch:
				if valid {
					fmt.Printf("Recebi do teclado: %s \n", x)
					buf := []byte(x)
					for i:=0; i<nServers; i++ {
						_,err := CliConn[i].Write(buf)
						PrintError(err)
					}
				} else {
					fmt.Println("Channel closed!")
				}
			default:
				// Do nothing in the non-blocking approach.
				time.Sleep(time.Second * 1)
		}
		// Wait a while
		time.Sleep(time.Second * 1)
	}
}