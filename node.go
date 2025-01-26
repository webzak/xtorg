package xtorg

import "errors"

var (
	ErrNodeDoesNotExist  = errors.New("node does not exist")
	ErrValueIsNotSubtree = errors.New("node value is not subtree")
	ErrValueIsNotRaw     = errors.New("node value is not raw")
)

var (
	undef = &Node{}
)

// Type contains node type
type Type int

const (
	// TRaw represents raw text, that is not classified or parsed
	TRaw Type = iota
	// TSubtree represents the part of text started with headline
	TSubtree
	// THeadline is section title
	THeadline
	// TParagraph is the text paragraph
	TParagraph
	// TBlock is the contents between #+BEGIN_<type>’  ‘#+END’
	TBlock
	// TList represents ordered or unordered list
	TList
	// TTable represents table
	TTable
	// TDrawer represents drawer
	TDrawer
	// TUndef represents special value
	TUndef
)

type Raw struct {
	Text string
}

type Subtree struct {
	Level      int
	Headline   *Headline
	Planning   *Planning
	Properties map[string]string
}

type Headline struct {
	Todo string
	Text string
	Tags []string
}

type Planning struct {
}

type Paragraph struct {
}

type Block struct {
}

type List struct {
}

type Table struct {
}

type Drawer struct {
}

// Node structure represents the element of parsed org document
type Node struct {
	parent   *Node
	idx      int
	value    any
	children []*Node
}

func (n *Node) Type() Type {
	if n == nil || n == undef {
		return TUndef
	}
	switch n.value.(type) {
	case *Raw:
		return TRaw
	case *Subtree:
		return TSubtree
	case *Paragraph:
		return TParagraph
	case *Block:
		return TBlock
	case *List:
		return TList
	case *Table:
		return TTable
	case *Drawer:
		return TDrawer
	default:
		panic("node value type is not supported")
	}
}

func (n *Node) IsSubtree() bool {
	return n.Type() == TSubtree
}

func (n *Node) Subtree() (*Subtree, error) {
	if n == nil || n == undef {
		return nil, ErrNodeDoesNotExist
	}
	v, ok := n.value.(*Subtree)
	if !ok {
		return nil, ErrValueIsNotSubtree
	}
	return v, nil
}

func (n *Node) IsRaw() bool {
	return n.Type() == TSubtree
}

func (n *Node) Raw() (*Raw, error) {
	if n == nil || n == undef {
		return nil, ErrNodeDoesNotExist
	}
	v, ok := n.value.(*Raw)
	if !ok {
		return nil, ErrValueIsNotRaw
	}
	return v, nil
}

func (n *Node) Level() int {
	if n == nil || n == undef {
		return -1
	}
	st, err := n.Subtree()
	if err == nil {
		return st.Level
	}
	if n.parent == nil {
		return 0
	}
	st, err = n.parent.Subtree()
	if err == nil {
		return st.Level
	}
	return -1
}
