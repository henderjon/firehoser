[![GoDoc](https://godoc.org/github.com/henderjon/omnilogger/writesplitter?status.svg)](https://godoc.org/github.com/henderjon/omnilogger/writesplitter)

WriteSplitter represents a disk bound io.WriteCloser that splits the input
across multiple files based on either the number of bytes or the number of
lines. Splitting does not guarantee true byte/line split precision as it does
not parse the incoming data. The decision to split is before the underlying
write operation based on the previous invocation. In other words, if a []byte
sent to `Write()` contains enough bytes or new lines ('\n') to exceed the
given limit, a new file won't be generated until the *next* invocation of
`Write()`. If both LineLimit and ByteLimit is set, preference is given to
LineLimit. By default, no splitting occurs because both LineLimit and
ByteLimit are zero (0).


