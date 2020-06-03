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

func main() {
	patternParameter := flag.String("e", "", "use PATTERN for matching.")
	flag.Parse()

	if *patternParameter == "" {
		fmt.Fprintln(os.Stderr, "grep: at least the -e parameter is needed.")
		os.Exit(1)
	}

	re, regErr := regexp.Compile(*patternParameter)
	if regErr != nil {
		fmt.Fprintln(os.Stderr, regErr)
		os.Exit(1)
	}

	buf := bufio.NewReader(os.Stdin)

	for {
		data, dataErr := buf.ReadBytes('\n')
		if dataErr != nil && dataErr != io.EOF {
			fmt.Fprintln(os.Stderr, dataErr)
			os.Exit(1)
		}

		if len(data) == 0 {
			break
		}

		if re.Match(data) {
			if !bytes.HasSuffix(data, []byte{'\n'}) {
				data = append(data, '\n')
			}
			_, _ = os.Stdout.Write(data)
		}

		if dataErr == io.EOF {
			break
		}
	}
}
