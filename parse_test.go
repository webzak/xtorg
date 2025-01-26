package xtorg

import "testing"

func TestParse(t *testing.T) {
	raw := `
#+TITLE: Sample document

Occationally it has some text
at the top level.

* First level
There is a description of the first level
** Second level one :foo:boo:
:PROPERTIES:
:CUSTOM_ID: initialization
:END:
The start of the content

** Second level two
*** Third level two - one
Some text about it
`
	root, err := ParseString(raw)
	if err != nil {
		t.Fatalf("ParseString error: %v", err)
	}
	if !root.IsSubtree() {
		t.Fatal("Root element must be subtree")
	}
	if root.Level() != 0 {
		t.Fatalf("Root level expected to be 0, got: %d", root.Level())
	}
	if len(root.children) != 2 {
		t.Fatalf("Expected 2 children, but got: %d", len(root.children))
	}
	if root.children[0].Type() != TRaw {
		t.Fatalf("Expected TRaw type. Got: %d", root.children[0].Type())
	}

	st1, err := root.children[1].Subtree()
	if err != nil {
		t.Fatalf("Second root child expected to be subtree, got error: %v", err)
	}
	if st1.Level != 1 {
		t.Fatalf("Level expected to be 1, got: %d", st1.Level)
	}
	if st1.Headline.Text != "First level" {
		t.Fatal("Wrong text for top header")
	}

	st2, err := root.children[1].children[2].Subtree()
	if err != nil {
		t.Fatalf("Expected to be subtree, got error: %v", err)
	}
	if st2.Level != 2 {
		t.Fatalf("Level expected to be 2, got: %d", st2.Level)
	}
	if st2.Headline.Text != "Second level two" {
		t.Fatal("Wrong text for header")
	}

	st3, err := root.children[1].children[2].children[0].Subtree()
	if err != nil {
		t.Fatalf("Expected to be subtree, got error: %v", err)
	}
	if st3.Level != 3 {
		t.Fatalf("Level expected to be 3, got: %d", st3.Level)
	}
	if st3.Headline.Text != "Third level two - one" {
		t.Fatal("Wrong text for header")
	}
}
