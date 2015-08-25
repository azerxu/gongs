package biofile

type Seqer interface {
	GetName() string
	GetSeq() []byte
	GetQual() []byte
	String() string
}

type SeqIter interface {
	Next() bool
	Value() (string, []byte, []byte)
}
