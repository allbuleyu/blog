package framework

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strconv"
	"sync"
	"unicode"
)

var (
	bCommont = []byte{'#'}
	bEmpty = []byte{}
	bEqual = []byte{'='}
	bDQuote = []byte{'"'}
)


type Config struct {
	filename string
	comment  map[int][]string  	// id: []{comment, key...}; id 1 is for main comment.
	data map[string]string		// key: value
	offset   map[string]int64  	// key: offset; for editing.
	mu sync.RWMutex
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cfg := &Config{
		filename:file.Name(),
		comment: map[int][]string{},
		data: map[string]string{},
		offset: map[string]int64{},
		mu: sync.RWMutex{},
	}

	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	var comment bytes.Buffer

	buf := bufio.NewReader(file)

	for line, _, err := buf.ReadLine(); err != io.EOF; {
		if err != nil {
			return nil, err
		}
		nComment, off := 0, int64(1)

		if bytes.Equal(line, bEmpty) {
			continue
		}

		off += int64(len(line))

		if bytes.HasPrefix(line, bCommont) {
			line = bytes.TrimLeft(line, string(bCommont))
			line = bytes.TrimFunc(line, unicode.IsSpace)

			comment.Write(line)
			comment.WriteByte('\n')
			continue
		}

		if comment.Len() != 0 {
			cfg.comment[nComment] = []string{comment.String()}

			nComment++
			comment.Reset()
		}

		val := bytes.SplitN(line, bEqual, 2)
		if bytes.HasPrefix(val[1], bDQuote) {
			val[1] = bytes.Trim(val[1], string(bDQuote))
		}

		key := string(bytes.TrimSpace(val[0]))
		cfg.data[key] = string(bytes.TrimSpace(val[1]))

		cfg.comment[nComment-1] = append(cfg.comment[nComment-1], key)
		cfg.offset[key] = off
	}

	return cfg, nil
}


// Bool returns the boolean value for a given key.
func (c *Config) Bool(key string) (bool, error) {
	return strconv.ParseBool(c.data[key])
}

// Int returns the integer value for a given key.
func (c *Config) Int(key string) (int, error) {
	return strconv.Atoi(c.data[key])
}

// Float returns the float value for a given key.
func (c *Config) Float(key string) (float64, error) {
	return strconv.ParseFloat(c.data[key], 64)
}

// String returns the string value for a given key.
func (c *Config) String(key string) string {
	return c.data[key]
}