## log all the things

[![GoDoc](https://godoc.org/github.com/henderjon/omnilogger?status.svg)](https://godoc.org/github.com/henderjon/omnilogger)

Omnilogger is an HTTP (or TCP) server that coalesces log data (line by line) from
multiple sources to a common destination (defaults to consecutively named log
files of ~5000 lines).

Use `-h` to view the available options

## faq

### Why?

The intended functionality was to quickly ingest line based (CSV/TSV)
log data from many different EC2 instances being auto-scaled.

### Won't I fill up my mom's 250GB hard drive really fast?

Potentially, yes. I'd recommend a cron job that rotates logs to a long-term
storage facility--something like AWS S3.

### Is there anything else I ought to know?

  - The default HTTP server is better than the TCP server.
  - If you use the TCP server, idle connections are closed after 3 seconds.
  - The sub-library `writesplitter` is useful on it's own, and has it's own [README.md](writesplitter/README.md).
  - All HTTP requests must be a POST, but the body is not parsed (e.g. form-encoded data will get logged as is)
  - All HTTP requests have to send a custom header ('X-Omnilog-Stream'). As of now
    it only checks to see if it's there. In the future it **might** use this header
    to divert data to separate destinations.


