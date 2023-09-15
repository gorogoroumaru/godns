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
    return &DnsHeader{}
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
