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
var myProcess int //numero do meu processo
var nServers int //qtde de outros processo
var myID int // meu ID que eh digitado no prompt
var leaderID int // guarda o ID do lider
var isParticipant bool // booleana para participacao
var CliConn []*net.UDPConn //vetor com conexões para os servidores dos outros processos
var ServConn *net.UDPConn //conexão do meu servidor (onde recebo mensagens dos outros processos)

//Estruturas para o processo
type Message struct { //crio a estrutura para o vector timestamp (aqui chamado de VM)
	Id int
	Type string
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

func doServerJob() {
	//Ler (uma vez somente) da conexão UDP a mensagem
	buf := make([]byte, 1024)
	n,_,err := ServConn.ReadFromUDP(buf)
	CheckError(err)

	err = json.Unmarshal(buf[:n], &MessageReceived) //interpreto por meio do json e passo pra estrutura de dados
	PrintError(err)

	fmt.Print("Recebi: ", MessageReceived.Type, MessageReceived.Id, "\n") //Escreve na tela a msg recebida
	
	if MessageReceived.Type=="S" {
		go stage1(MessageReceived.Id)
	} else if MessageReceived.Type=="F" {
		go stage2(MessageReceived.Id)
	} else {
		fmt.Println("Entrada invalida. Fim do programaa")
		os.Exit(0)
	}
	fmt.Println("-----------------------------")
	PrintError(err)
}

func stage1(receivedID int) {
	time.Sleep(time.Second)
	if receivedID<myID {
		receivedID = myID
	} else if receivedID==myID && isParticipant{
		leaderID = myID
		stage2(receivedID)
	}
	isParticipant = true
	doClientJob((myID+1)%(nServers+1))
}

func stage2(receivedID int) {
	time.Sleep(time.Second)
	isParticipant = false // nem precisava, a booleana ja vem como false
	MessageReceived.Type = "F"
	doClientJob((myID+1)%(nServers+1))
}

func doClientJob(x int) {
	if x==0 {
		x = nServers + 1
	}
	fmt.Print("Enviando ", MessageReceived.Type, MessageReceived.Id)
	fmt.Println("\n-----------------------------")
	jsonRequest, err := json.Marshal(MessageReceived) //reescrevo os dados por meio do json
	CheckError(err)
	if x>myProcess {
		x--
	}
	_, err = CliConn[x-1].Write(jsonRequest) //envio os dados reescritos pelo canal
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


/////////////////////////////////////////////ABAIXO//////////////////////////////////////////////////////////////

func initConnections() {
	myProcess, _ = strconv.Atoi(os.Args[1])
	fmt.Println("Este eh o processo ", myProcess)
	myPort = os.Args[myProcess+1]
	nServers = len(os.Args) - 3
	/*Esse 3 tira o nome (no caso Process), o numero do (meu) processo e tira a porta que é minha. As demais portas são dos outros processos*/

	fmt.Print("Digite meu ID: ")
	reader := bufio.NewReader(os.Stdin)
	text, _, _ := reader.ReadLine()
	myID, _ = strconv.Atoi(string(text))


	//Outros códigos para deixar ok a conexão do meu servidor
	ServAddr,err := net.ResolveUDPAddr("udp","127.0.0.1"+myPort)
	CheckError(err)
	ServConn, err = net.ListenUDP("udp", ServAddr)
	CheckError(err)
	CliConn = make([]*net.UDPConn, nServers)

	//Outros códigos para deixar ok as conexões com os servidores dos outros processos
	j:=0 //esse j eh apenas para "pular" o i correspondente ao meu servidor
	for i:=0; i<nServers+1; i++ {
		if i!=myProcess-1 {
			ServAddr,err = net.ResolveUDPAddr("udp","127.0.0.1"+os.Args[i+2])
			CheckError(err)
			LocalAddr, err := net.ResolveUDPAddr("udp","127.0.0.1:0")
			CheckError(err)
			CliConn[j], err = net.DialUDP("udp", LocalAddr, ServAddr)
			CheckError(err)
			j++
		}
	}
	fmt.Println("Conexoes inicializadas")
	fmt.Println("-----------------------------")
}

/////////////////////////////////////////////ACIMA///////////////////////////////////////////////////////////////

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
		go readInput(ch)
		//Server
		go doServerJob()
		// When there is a request (from stdin). Do it!
		select {
			case x, valid := <-ch:
				if valid && x=="S" {
					fmt.Println("Enviando start")
					MessageReceived.Type = "S"
					MessageReceived.Id = myID
					go stage1(myID)
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