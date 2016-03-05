package cson

import (
	"bytes"
	"encoding/json"
	"io"
	//	"log"
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
	needComma := false
	out := &bytes.Buffer{}

	s := raw
	i := 0
	nest := 1
	currentIndent := 0
	lastIndent := 0

	out.WriteByte('{')

	for i < len(s) {
		switch s[i] {
		case ' ', '\t':
			i++
			currentIndent++
		case '\n', '\r':
			i++
			currentIndent = 0
		case '#':
			comment := getComment(s[i:])
			i += len(comment)
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
			if currentIndent == -1 {
				// nop
			} else if currentIndent < lastIndent {
				nest--
				// close off object
				out.WriteByte('}')
				out.WriteByte(',')
				lastIndent = currentIndent
			} else if currentIndent == lastIndent {
				// bare word
				// could be a keyword, or a un-quoted string
				needComma = writeComma(out, needComma)
			} else {
				nest++
				// new object
				out.WriteByte('{')
				lastIndent = currentIndent
			}
			currentIndent = -1
			word := getWord(s[i:])
			writeWord(out, word, !isKeyword(word))
			i += len(word)
		}
	}

	for i := 0; i < nest; i++ {
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

// gets single or multiline comment
// ### means start/end multiline
// #, ##, or #### (or more) is a single line
//
func getComment(s []byte) []byte {
	// should never happen but be defensive
	if len(s) == 0 || s[0] != '#' {
		return nil
	}
	// get first line
	idx := bytes.IndexByte(s, '\n')
	if idx < 3 {
		// includes no ending newline
		return s
	}
	if s[1] != '#' || s[2] != '#' || s[3] == '#' {
		// single line comment
		// # ...
		// ## ...
		// ###x ...
		return s[:idx]
	}

	// multi-line
	idx = bytes.Index(s[4:], []byte("###"))
	if idx == -1 {
		// with no ending
		return s
	}
	return s[:idx+7]
}

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
