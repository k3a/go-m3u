package m3u

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

const m3uHeader = "#EXTM3U"
const newLine = "\n"

// Record is a single record of M3U file
type Record struct {
	Duration   float64
	Attributes map[string]string
	Title      string
	URL        string
}

// NewRecord creates a new empty M3U record instance
func NewRecord() *Record {
	r := new(Record)
	r.Attributes = make(map[string]string)
	return r
}

// M3U holds m3u file with records
type M3U []*Record

func newM3U() *M3U {
	return nil
}

func (m *M3U) Write(wr io.Writer) (err error) {
	if _, err = wr.Write([]byte(m3uHeader + newLine)); err != nil {
		return
	}

	for _, r := range *m {
		// header with duration
		extInf := fmt.Sprintf("#EXTINF:%.0f", r.Duration)

		// attrribute list
		for ak, av := range r.Attributes {
			extInf += fmt.Sprintf(" %s=\"%s\"", ak, av)
		}

		// title and line ending
		extInf += fmt.Sprintf(",%s"+newLine, r.Title)

		// write EXTINF
		if _, err = wr.Write([]byte(extInf)); err != nil {
			return
		}

		// write URL
		if _, err = wr.Write([]byte(r.URL + newLine)); err != nil {
			return
		}
	}

	return
}

func (m *M3U) String() string {
	dbg := ""
	for _, r := range *m {
		dbg += fmt.Sprintf("- title: %s, duration: %.0f, attrs: %v\n  %s\n",
			r.Title, r.Duration, r.Attributes, r.URL)
	}
	return dbg
}

// Add adds a new record to M3U
func (m *M3U) Add(r *Record) {
	*m = append(*m, r)
}

// Records returns mutable list of M3U records
func (m *M3U) Records() []*Record {
	return *m
}

var m3uExtInfRE = regexp.MustCompile(`#EXTINF:(\S+)\s*(.*),(.*)`)
var m3uAttrRE = regexp.MustCompile(`(.+)="([^"]*)"`)

func (m *M3U) parseExtInf(r *Record, line string) (err error) {
	matches := m3uExtInfRE.FindStringSubmatch(line)
	if len(matches) < 4 {
		return fmt.Errorf("m3u: wrongly-formatted line: %s", line)
	}

	// duration
	r.Duration, err = strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return
	}

	for _, attr := range m3uAttrRE.FindAllString(matches[2], -1) {
		// clean attribute pair
		attr = strings.TrimSpace(attr)
		if len(attr) == 0 {
			// skip empty attribute
			continue
		}

		// parse attribute to key-value pair
		am := m3uAttrRE.FindStringSubmatch(attr)
		if len(am) <= 1 {
			return fmt.Errorf("m3u: wrongly-formatted attribute '%s' on the line '%s'", attr, line)
		}

		// store attribute
		if len(am) > 2 {
			r.Attributes[am[1]] = am[2]
		} else {
			r.Attributes[am[1]] = ""
		}
	}

	// title
	r.Title = matches[3]

	return
}

func (m *M3U) Read(rdr io.Reader) error {
	scanner := bufio.NewScanner(rdr)

	r := NewRecord()

	for scanner.Scan() {
		line := scanner.Text()

		// skip header
		if line == m3uHeader {
			continue
		} else if strings.HasPrefix(line, "#EXTINF") {
			// EXTINF
			if err := m.parseExtInf(r, line); err != nil {
				return err
			}
		} else if strings.HasPrefix(line, "#") {
			// unknown extension
			return fmt.Errorf("m3u: unknown line %s", line)
		} else {
			// pass-through URL and finish record
			r.URL = line
			*m = append(*m, r)
			r = NewRecord()
		}

	}

	// parsing error
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
