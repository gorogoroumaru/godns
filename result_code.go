package main

type ResultCode int

const (
    NOERROR  ResultCode = 0
    FORMERR  ResultCode = 1
    SERVFAIL ResultCode = 2
    NXDOMAIN ResultCode = 3
    NOTIMP   ResultCode = 4
    REFUSED  ResultCode = 5
)

func ResultCodeFromNum(num uint8) ResultCode {
    switch num {
    case 1:
        return FORMERR
    case 2:
        return SERVFAIL
    case 3:
        return NXDOMAIN
    case 4:
        return NOTIMP
    case 5:
        return REFUSED
    default:
        return NOERROR
    }
}
