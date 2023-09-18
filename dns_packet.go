package main

import (
	"net"
	"strings"
)

type DnsPacket struct {
	Header       *DnsHeader
	Questions    []*DnsQuestion
	Answers      []DnsRecord
	Authorities  []DnsRecord
	Resources    []DnsRecord
}

func NewDnsPacket() *DnsPacket {
	return &DnsPacket{
		Header:       NewDnsHeader(),
		Questions:    []*DnsQuestion{},
		Answers:      []DnsRecord{},
		Authorities:  []DnsRecord{},
		Resources:    []DnsRecord{},
	}
}

func ReadDnsPacket(buffer *BytePacketBuffer) (*DnsPacket, error) {
	header := NewDnsHeader()
	if err := header.Read(buffer); err != nil {
		return nil, err
	}

	dnsPacket := NewDnsPacket()
	dnsPacket.Header = header

	for i := 0; i < int(header.Questions); i++ {
		question := &DnsQuestion{}
		if err := question.Read(buffer); err != nil {
			return nil, err
		}
		dnsPacket.Questions = append(dnsPacket.Questions, question)
	}

	for i := 0; i < int(header.Answers); i++ {
		record, err := ReadDnsRecord(buffer)
		if err != nil {
			return nil, err
		}
		dnsPacket.Answers = append(dnsPacket.Answers, record)
	}

	for i := 0; i < int(header.AuthoritativeEntries); i++ {
		record, err := ReadDnsRecord(buffer)
		if err != nil {
			return nil, err
		}
		dnsPacket.Authorities = append(dnsPacket.Authorities, record)
	}

	for i := 0; i < int(header.ResourceEntries); i++ {
		record, err := ReadDnsRecord(buffer)
		if err != nil {
			return nil, err
		}
		dnsPacket.Resources = append(dnsPacket.Resources, record)
	}

	return dnsPacket, nil
}

func (packet *DnsPacket) Write(buffer *BytePacketBuffer) error {
	packet.Header.Questions = uint16(len(packet.Questions))
	packet.Header.Answers = uint16(len(packet.Answers))
	packet.Header.AuthoritativeEntries = uint16(len(packet.Authorities))
	packet.Header.ResourceEntries = uint16(len(packet.Resources))

	if err := packet.Header.Write(buffer); err != nil {
		return err
	}

	for _, question := range packet.Questions {
		if err := question.Write(buffer); err != nil {
			return err
		}
	}

	for _, rec := range packet.Answers {
		if _, err := rec.Write(buffer); err != nil {
			return err
		}
	}

	for _, rec := range packet.Authorities {
		if _, err := rec.Write(buffer); err != nil {
			return err
		}
	}

	for _, rec := range packet.Resources {
		if _, err := rec.Write(buffer); err != nil {
			return err
		}
	}

	return nil
}

func (p *DnsPacket) GetRandomA() <-chan net.IP {
	resultChan := make(chan net.IP)

	go func() {
		defer close(resultChan)

		for _, record := range p.Answers {
			if record.getType() == A {
				resultChan <- record.(*ARecord).Addr
			}
		}
	}()

	return resultChan
}

func (p *DnsPacket) GetNs(qname string) []string {
	nsList := make([]string, 0)
	for _, record := range p.Authorities {
		if record.getType() == NS && strings.HasSuffix(qname, record.(*NSRecord).Domain){
			nsList = append(nsList, record.(*NSRecord).Host)
		}
	}
	return nsList
}

func (p *DnsPacket) GetResolvedNs(qname string) <-chan net.IP {
	resultChan := make(chan net.IP)

	go func() {
		defer close(resultChan)

		nsList := p.GetNs(qname)
		for _, record := range p.Resources {
			if record.getType() == A && contains(nsList, record.(*ARecord).Domain) {
				resultChan <- record.(*ARecord).Addr
			}
		}
	}()

	return resultChan
}

func (p *DnsPacket) GetUnresolvedNs(qname string) <-chan string {
	resultChan := make(chan string)

	go func() {
		defer close(resultChan)

		nsList := p.GetNs(qname)
		for _, ns := range nsList {
			resultChan <- ns
		}
	}()

	return resultChan
}

func contains(slice []string, s string) bool {
	for _, str := range slice {
		if str == s {
			return true
		}
	}
	return false
}
