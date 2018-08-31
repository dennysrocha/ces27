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
var CliConn []*net.UDPConn //vetor com conexões para os servidores dos outros processos
var ServConn *net.UDPConn //conexão do meu servidor (onde recebo mensagens dos outros processos)
var Data DataClocks

//Estruturas para o processo
type DataClocks struct { //crio a estrutura para o vector timestamp (aqui chamado de VM)
	Id int
	LogicalClock []int
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
	n,_,err := ServConn.ReadFromUDP(buf)
	CheckError(err)

	//Criar uma struct para futuramente armazenar os dados recebidos
	var DataReceived DataClocks
	DataReceived.Id = 0
	for i:=0; i<nServers+1; i++ { //aloco espaco para cada dimensao do meu vetor
		DataReceived.LogicalClock = append(DataReceived.LogicalClock, 1)
	}

	err = json.Unmarshal(buf[:n], &DataReceived) //interpreto por meio do json e passo pra estrutura de dados
	PrintError(err)

	fmt.Println("Recebi msg. O VT recebido eh: ", DataReceived.LogicalClock)
	fmt.Println("Meu VT era: ", Data.LogicalClock)
	//Escreve na tela a msg recebida
	for i:=0; i<nServers+1; i++ { //comparo e sempre pego o maior
		if DataReceived.LogicalClock[i] > Data.LogicalClock[i] {
			Data.LogicalClock[i] = DataReceived.LogicalClock[i]
		}
	}
	Data.LogicalClock[myProcess-1]++
	fmt.Println("Agora, meu VT eh ", Data.LogicalClock)
	fmt.Println("-----------------------------")
	PrintError(err)
}

func doClientJob(x int) {
	Data.LogicalClock[myProcess-1]++
	fmt.Println("Meu VT eh: ", Data.LogicalClock)
	fmt.Printf("Enviando msg pro proc %v\n", x)
	fmt.Println("-----------------------------")
	if x!=myProcess {
		jsonRequest, err := json.Marshal(Data) //reescrevo os dados por meio do json
		CheckError(err)
		if x>myProcess {
			x--
		}
		_, err = CliConn[x-1].Write(jsonRequest) //envio os dados reescritos pelo canal
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


/////////////////////////////////////////////ABAIXO//////////////////////////////////////////////////////////////

func initConnections() {
	myProcess, _ = strconv.Atoi(os.Args[1])
	fmt.Println("Este eh o processo ", myProcess)
	myPort = os.Args[myProcess+1]
	nServers = len(os.Args) - 3
	/*Esse 3 tira o nome (no caso Process), o numero do (meu) processo e tira a porta que é minha. As demais portas são dos outros processos*/

	Data.Id = myProcess

	//Outros códigos para deixar ok a conexão do meu servidor
	ServAddr,err := net.ResolveUDPAddr("udp","127.0.0.1"+myPort)
	CheckError(err)
	ServConn, err = net.ListenUDP("udp", ServAddr)
	CheckError(err)
	CliConn = make([]*net.UDPConn, nServers)

	//Outros códigos para deixar ok as conexões com os servidores dos outros processos
	j:=0 //esse j eh apenas para "pular" o i correspondente ao meu servidor
	for i:=0; i<nServers+1; i++ {
		Data.LogicalClock = append(Data.LogicalClock, 1)
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