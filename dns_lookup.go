package main

import (
	"net"
)

func Lookup(qname string, qtype QueryType) (*DnsPacket, error) {
	server := "8.8.8.8:53"

	serverAddr, err := net.ResolveUDPAddr("udp", server)
	if err != nil {
		return nil, err
	}

	localAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:43210")
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", localAddr, serverAddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()


	question := &DnsQuestion{
		Name:  qname,
		QType: qtype,
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
		return nil, err
	}

	_, err = conn.Write(reqBuffer.buf[:])
	if err != nil {
		return nil, err
	}

	resBuffer := NewBytePacketBuffer()
	_, _, err = conn.ReadFromUDP(resBuffer.buf[:])
	if err != nil {
		return nil, err
	}

	resPacket, err := ReadDnsPacket(resBuffer)
	if err != nil {
		return nil, err
	}

	return resPacket, nil
}