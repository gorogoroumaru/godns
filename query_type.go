package main

const (
	Unknown = iota
	A = 1
    NS = 2
    CNAME = 5
    MX = 15
    AAAA = 28
)

type QueryType struct {
	query_type	uint16
    val	uint16
}

func NewQueryType(query_type uint16, val uint16) *QueryType {
    return &QueryType{query_type, val}
}

func (qt QueryType) ToNum() uint16 {
    switch qt.query_type {
    case A:
        return 1
    case NS:
        return 2
    case CNAME:
        return 5
    case MX:
        return 15
    case AAAA:
        return 28
    default:
        return uint16(qt.val)
    }
}

func QueryTypeFromNum(num uint16) QueryType {
    switch num {
    case 1:
        return *NewQueryType(A, num)
    case 2:
        return *NewQueryType(NS, num)
    case 5:
        return *NewQueryType(CNAME, num)
    case 15:
        return *NewQueryType(MX, num)
    case 28:
        return *NewQueryType(AAAA, num)
    default:
        return *NewQueryType(Unknown, num)
    }
}
