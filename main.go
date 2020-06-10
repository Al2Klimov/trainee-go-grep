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

var match bool
var regArr []*regexp.Regexp
var nInvertParameter = flag.Bool("v", false, "use PATTERN as non-matching lines.")
var quietParameter = flag.Bool("q", false, "suppress all normal output")
var maxCountParameter = flag.Int("m", -1, "stop after NUM selected lines.")

func main() {
	var flagTest arrayFlags
	flag.Var(&flagTest, "e", "use PATTERN for matching.")
	fixedStringsParameter := flag.Bool("F", false, "use PATTERN not as a regular expression but as a string.")
	wordsParameter := flag.Bool("w", false, "use PATTERN that only matches words.")
	linesParameter := flag.Bool("x", false, "use PATTERN that only matches whole lines.")
	ignoreCaseParameter := flag.Bool("i", false, "ignore case distinctions")
	flag.Parse()

	if len(flagTest) == 0 {
		fmt.Fprintln(os.Stderr, "grep: at least the -e parameter is needed.")
		os.Exit(2)
	}

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

	if len(flag.Args()) == 0 {
		compareAndPrint(os.Stdin)
	} else {
		for _, fileName := range flag.Args() {
			if fileName == "-" {
				compareAndPrint(os.Stdin)
			} else {
				file, fileErr := os.Open(fileName)
				if fileErr != nil {
					fmt.Fprintln(os.Stderr, fileErr)
					os.Exit(2)
				}

				compareAndPrint(file)
				file.Close()
			}
		}
	}

	if !match {
		os.Exit(1)
	}
}

func compareAndPrint(file *os.File) {
	counter := 0
	equalComparator := false
	buf := bufio.NewReader(file)

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

			match = true

			if !*quietParameter {
				if len(flag.Args()) > 1 {
					fmt.Printf("%s:", file.Name())
				}
				_, _ = os.Stdout.Write(data)
			}
			counter++
		}

		if dataErr == io.EOF {
			break
		}
	}
}
