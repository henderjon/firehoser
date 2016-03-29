## firehoser

Similar to AWS' Firehose, this application will run listening for either
HTTP or TCP connections (not simultaneously), scan the incoming data line by line and write it to
STDOUT. The intended functionality was to quickly ingest line based (CSV/TSV)
log data. Instead of writing to a file directly, writing to STDOUT allows the
user to pipe the output to other programs. For instance, redirecting the output
to `split` will write canonically named text files of a given size (e.g.
`./firehoser | split -a 9 -l 100 - mylogfile`).

## todo

  - graceful shutdown (catch broken pipes)
  - ~~command line args to alter behavior~~
  - errors for non-POST methods/more verbose responses
