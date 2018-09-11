func main() {
	Address, err := net.ResolveUDPAddr("udp", ":10001")
	CheckError(err)
	Connection, err := net.ListenUDP("udp", Address)
	CheckError(err)
	defer Connection.Close()
	for {
		//Loop infinito para receber mensagem e escrever todo
		//conteúdo (processo que enviou, seu relógio e texto)
		//na tela
		//FALTA FAZER
	}
}