## log all the things

Similar to AWS' Firehose, this application will run listening for either
HTTP or TCP connections (not simultaneously), scan the incoming data line by line and write it to
STDOUT. The intended functionality was to quickly ingest line based (CSV/TSV)
log data from many different EC2 instances being auto-scaled. Instead of writing to a file directly, writing to STDOUT allows the
user to pipe the output to other programs. For instance, redirecting the output
to `split` will write canonically named text files of a given size (e.g.
`./omnilogger | split -a 9 -l 100 - mylogfile`). To handle many streams, you could run multiple instances of this application, each on a different port or use one instance to coallesce them and parse them after they're on disk.

## todo

  - graceful shutdown (catch broken pipes)
  - how does stdout handle broken pipes
  - ~~command line args to alter behavior~~
  - ~~errors for non-POST methods~~/more verbose responses
  - ~~how does split handle filename collisions~~
  - ~~where should gzip take place~~
