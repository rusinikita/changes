package git

import (
	"strings"

	"github.com/go-git/go-git/v5/plumbing/format/diff"
)

// FileChange represents the necessary steps to transform one file into another.
type FileChange struct {
	// Path is file path after change
	// empty if file deleted
	Path string
	// PrevPath is file path before change
	// empty if file created
	// diff from Path if file renamed
	PrevPath string
	// Chunks returns a slice of ordered changes to transform "from" File into
	// "to" File. If the file is a binary one, Chunks will be empty.
	Chunks []Chunk
}

// Chunk represents a portion of a file transformation into another.
type Chunk struct {
	// Content contains the portion of the file.
	Type Operation
	// Type contains the Operation to do with this Chunk.
	Content string
}

type Stat struct {
	Additions, Deletions int
}

func (c FileChange) Stats() (s Stat) {
	var (
		content string
		lines   int
	)

	for _, chunk := range c.Chunks {
		content = chunk.Content

		lines = strings.Count(content, "\n")
		if content[len(content)-1] != '\n' {
			lines++
		}

		switch chunk.Type {
		case Add:
			s.Additions += lines
		case Delete:
			s.Deletions += lines
		}
	}

	return s
}

// Operation defines the operation of a diff item.
type Operation interface {
	_operation()
}

const (
	// Equal item represents an equals diff.
	Equal operation = iota
	// Add item represents an insert diff.
	Add
	// Delete item represents a delete diff.
	Delete
)

type operation int

func (operation) _operation() {}

func (g *change) FilesDiff() (result []FileChange, err error) {
	patch, err := g.repository.CommonCommit.Patch(g.repository.SourceLastCommit)
	if err != nil {
		return nil, err
	}

	for _, file := range patch.FilePatches() {
		change := FileChange{}

		from, to := file.Files()

		if to != nil {
			change.Path = to.Path()
		}

		if from != nil {
			change.PrevPath = from.Path()
		}

		change.Chunks = make([]Chunk, 0, len(file.Chunks()))
		for _, c := range file.Chunks() {
			change.Chunks = append(change.Chunks, convertChunk(c))
		}

		result = append(result, change)
	}

	return result, nil
}

func convertChunk(c diff.Chunk) Chunk {
	t := Equal

	switch c.Type() {
	case diff.Add:
		t = Add
	case diff.Delete:
		t = Delete
	}

	return Chunk{
		Type:    t,
		Content: c.Content(),
	}
}
