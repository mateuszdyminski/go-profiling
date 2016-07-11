*** Live-coding session:

- run http server and enter urls in chrome
- run http server with -printStats flag
- enter urls in chrome - show logs from http server

```
boom -c 10 -n 100000 http://localhost:9090/hello
boom -c 10 -n 100000 http://localhost:9090/simple
```

- run long job 
```
boom -c 10 -n 100000000 http://localhost:9090/simple
```

goto http://localhost:9090/debug/pprof/

- show goroutines
- show heap - not interesting at top but nice at bottom - GC num, gc times, alloc size

-  start profiling

```
go tool pprof -seconds 5 http://localhost:9090/debug/pprof/profile
```

pprof - top10
pprof - top10 -cum
pprof - web

- go-torch example:

```
go-torch -seconds 5 http://localhost:9090/debug/pprof/profile
```
open torch.svg

- it shows that the biggest problem is in os.Hostname() call in getStatsTag() func
- move it outside the func and call only once at the beginning
- start next perf tests

- run new long run test
- do next pprof by go-torch
- find that addTagsToName causes the problems 

- writer benchmark for addTagsToName func
```
go test -bench=. -benchmem -cpuprofile=cpu.prof // in that way we will find that problem might be caused by too many allocations
	top
	list addTagsToName
	list clean
	
	remove regex and replace it by:
	func clean(value string) string {
    	newStr := make([]byte, len(value))
    	for i := 0; i < len(value) ; i++ {
    		switch c := value[i]; c {
    		case '[', '{', '}', '/', '\\', ':', ' ', '\t', '.':
    			newStr[i] = '-'
    		default:
    			newStr[i] = c
    		}
    	}
    
    	return string(newStr)
    }
	
	run benchmark once again
	show growslice or makeslice
	disasm clean -> show CALL runtime.slicebytetostring(SB)
	run list addTagsToName
	
	show that problem is with:  var keyOrder []string and parts := []string{name}
	change it to make([]string, 0, 4)
	and make([]string, 0, 5) parts[0] = name
	 
	run benchmark one with cpu2.prof as output
	show benchcmp cpu.prof cpu2.prof
	 
	run pprof once again and run 'list addTagsToName'
	show that slice of strings in not needed - the only thing which we really need is one single string
	change slice of string to bytes.Buffer
	
	change func 'clear' to get bytes.Buffer as param and write to that buffer
	run benchmark once again and show that there are 2 allocs and time ~ 400ns/op
	
	can we get rid of those 2 allocations?
	go test -bench=. -benchmem -memprofile=mem.prof
	go tool pprof -alloc_objects stats.test mem.prof
	
	show that there are only 2 allocations: 'buf := bytes.Buffer{}' and 'return buf.String()' and we can remove the first one
	describe what syncpool is
	create sync pool of buffers
	
	var bufPool = sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}
	
	and: 
	
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufPool.Put(buf)

	run benchmark once again
	show that the time is bigger - why is that? defer
	
	run cpuprofile once again
		top
		list addTagsToName
		show cost of the 'defer'
		
	run go-torch stats.test cpu.prof
	show the cost of defer
	
	end of the func: + remove the defer

		final := buf.String()
    	bufPool.Put(buf)
    	return final
	
	
	run benchmark once again
	
	check final performance of the whole http server
	