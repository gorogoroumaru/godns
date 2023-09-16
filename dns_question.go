package main

type DnsQuestion struct {
    Name  string
    QType QueryType
}

func NewDnsQuestion(name string, qType QueryType) *DnsQuestion {
    return &DnsQuestion{
        Name:  name,
        QType: qType,
    }
}

func (q *DnsQuestion) Read(buffer *BytePacketBuffer) error {
    if err := buffer.ReadQName(&q.Name); err != nil {
        return err
    }
    qTypeNum, err := buffer.ReadU16()
    if err != nil {
        return err
    }
    q.QType = QueryTypeFromNum(qTypeNum)

    _, err = buffer.ReadU16()
    if err != nil {
        return err
    }

    return nil
}

func (q *DnsQuestion) Write(buffer *BytePacketBuffer) error {
	if err := buffer.WriteQName(&q.Name); err != nil {
        return err
    }

	typeNum := q.QType.ToNum()
	if err :=  buffer.WriteU16(typeNum); err != nil {
        return err
    }

	if err :=  buffer.WriteU16(1); err != nil {
        return err
    }

	return nil
}