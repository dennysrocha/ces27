package main

import (
	"fmt"
	"net"
	"os"
	"time"
	"bufio"
	"strconv"
	"encoding/json"
)

//Variáveis globais interessantes para o processo
var err string
var myPort string //porta do meu servidor
var myProcess int //
var nServers int //qtde de outros processo
var CliConn []*net.UDPConn //vetor com conexões para os servidores dos outros processos
var ServConn *net.UDPConn //conexão do meu servidor (onde recebo mensagens dos outros processos)
//var logicalClock int
var data DataClocks

type DataClocks struct {
	id int
	logicalClock []int
}

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

	var dataReceived DataClocks
	err = json.Unmarshal(buf[:n], &dataReceived)
	PrintError(err)

	//Escreve na tela a msg recebida
	if data.logicalClock[myProcess]>dataReceived.logicalClock[dataReceived.id] {
		logicalClock = 1+clock
	} else {
		logicalClock++
	}
	fmt.Println("Recebi um logicalClock igual a ",clock, " do ",addr)
	fmt.Println("Agora meu logicalClock eh ",logicalClock)
	PrintError(err)
}

func readInput(ch chan int) {
	// Non-blocking async routine to listen for terminal input
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _, _ := reader.ReadLine()
		aux, _ := strconv.Atoi(string(text))
		ch <- aux
	}
}

func initConnections() {
	myProcess, _ = strconv.Atoi(os.Args[1])
	fmt.Println("A variavel myProcess recebeu", myProcess)
	myPort = os.Args[myProcess+1]
	nServers = len(os.Args) - 3
	data := DataClocks {
		id: myProcess,
		logicalClock: make([]int, nServers+1),
	}
	/*Esse 3 tira o nome (no caso Process), o numero do (meu) processo e tira a porta que é minha. As demais portas são dos outros processos*/


	//Outros códigos para deixar ok a conexão do meu servidor
	ServAddr,err := net.ResolveUDPAddr("udp","127.0.0.1"+myPort)
	CheckError(err)
	ServConn, err = net.ListenUDP("udp", ServAddr)
	CheckError(err)
	CliConn = make([]*net.UDPConn, nServers)

	//Outros códigos para deixar ok as conexões com os servidores dos outros processos
	j:=0
	for i:=0; i<nServers+1; i++ {
		data.logicalClock[i]=1
		fmt.Println("inicializei o LC do id ", i ," com o valor de ", data.logicalClock[i])
		if i!=myProcess-1 {
			ServAddr,err = net.ResolveUDPAddr("udp","127.0.0.1"+os.Args[i+2])
			LocalAddr, err := net.ResolveUDPAddr("udp","127.0.0.1:0")
			CheckError(err)
			CliConn[j], err = net.DialUDP("udp", LocalAddr, ServAddr)
			CheckError(err)
			j++
		}
	}
}

func main() {
	logicalClock=1
	initConnections()
	//O fechamento de conexões devem ficar aqui, assim só fecha conexão quando a main morrer
	defer ServConn.Close()
	for i := 0; i < nServers; i++ {
		defer CliConn[i].Close()
	}
	//Todo Process fará a mesma coisa: ouvir msg e mandar infinitos i’s para os outros processos
	ch := make(chan int)
	for {
		//Server
		go doServerJob()
		// When there is a request (from stdin). Do it!
		go readInput(ch)
		select {
			case x, valid := <-ch:
				if valid {
					fmt.Printf("Enviei para o processo %v ", x)
					fmt.Printf("o meu struct data\n")
					buf := []byte(strconv.Itoa(logicalClock))
					if x!=myProcess {
						if x>myProcess {
							x--
						}
						_, err := CliConn[x-1].Write(buf)
						PrintError(err)
					} else {
						logicalClock++
						fmt.Println("Estou incrementando meu logicalClock em 1 unidade. Agora ele vale logicalClock =",logicalClock)
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