package xtorg

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	ErrScan = errors.New("scanner error")
)

func checkHeadline(s string) (*Headline, int) {
	level := 0
	for i, c := range s {
		if c != '*' {
			break
		}
		level = i + 1
	}
	if level == 0 || len(s) <= level+2 || s[level] != ' ' {
		return nil, 0
	}
	text := s[level+1:]
	//TODO parse todo and tags
	return &Headline{Text: text}, level
}

func parse(scanner *bufio.Scanner) (*Node, error) {

	root := &Node{
		value: &Subtree{
			Level: 0,
		},
	}

	chain := []*Node{root}
	raw := ""
	var parent *Node

	for scanner.Scan() {
		line := scanner.Text()
		h, level := checkHeadline(line)
		if h == nil {
			raw += line + "\n"
			continue
		}
		if raw != "" {
			parent = chain[len(chain)-1]
			node := &Node{
				parent: parent,
				idx:    len(parent.children),
				value:  &Raw{raw},
			}
			parent.children = append(parent.children, node)
			raw = ""
		}
		cutIdx := 1
		parent = root
		for i, node := range chain {
			if node.value.(*Subtree).Level < level {
				parent = node
				cutIdx = i + 1
			}
		}
		if cutIdx < len(chain) {
			chain = chain[:cutIdx]
		}
		node := &Node{
			parent: parent,
			idx:    len(parent.children),
			value: &Subtree{
				Level:    level,
				Headline: h,
			},
		}
		parent.children = append(parent.children, node)
		chain = append(chain, node)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrScan, err.Error())
	}
	if raw != "" {
		parent = chain[len(chain)-1]
		node := &Node{
			parent: parent,
			idx:    len(parent.children),
			value:  &Raw{raw},
		}
		parent.children = append(parent.children, node)
	}
	return root, nil
}

// Parse
func Parse(stream io.Reader) (*Node, error) {
	scanner := bufio.NewScanner(stream)
	return parse(scanner)
}

// ParseString
func ParseString(s string) (*Node, error) {
	scanner := bufio.NewScanner(strings.NewReader(s))
	return parse(scanner)
}

// ParseBytes
func ParseBytes(b []byte) (*Node, error) {
	scanner := bufio.NewScanner(bytes.NewReader(b))
	return parse(scanner)
}

// ParseFile
func ParseFile(name string) (*Node, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return parse(bufio.NewScanner(f))
}
