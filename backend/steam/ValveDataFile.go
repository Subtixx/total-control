package steam

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
	"unicode"
)

func ReadVDF(filePath string) (map[string]interface{}, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return ParseVDF(data)
}

// ParseVDF parses a Valve Data Format (VDF) file and returns a map[string]interface{}.
func ParseVDF(data []byte) (map[string]interface{}, error) {
	reader := bufio.NewReader(bytes.NewReader(data))
	return parseObject(reader)
}

func parseObject(r *bufio.Reader) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	var key string
	for {
		token, err := readToken(r)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if token == "}" {
			break
		}
		key = token
		token, err = readToken(r)
		if err != nil {
			return nil, err
		}
		if token == "{" {
			val, err := parseObject(r)
			if err != nil {
				return nil, err
			}
			result[key] = val
		} else {
			result[key] = token
		}
	}
	return result, nil
}

func readToken(r *bufio.Reader) (string, error) {
	var sb strings.Builder
	inQuotes := false
	for {
		ch, _, err := r.ReadRune()
		if err != nil {
			return "", err
		}
		if unicode.IsSpace(ch) && !inQuotes {
			if sb.Len() > 0 {
				break
			}
			continue
		}
		if ch == '"' {
			if inQuotes {
				break
			}
			inQuotes = true
			continue
		}
		if !inQuotes && (ch == '{' || ch == '}') {
			if sb.Len() == 0 {
				return string(ch), nil
			}
			err := r.UnreadRune()
			if err != nil {
				return "", err
			}
			break
		}
		sb.WriteRune(ch)
	}
	return sb.String(), nil
}
