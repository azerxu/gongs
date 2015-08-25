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

type SeqFiler interface {
	SeqIter
	Seqs() <-chan Seqer
	Close() error
	Err() error
}

type PairSeqer interface {
	GetRead1() Seqer
	GetRead2() Seqer
	String() string
}

type PairSeqIter interface {
	Next() bool
	Value() (Seqer, Seqer)
}

type PairSeqFiler interface {
	PairSeqIter
	Pairs() <-chan PairSeqer
	Close() error
	Err() error
}
