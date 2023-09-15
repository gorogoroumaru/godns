package main

import (
	"net"
)

type DnsRecord interface {
	isDnsRecord()
}

type UnknownRecord struct {
	Domain  string
	QType   uint16
	DataLen uint16
	TTL     uint32
}

func (u *UnknownRecord) isDnsRecord() {}

type ARecord struct {
	Domain string
	Addr   net.IP
	TTL    uint32
}

func (a *ARecord) isDnsRecord() {}

func ReadDnsRecord(buffer *BytePacketBuffer) (DnsRecord, error) {
	var domain string
	if err := buffer.ReadQName(&domain); err != nil {
		return nil, err
	}

	qTypeNum, err := buffer.ReadU16()
	if err != nil {
		return nil, err
	}
	qType := QueryTypeFromNum(qTypeNum)

	if _, err := buffer.ReadU16(); err != nil {
		return nil, err
	}

	ttl, err := buffer.ReadU32()
	if err != nil {
		return nil, err
	}

	dataLen, err := buffer.ReadU16()
	if err != nil {
		return nil, err
	}

	switch qType.query_type {
	case A:
		rawAddr, err := buffer.ReadU32()
		if err != nil {
			return nil, err
		}
		addr := net.IPv4(
			byte(rawAddr>>24),
			byte(rawAddr>>16),
			byte(rawAddr>>8),
			byte(rawAddr),
		)
		return &ARecord{
			Domain: domain,
			Addr:   addr,
			TTL:    ttl,
		}, nil
	default:
		if err := buffer.Step(uint16(dataLen)); err != nil {
			return nil, err
		}
		return &UnknownRecord{
			Domain:  domain,
			QType:   qTypeNum,
			DataLen: dataLen,
			TTL:     ttl,
		}, nil
	}
}
