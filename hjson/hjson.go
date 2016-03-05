package hjson

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
)

type readerState struct {
	source io.Reader
	br     *bytes.Reader
}

// New returns an io.Reader that converts a HJSON input to JSON
func New(r io.Reader) io.Reader {
	return &readerState{source: r}
}

// Read implements the io.Reader interface
func (st *readerState) Read(p []byte) (int, error) {
	if st.br == nil {
		buf := &bytes.Buffer{}
		if _, err := io.Copy(buf, st.source); err != nil {
			return 0, err
		}
		st.br = bytes.NewReader(ToJSON(buf.Bytes()))
	}
	return st.br.Read(p)
}

// Unmarshal is the same as JSON.Unmarshal but for HJSON files
func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(ToJSON(data), v)
}

// ToJSON converts a hjson format to JSON
func ToJSON(raw []byte) []byte {
	needEnding := false
	needComma := false
	out := &bytes.Buffer{}

	s := raw
	i := 0

	// skip over initial whitespace.
	// if first char is NOT a '{' then add it
	for i < len(s) {
		if isWhitespace(s[i]) {
			i++
		} else if s[i] == '{' {
			break
		} else {
			out.WriteByte('{')
			needEnding = true
			break
		}
	}

	for i < len(s) {
		switch s[i] {
		case ' ', '\n', '\t', '\r':
			i++
		case ':':
			// next value does not need an auto-comma
			needComma = false
			out.WriteByte(':')
			i++
		case '{':
			writeComma(out, needComma)
			needComma = false
			out.WriteByte('{')
			i++
		case '[':
			writeComma(out, needComma)
			needComma = false
			out.WriteByte('[')
			i++
		case '}':
			// next value may need a comma, e.g. { ...},{...}
			needComma = true
			out.WriteByte('}')
			i++
		case ']':
			// next value may need a comma, e.g. { ...},{...}
			needComma = true
			out.WriteByte(']')
			i++
		case '/':
			if i+1 < len(s) && s[i+1] == '/' {
				idx := bytes.IndexByte(s[2:], '\n')
				if idx == -1 {
					i = len(s)
				} else {
					i += idx + 2
				}
			} else if i+1 < len(s) && s[i+1] == '*' {
				idx := bytes.Index(s[2:], []byte("*/"))
				if idx == -1 {
					i = len(s)
				} else {
					i += idx + 4
				}
			} else {
				// bare word
				needComma = writeComma(out, needComma)
				word := getWord(s[i:])
				writeWord(out, word, !isKeyword(word))
				i += len(word)
			}
		case '#':
			// Scan to EOL
			for ; i < len(s); i++ {
				if s[i] == '\n' {
					break
				}
			}
		case ',':
			// we pretend we didn't see this and let the auto-comma code add it if necessary
			// if the next token is value, it will get added
			// if the next token is a '}' or '], then it will NOT get added (fixes ending comma problem in JSON)
			needComma = true
			i++
		case '\'', '"':
			needComma = writeComma(out, needComma)
			content, offset := getString(s[i:])
			out.WriteByte('"')
			out.Write(content)
			out.WriteByte('"')
			i += offset
		case '+', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			needComma = writeComma(out, needComma)
			word := getWord(s[i:])
			// captured numeric input... does it parse as a number?
			// if not, then quote it
			_, err := strconv.ParseFloat(string(word), 64)
			writeWord(out, word, err != nil)
			i += len(word)
		default:
			// bare word
			// could be a keyword, or a un-quoted string
			needComma = writeComma(out, needComma)
			word := getWord(s[i:])
			writeWord(out, word, !isKeyword(word))
			i += len(word)
		}
	}

	if needEnding {
		out.WriteByte('}')
	}

	return out.Bytes()
}

func isWhitespace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

func isDelimiter(c byte) bool {
	return c == ':' || c == '}' || c == ']' || c == ',' || c == '\n'
}

func getWord(s []byte) []byte {
	for j := 0; j < len(s); j++ {
		if isDelimiter(s[j]) {
			return bytes.TrimSpace(s[:j])
		}
	}
	return s
}

func isKeyword(s []byte) bool {
	return bytes.Equal(s, []byte("false")) || bytes.Equal(s, []byte("true")) || bytes.Equal(s, []byte("null"))
}

func writeComma(buf *bytes.Buffer, comma bool) bool {
	if comma {
		buf.WriteByte(',')
	}
	return true
}

func writeWord(buf *bytes.Buffer, word []byte, quote bool) {
	if quote {
		buf.WriteByte('"')
	}

	// to JS escape word
	buf.Write(word)

	if quote {
		buf.WriteByte('"')
	}
}

// handles single line and multi-line strings
func getString(s []byte) ([]byte, int) {
	if len(s) == 0 {
		return nil, 0
	}
	char := s[0]
	if char != '\'' && char != '"' {
		return nil, 0
	}
	if len(s) > 3 && s[1] == char && s[2] == char {
		// we have multi-line

		// assume not ended correctly
		offset := len(s)
		content := s[3:]

		idx := bytes.Index(content, []byte{char, char, char})
		if idx > -1 {
			// with ending
			content = content[:idx]
			offset = idx + 7
		}
		// now figure out whitespace stuff
		if len(content) > 0 && content[0] == '\n' {
			content = content[1:]
		}
		if len(content) > 0 && content[len(content)-1] == '\n' {
			content = content[:len(content)-1]
		}
		minIndent := len(content)
		lines := bytes.Split(content, []byte{'\n'})
		for _, line := range lines {
			for i := 0; i < len(line) && i < minIndent; i++ {
				if line[i] != ' ' {
					minIndent = i
					break
				}
			}
		}

		if minIndent > 0 {
			for i, line := range lines {
				lines[i] = line[minIndent:]
			}
		}
		content = bytes.Join(lines, []byte{'\\', 'n'})
		return content, offset
	}

	// single line string
	j := 1
	for j < len(s) {
		if s[j] == char {
			break
		} else if s[j] == '\\' && j+1 < len(s) {
			j++
		}
		j++
	}

	// not sure if other things need replacing or not
	content := s[1:j]
	content = bytes.Replace(content, []byte{'\n'}, []byte{'\\', 'n'}, -1)
	return content, j + 1
}
