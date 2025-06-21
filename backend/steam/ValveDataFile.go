package steam

import (
	"bufio"
	"errors"
	"io"
	"strings"
	"unicode"
)

// VDFNode represents a node in the VDF structure.
type VDFNode struct {
	Key      string
	Value    string
	Children []*VDFNode
}

// ParseVDF parses a Valve Data File (VDF) from an io.Reader.
func ParseVDF(r io.Reader) (*VDFNode, error) {
	scanner := bufio.NewScanner(r)
	var stack []*VDFNode
	var root *VDFNode

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		if line == "{" {
			continue
		}
		if line == "}" {
			if len(stack) > 1 {
				stack = stack[:len(stack)-1]
			}
			continue
		}
		key, value, err := parseVDFLine(line)
		if err != nil {
			return nil, err
		}
		node := &VDFNode{Key: key, Value: value}
		if len(stack) == 0 {
			root = node
			stack = append(stack, node)
		} else {
			parent := stack[len(stack)-1]
			parent.Children = append(parent.Children, node)
			if value == "" {
				stack = append(stack, node)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return root, nil
}

func parseVDFLine(line string) (string, string, error) {
	parts := splitVDFLine(line)
	if len(parts) == 1 {
		return parts[0], "", nil
	}
	if len(parts) == 2 {
		return parts[0], parts[1], nil
	}
	return "", "", errors.New("invalid VDF line: " + line)
}

func splitVDFLine(line string) []string {
	var parts []string
	var inQuote bool
	var sb strings.Builder
	for _, r := range line {
		switch {
		case r == '"':
			inQuote = !inQuote
			if !inQuote {
				parts = append(parts, sb.String())
				sb.Reset()
			}
		case inQuote:
			sb.WriteRune(r)
		case unicode.IsSpace(r):
			continue
		}
	}
	return parts
}
