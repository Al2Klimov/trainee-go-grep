// grep | (c) 2020 NETWAYS GmbH | GPLv2+

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

type arrayFlags []string

func (*arrayFlags) String() string {
	return ""
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	var flagTest arrayFlags
	equalComparator := false
	flag.Var(&flagTest, "e", "use PATTERN for matching.")
	fixedStringsParameter := flag.Bool("F", false, "use PATTERN not as a regular expression but as a string")
	flag.Parse()

	if len(flagTest) == 0 {
		fmt.Fprintln(os.Stderr, "grep: at least the -e parameter is needed.")
		os.Exit(1)
	}

	buf := bufio.NewReader(os.Stdin)
	var regArr []*regexp.Regexp

	if !*fixedStringsParameter {
		for _, i := range flagTest {
			re, regErr := regexp.Compile(i)
			if regErr != nil {
				fmt.Fprintln(os.Stderr, regErr)
				os.Exit(1)
			}
			regArr = append(regArr, re)
		}
	}

	for {
		data, dataErr := buf.ReadBytes('\n')
		if dataErr != nil && dataErr != io.EOF {
			fmt.Fprintln(os.Stderr, dataErr)
			os.Exit(1)
		}

		if len(data) == 0 {
			break
		}

		if *fixedStringsParameter {
			for _, i := range flagTest {
				if strings.Contains(string(data), i) {
					equalComparator = true
					break
				}
			}
		} else {
			for _, i := range regArr {
				if i.Match(data) {
					equalComparator = true
					break
				}
			}
		}

		if equalComparator {
			if !bytes.HasSuffix(data, []byte{'\n'}) {
				data = append(data, '\n')
			}
			_, _ = os.Stdout.Write(data)
			equalComparator = false
		}

		if dataErr == io.EOF {
			break
		}
	}
}
