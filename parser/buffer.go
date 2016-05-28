package parser

import (
	"bytes"
)

type Buffer struct {
	*bytes.Buffer
}

func NewBuffer(data []byte) *Buffer {
	buf := &Buffer{}
	buf.Buffer = bytes.NewBuffer(data)
	return buf
}

func (buf *Buffer) removeBytes(count int) {
	buf.Buffer = bytes.NewBuffer(buf.Bytes()[:buf.Len()-count])
}

// Grab n bytes (length of token) from end of buf and compare to token
func (buf *Buffer) checkForToken(token string) bool {
	//fmt.Printf("checkBufForToken('%s', '%s')\n", buf.String(), token)
	token_len := len(token)
	if token_len == 0 {
		return false
	}
	buf_len := buf.Len()
	if buf_len < token_len {
		return false
	}
	//fmt.Printf("buf_len = %d, token_len = %d\n", buf_len, token_len)
	buf_bytes := buf.Bytes()[buf.Len()-token_len:]
	for i, b := range []byte(token) {
		if buf_bytes[i] != b {
			return false
		}
	}
	// Look backwards through the buffer to see if the beginning of our
	// token was escaped
	escape := false
	buf_bytes = buf.Bytes()
	for i := (buf_len - token_len - 1); i >= 0; i-- {
		//fmt.Printf("buf is '%s', char at %d is '%c'\n", buf_bytes, i, buf_bytes[i])
		if buf_bytes[i] == '\\' {
			//fmt.Printf("flipping 'escape' from %t to %t\n", escape, !escape)
			escape = !escape
		} else {
			break
		}
	}
	//fmt.Printf("'escape' is %t\n", escape)
	if escape {
		return false
	}
	return true
}
