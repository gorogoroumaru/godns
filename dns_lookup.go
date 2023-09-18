package main

import (
	"errors"
	"fmt"
	"net"

	"github.com/k0kubun/pp/v3"
)

func Lookup(qname string, qtype QueryType, serverAddr *net.UDPAddr) (*DnsPacket, error) {
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

	packet := NewDnsPacket()

	packet.Header = NewDnsHeader()
	packet.Header.ID = 5555
	packet.Header.Questions = 1
	packet.Header.RecursionDesired = true

	packet.Questions = append(packet.Questions, question)

	reqBuffer := NewBytePacketBuffer()
	if err := packet.Write(reqBuffer); err != nil {
		return nil, err
	}

	_, err = conn.Write(reqBuffer.buf[:])
	if err != nil {
		return nil, err
	}

	resBuffer := NewBytePacketBuffer()
	_, err = conn.Read(resBuffer.buf[:])
	if err != nil {
		return nil, err
	}

	resPacket, err := ReadDnsPacket(resBuffer)
	if err != nil {
		return nil, err
	}

	pp.Print(resPacket)

	return resPacket, nil
}

func RecursiveLookup(qname string, qtype QueryType) (*DnsPacket, error) {
    ns := net.ParseIP("192.5.5.241").To4()
    if ns == nil {
        return nil, errors.New("Failed to parse root server IP address")
    }

    for {
        fmt.Printf("attempting lookup of %v %s with ns %s\n", qtype, qname, ns.String())

        server := &net.UDPAddr{
            IP:   ns,
            Port: 53,
        }
        response, err := Lookup(qname, qtype, server)

        if err != nil {
            return nil, err
        }
        if len(response.Answers) > 0 && response.Header.ResCode == NOERROR {
            return response, nil
        }

        if response.Header.ResCode == NXDOMAIN {
            return response, nil
        }

        newNs := <-response.GetResolvedNs(qname)
        if newNs != nil {
            ns = newNs
            continue
        }

        newNsName := <-response.GetUnresolvedNs(qname)
        if newNsName == "" {
            return response, nil
        }

		AQueryType := NewQueryType(A,A)
        recursiveResponse, err := RecursiveLookup(newNsName, *AQueryType)
        if err != nil {
            return nil, err
        }

        newNs = <-recursiveResponse.GetRandomA()
        if newNs != nil {
            ns = newNs.To4()
        } else {
            return response, nil
        }
    }
}
