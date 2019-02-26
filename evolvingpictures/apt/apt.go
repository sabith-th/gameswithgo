package apt

import "math"

// Node is the basic node interface
type Node interface {
	Eval(x, y float32) float32
	String() string
}

// LeafNode is at the end, has no child
type LeafNode struct{}

// SingleNode has exactly one child
type SingleNode struct {
	Child Node
}

// DoubleNode has exactly two children - a left node and a right node
type DoubleNode struct {
	LeftChild  Node
	RightChild Node
}

// OpX is the operand x
type OpX LeafNode

// Eval returns value of operand x
func (op *OpX) Eval(x, y float32) float32 {
	return x
}

// String returns X
func (op *OpX) String() string {
	return "X"
}

// OpY is the operand x
type OpY LeafNode

// Eval returns value of operand y
func (op *OpY) Eval(x, y float32) float32 {
	return y
}

// String returns Y
func (op *OpY) String() string {
	return "Y"
}

// OpPlus is a double node which does addition
type OpPlus struct {
	DoubleNode
}

// Eval returns the sum of OpPlus nodes's two children
func (op *OpPlus) Eval(x, y float32) float32 {
	return op.LeftChild.Eval(x, y) + op.RightChild.Eval(x, y)
}

// String returns string x + y
func (op *OpPlus) String() string {
	return "( + " + op.LeftChild.String() + " " + op.RightChild.String() + " )"
}

// OpSin is the sin operator
type OpSin SingleNode

// Eval returns the sin of the child
func (op *OpSin) Eval(x, y float32) float32 {
	return float32(math.Sin(float64(op.Child.Eval(x, y))))
}

// String returns string Sin(x)
func (op *OpSin) String() string {
	return "( Sin " + op.Child.String() + " )"
}
