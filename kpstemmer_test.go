package kpstemmer

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKpStemmer(t *testing.T) {

	deviationMap := loadDeviationMap()

	fi, err := os.Open("test_diffs.txt")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()

	count := 0
	bufioReader := bufio.NewReader(fi)
	line, err := bufioReader.ReadString('\n')
	for err == nil {
		if len(strings.TrimSpace(line)) > 0 {
			runes := bytes.Runes([]byte(line))
			var token, expected string
			if runes[29] == ' ' {
				token = strings.TrimRight(string(runes[:30]), " ")
				sub2 := runes[30:]
				if len(sub2) >= 31 && sub2[30] == '*' {
					sub2 = sub2[32:]
				}
				expected = strings.TrimRight(strings.TrimSuffix(string(sub2), "\n"), " ")
			} else {
				token = strings.TrimSuffix(line, "\n")
				line, err = bufioReader.ReadString('\n')
				assert.NoError(t, err)
				runes = bytes.Runes([]byte(line))
				expected = strings.TrimSuffix(string(runes[30:]), "\n")
			}

			result := Stem(token)
			if deviation, ok := deviationMap[token]; ok {
				expected = deviation
			}
			assert.Equal(t, expected, result, `stemmed:  "%s"`, token)
			count += 1
		}
		line, err = bufioReader.ReadString('\n')
	}
	if err != io.EOF {
		assert.NoError(t, err)
	}
	fmt.Printf("Processed words: %d\n", count)
}

func loadDeviationMap() map[string]string {
	deviationMap := map[string]string{}
	bufioReader := bufio.NewReader(bytes.NewReader(deviationData))
	line, err := bufioReader.ReadString('\n')
	for err == nil {
		tokens := strings.Fields(line)
		if len(tokens) >= 3 {
			deviationMap[tokens[0]] = tokens[2]
		}
		line, err = bufioReader.ReadString('\n')
	}
	return deviationMap
}

var deviationData = []byte(`
emiliètje		emiliè		emilièt
liètje			liè			lièt
mariètje		mariè		marièt
ruïneren		ruïner		ruïneer
`)
