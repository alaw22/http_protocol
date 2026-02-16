package main


import (
	"fmt"
	"os"
	"net"
	"log"
	"bufio"
)

const address = "localhost:42069"

func main(){
	

	// Resolve UDP address
	udpaddr, err := net.ResolveUDPAddr("udp",address)
	if err != nil{
		log.Fatalf("Unable to resolve host '%s': %s\n",address,err.Error())
	}

	conn, err := net.DialUDP("udp", nil, udpaddr)
	if err != nil{
		log.Fatalf("Unable to establish connection to %s: %s\n",address,err.Error())
	}

	defer conn.Close()

	// Create reader
	reader := bufio.NewReader(os.Stdin)

	for {
		// Start of prompt
		fmt.Printf("> ")

		message, err := reader.ReadString('\n')
		if err != nil{
			log.Printf("Error in reader: %s\n",err.Error())
			os.Exit(1)
		}

		// Write string to UDP connection
		_, err = conn.Write([]byte(message))
		if err != nil{
			log.Printf("Error in conn.Write(): %s\n",err.Error())
			os.Exit(1)
		}

		fmt.Printf("Message sent: %s", message)
	}

}
