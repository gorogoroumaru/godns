package main

import (
	"fmt"
	"net"

	"github.com/k0kubun/pp/v3"
)

func handleQuery(socket *net.UDPConn) error {
	reqBuffer := NewBytePacketBuffer()

	_, src, err := socket.ReadFromUDP(reqBuffer.buf[:])
	if err != nil {
		return err
	}

	request, err := ReadDnsPacket(reqBuffer)
	if err != nil {
		return err
	}

	packet := &DnsPacket{
		Header: &DnsHeader{
			ID:                request.Header.ID,
			RecursionDesired:  true,
			RecursionAvailable: true,
			Response:          true,
		},
	}

	if len(request.Questions) > 0 {
		question := request.Questions[0]
		fmt.Printf("Received query: %+v\n", question)

		result, err := RecursiveLookup(question.Name, question.QType)

		if err != nil {
			packet.Header.ResCode = SERVFAIL
		} else {
			packet.Questions = append(packet.Questions, question)
			packet.Header.ResCode = result.Header.ResCode

			for _, rec := range result.Answers {
				fmt.Printf("Answer: ")
				pp.Print(rec)
				packet.Answers = append(packet.Answers, rec)
			}
			for _, rec := range result.Authorities {
				fmt.Printf("Authority: ")
				pp.Print(rec)
				packet.Authorities = append(packet.Authorities, rec)
			}
			for _, rec := range result.Resources {
				fmt.Printf("Resource: ")
				pp.Print(rec)
				packet.Resources = append(packet.Resources, rec)
			}
		}
	} else {
		packet.Header.ResCode = FORMERR
	}

	resBuffer := NewBytePacketBuffer()
	err = packet.Write(resBuffer)
	if err != nil {
		return err
	}

	_, err = socket.WriteToUDP(resBuffer.buf[:resBuffer.Pos()], src)
	return err
}

func main() {
	socket, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 2053})
	if err != nil {
		fmt.Printf("Failed to bind UDP socket: %v\n", err)
		return
	}
	defer socket.Close()

	fmt.Println("Listening on UDP port 2053...")

	for {
		err := handleQuery(socket)
		if err != nil {
			fmt.Println("An error occurred", err)
		}
	}
}
