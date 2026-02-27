// This program takes the structured log output and makes it readable.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	// experimental package for json2 by go team
	"github.com/go-json-experiment/json"
)

var service string

func init() {
	flag.StringVar(&service, "service", "", "filter which service to see")
}

// Keep in mind that this is a separate go program which can be run from anywhere from my main program.
// you can pipe things from other programs to this program.
// In fact how I am using this right now is that I am piping the output to stdOut from other program
// to stdIn for this program.
func main() {
	flag.Parse()
	var b strings.Builder

	service := strings.ToLower(service)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		// read each line
		s := scanner.Text()

		// Create a map
		m := make(map[string]any)

		// unmarshall the json
		err := json.Unmarshal([]byte(s), &m)
		if err != nil {
			if service == "" {
				fmt.Println(s)
			}
			continue
		}

		// If a service filter was provided, check.
		if service != "" && strings.ToLower(m["service"].(string)) != service {
			continue
		}

		// I like always having a traceid present in the logs.
		traceID := "00000000-0000-0000-0000-000000000000"
		if v, ok := m["trace_id"]; ok {
			traceID = fmt.Sprintf("%v", v)
		}

		// {"time":"2023-06-01T17:21:11.13704718Z","level":"INFO","msg":"startup","service":"SALES-API","GOMAXPROCS":1}

		// Build out the know portions of the log in the order
		// I want them in.
		b.Reset()
		b.WriteString(fmt.Sprintf("%s: %s: %s: %s: %s: %s: ",
			m["service"],
			m["time"],
			m["file"],
			m["level"],
			traceID,
			m["msg"],
		))

		// Add the rest of the keys ignoring the ones we already
		// added for the log.
		for k, v := range m {
			switch k {
			case "service", "time", "file", "level", "trace_id", "msg":
				continue
			}

			// It's nice to see the key[value] in this format
			// especially since map ordering is random.
			b.WriteString(fmt.Sprintf("%s[%v]: ", k, v))
		}

		// Write the new log format, removing the last :
		out := b.String()
		fmt.Println(out[:len(out)-2])
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}
