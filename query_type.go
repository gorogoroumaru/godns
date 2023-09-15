package main

const (
	Unknown = iota
	A
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
    default:
        return uint16(qt.val)
    }
}

func QueryTypeFromNum(num uint16) QueryType {
    switch num {
    case 1:
        return *NewQueryType(A, num)
    default:
        return *NewQueryType(Unknown, num)
    }
}
