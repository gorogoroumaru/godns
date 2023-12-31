package main

type DnsHeader struct {
    ID                   uint16
    RecursionDesired     bool
    TruncatedMessage     bool
    AuthoritativeAnswer  bool
    Opcode               uint8
    Response             bool
    ResCode              ResultCode
    CheckingDisabled     bool
    AuthedData           bool
    Z                    bool
    RecursionAvailable   bool
    Questions            uint16
    Answers              uint16
    AuthoritativeEntries uint16
    ResourceEntries      uint16
}

func NewDnsHeader() *DnsHeader {
    return &DnsHeader{
        ID: 0,
        RecursionDesired: false,
        TruncatedMessage: false,
        AuthoritativeAnswer: false,
        Opcode: 0,
        Response: false,
        ResCode: NOERROR,
        CheckingDisabled: false,
        AuthedData: false,
        Z: false,
        RecursionAvailable: false,
        Questions: 0,
        Answers: 0,
        AuthoritativeEntries: 0,
        ResourceEntries: 0,
    }
}

func (h *DnsHeader) Read(buffer *BytePacketBuffer) error {
    var err error
    h.ID, err = buffer.ReadU16()
    if err != nil {
        return err
    }

    flags, err := buffer.ReadU16()
    if err != nil {
        return err
    }
    a := byte(flags >> 8)
    b := byte(flags & 0xFF)

    h.RecursionDesired = (a & (1 << 0)) > 0
    h.TruncatedMessage = (a & (1 << 1)) > 0
    h.AuthoritativeAnswer = (a & (1 << 2)) > 0
    h.Opcode = (a >> 3) & 0x0F
    h.Response = (a & (1 << 7)) > 0

    h.ResCode = ResultCode(b & 0x0F)
    h.CheckingDisabled = (b & (1 << 4)) > 0
    h.AuthedData = (b & (1 << 5)) > 0
    h.Z = (b & (1 << 6)) > 0
    h.RecursionAvailable = (b & (1 << 7)) > 0

    h.Questions, err = buffer.ReadU16()
    if err != nil {
        return err
    }
    h.Answers, err = buffer.ReadU16()
    if err != nil {
        return err
    }
    h.AuthoritativeEntries, err = buffer.ReadU16()
    if err != nil {
        return err
    }
    h.ResourceEntries, err = buffer.ReadU16()
    if err != nil {
        return err
    }

    return nil
}

func (header *DnsHeader) Write(buffer *BytePacketBuffer) error {
	if err := buffer.WriteU16(header.ID); err != nil {
		return err
	}

	header_byte_1 := uint8(0)
	if header.RecursionDesired {
		header_byte_1 |= 0x01
	}
	if header.TruncatedMessage {
		header_byte_1 |= 0x02
	}
	if header.AuthoritativeAnswer {
		header_byte_1 |= 0x04
	}
	header_byte_1 |= header.Opcode << 3
	if header.Response {
		header_byte_1 |= 0x80
	}

	if err := buffer.Write(header_byte_1); err != nil {
		return err
	}

	header_byte_2 := uint8(0)
	header_byte_2 |= (uint8(header.ResCode) & 0x0F)
	if header.CheckingDisabled {
		header_byte_2 |= 0x10
	}
	if header.AuthedData {
		header_byte_2 |= 0x20
	}
	if header.Z {
		header_byte_2 |= 0x40
	}
	if header.RecursionAvailable {
		header_byte_2 |= 0x80
	}

	if err := buffer.Write(header_byte_2); err != nil {
		return err
	}

	if err := buffer.WriteU16(header.Questions); err != nil {
		return err
	}

	if err := buffer.WriteU16(header.Answers); err != nil {
		return err
	}

	if err := buffer.WriteU16(header.AuthoritativeEntries); err != nil {
		return err
	}

	if err := buffer.WriteU16(header.ResourceEntries); err != nil {
		return err
	}

	return nil
}
