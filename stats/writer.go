package stats

import (
	"bytes"
	"flag"
	"fmt"
	"sync"
	"time"
)

var printStats = flag.Bool("printStats", false, "Print stats to console")

var bufPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

// IncCounter increments a counter.
func IncCounter(name string, tags map[string]string, value int64) {
	name = addTagsToName(name, tags)
	if *printStats {
		fmt.Printf("IncCounter: %v = %v\n", name, value)
	}
}

// RecordTimer records a timer.
func RecordTimer(name string, tags map[string]string, d time.Duration) {
	name = addTagsToName(name, tags)
	if *printStats {
		fmt.Printf("RecordTimer: %v = %v\n", name, d)
	}
}

func addTagsToName(name string, tags map[string]string) string {
	// The format we want is: host.endpoint.os.browser
	// if there's no host tag, then we don't use it.
	keyOrder := make([]string, 0, 4)
	if _, ok := tags["host"]; ok {
		keyOrder = append(keyOrder, "host")
	}
	keyOrder = append(keyOrder, "endpoint", "os", "browser")

	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()

	buf.WriteString(name)
	for _, k := range keyOrder {
		buf.WriteByte('.')
		v, ok := tags[k]
		if !ok || v == "" {
			buf.WriteString("no-")
			buf.WriteString(k)
			continue
		}
		writeClean(buf, v)
	}

	final := buf.String()
	bufPool.Put(buf)
	return final
}

// writeClean takes a string that may contain special characters, and replaces these
// characters with a '-'.
func writeClean(buf *bytes.Buffer, value string) {
	for i := 0; i < len(value); i++ {
		switch c := value[i]; c {
		case '[', '{', '}', '/', '\\', ':', ' ', '\t', '.':
			buf.WriteByte('-')
		default:
			buf.WriteByte(c)
		}
	}
}
