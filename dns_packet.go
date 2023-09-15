package main

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
