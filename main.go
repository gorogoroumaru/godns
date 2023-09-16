// goでdnsサーバー
// https://github.com/EmilHernvall/dnsguide
// qname圧縮に関する仕様の解説
// http://park12.wakwak.com/~eslab/pcmemo/dns/dns5.html#condense

package main

import (
	"fmt"
	"net"
	"os"

	"github.com/k0kubun/pp/v3"
)


func main() {
	qname := "yahoo.com"
	qtype := CNAME

	server := "8.8.8.8:53"

	serverAddr, err := net.ResolveUDPAddr("udp", server)
	if err != nil {
		fmt.Println("Error resolving server address:", err)
		os.Exit(1)
	}

	localAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:43210")
	if err != nil {
		fmt.Println("Error resolving local address:", err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", localAddr, serverAddr)
	if err != nil {
		fmt.Println("Error creating UDP socket:", err)
		os.Exit(1)
	}
	defer conn.Close()


	question := &DnsQuestion{
		Name:  qname,
		QType: QueryType{
			query_type: uint16(qtype),
			val: uint16(qtype),
		},
	}

	packet := DnsPacket{
		Header: &DnsHeader{
			ID:               6666,
			Questions:        1,
			RecursionDesired: true,
		},
		Questions: []*DnsQuestion{question},
	}

	reqBuffer := NewBytePacketBuffer()
	if err := packet.Write(reqBuffer); err != nil {
		fmt.Println("Error writing request packet:", err)
		os.Exit(1)
	}

	_, err = conn.Write(reqBuffer.buf[:])
	if err != nil {
		fmt.Println("Error sending request packet:", err)
		os.Exit(1)
	}

	resBuffer := NewBytePacketBuffer()
	_, _, err = conn.ReadFromUDP(resBuffer.buf[:])
	if err != nil {
		fmt.Println("Error receiving response packet:", err)
		os.Exit(1)
	}

	resPacket, err := ReadDnsPacket(resBuffer)
	if err != nil {
		fmt.Println("Error parsing response packet:", err)
		os.Exit(1)
	}

	pp.Print(resPacket.Header)

	for _, q := range resPacket.Questions {
		pp.Print(q)
	}
	for _, rec := range resPacket.Answers {
		pp.Print(rec)
	}
	for _, rec := range resPacket.Authorities {
		pp.Print(rec)
	}
	for _, rec := range resPacket.Resources {
		pp.Print(rec)
	}
}
