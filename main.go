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
	var counter int
	equalComparator := false
	flag.Var(&flagTest, "e", "use PATTERN for matching.")
	nInvertParameter := flag.Bool("v", false, "use PATTERN as non-matching lines.")
	fixedStringsParameter := flag.Bool("F", false, "use PATTERN not as a regular expression but as a string.")
	wordsParameter := flag.Bool("w", false, "use PATTERN that only matches words.")
	linesParameter := flag.Bool("x", false, "use PATTERN that only matches whole lines.")
	ignoreCaseParameter := flag.Bool("i", false, "ignore case distinctions")
	quietParameter := flag.Bool("q", false, "suppress all normal output")
	maxCountParameter := flag.Int("m", -1, "stop after NUM selected lines.")
	flag.Parse()

	if len(flagTest) == 0 {
		fmt.Fprintln(os.Stderr, "grep: at least the -e parameter is needed.")
		os.Exit(2)
	}

	buf := bufio.NewReader(os.Stdin)
	var regArr []*regexp.Regexp

	if *fixedStringsParameter {
		for i := range flagTest {
			flagTest[i] = regexp.QuoteMeta(flagTest[i])
		}
	}

	if *wordsParameter {
		for i := range flagTest {
			flagTest[i] = "\\b(" + flagTest[i] + ")\\b"
		}
	}

	if *linesParameter {
		for i := range flagTest {
			flagTest[i] = "(?m:^(" + flagTest[i] + ")$)"
		}
	}

	if *ignoreCaseParameter {
		for i := range flagTest {
			flagTest[i] = "(?i)" + flagTest[i]
		}
	}

	for _, i := range flagTest {
		re, regErr := regexp.Compile(i)
		if regErr != nil {
			fmt.Fprintln(os.Stderr, regErr)
			os.Exit(2)
		}
		regArr = append(regArr, re)
	}

	for {
		equalComparator = false
		data, dataErr := buf.ReadBytes('\n')
		if dataErr != nil && dataErr != io.EOF {
			fmt.Fprintln(os.Stderr, dataErr)
			os.Exit(2)
		}

		if len(data) == 0 {
			break
		}

		for _, i := range regArr {
			if i.Match(data) {
				equalComparator = true
				break
			}
		}

		if *maxCountParameter != -1 && counter == *maxCountParameter {
			break
		}

		if equalComparator != *nInvertParameter {
			if !bytes.HasSuffix(data, []byte{'\n'}) {
				data = append(data, '\n')
			}
			if !*quietParameter {
				_, _ = os.Stdout.Write(data)
			}
			counter++
		}

		if dataErr == io.EOF {
			break
		}
	}

	if counter <= 0 {
		os.Exit(1)
	}
}
