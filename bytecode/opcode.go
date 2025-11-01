package bytecode

type Opcode int

const (
	OP_RETURN = Opcode(iota)
)

type Chunk struct {
	Count    int
	Capacity int
	Code     Opcode
}
