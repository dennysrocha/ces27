package main

import (
	"fmt"
	"net"
	"os"
	"time"
	"bufio"
	"strconv"
)

//Variáveis globais interessantes para o processo
var err string
var myPort string //porta do meu servidor
var myProcess int //
var nServers int //qtde de outros processo
var CliConn []*net.UDPConn //vetor com conexões para os servidores dos outros processos
var ServConn *net.UDPConn //conexão do meu servidor (onde recebo mensagens dos outros processos)
var logicalClock int

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
	clock, _ := strconv.Atoi(string(buf[0:n]))
	if clock>logicalClock { //comparo e pego o maior para a variavel logicalClock
		logicalClock = 1+clock
	} else {
		logicalClock++
	}
	fmt.Println("Recebi um logicalClock igual a ", clock, " do ", addr)
	fmt.Println("Agora meu logicalClock eh ", logicalClock)
	fmt.Println("----------------------------------------------------")
	PrintError(err)
}

func doClientJob(x int) {
	fmt.Printf("Enviei para o processo %v ", x)
	logicalClock++
	fmt.Printf("o meu logicalClock = %v\n", logicalClock)
	fmt.Println("----------------------------------------------------")
	buf := []byte(strconv.Itoa(logicalClock))
	if x!=myProcess {
		if x>myProcess { //essa correcao eh feita pra "pular" o indice que seria para este servidor
			x--
		}
		_, err := CliConn[x-1].Write(buf)
		PrintError(err)
	}
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
	fmt.Println("Eu sou o Processo ", myProcess, "!")
	fmt.Println("----------------------------------------------------")
	myPort = os.Args[myProcess+1]
	nServers = len(os.Args) - 3
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
					go doClientJob(x)
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