[![GoDoc](https://godoc.org/github.com/k3a/go-m3u?status.svg)](https://godoc.org/github.com/k3a/go-m3u)

# go-m3u

Simple library to read and write M3U records.
Supports EXTINF attributes.

## Usage

```
	m3u := new(M3U)

  // open file or other M3U stream
	file, err := os.Open("./test/working.m3u")
	if err != nil {
		panic(err)
	}

  // parse M3U input
	if err := m3u.Read(file); err != nil {
		panic(err)
	}

  // debug print
	fmt.Println(m3u.String())

  // loop stored m3u records
  for _, record := range m3u.Records() {
    // ...
  }

  // write M3U to output
	b := new(strings.Builder)
	if err := m3u.Write(b); err != nil {
		t.Fatal(err)
	}

  // print the written data
  fmt.Println(m3u.String())
```

## License

MIT
