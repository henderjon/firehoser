## log all the things

[![GoDoc](https://godoc.org/github.com/henderjon/omnilogger?status.svg)](https://godoc.org/github.com/henderjon/omnilogger)
[![Build Status](https://travis-ci.org/henderjon/omnilogger.svg?branch=master)](https://travis-ci.org/henderjon/omnilogger)

Omnilogger is an HTTP server that coalesces log data (line by line) from
multiple sources to a common destination (defaults to consecutively named log
files of ~5000 lines).

Use `-h` to view the available options

## faq

### Why?

The intended functionality was to quickly ingest line based (CSV/TSV)
log data from many different EC2 instances being auto-scaled up and down.

### Won't I fill up my mom's 250GB hard drive really fast?

Potentially, yes. I'd recommend a cron job that rotates logs to a long-term
storage facility--something like AWS S3.

### Is there anything else I ought to know?

  - The library `writesplitter` is useful on it's own, and has it's own [README.md](writesplitter).
  - All HTTP requests must be a POST, but the body is not parsed (e.g. form-encoded data will get logged as is)
  - All HTTP requests have to send a custom header ('X-Omnilog-Stream'). As of now
    it only checks to see if it's there. In the future it **might** use this header
    to divert data to separate destinations.

### You're HTTP errors are kinda cryptic.

Yeah, I'm lazy. Instead of making super cool error messages, I'm depending on
the default text related to a given [HTTP status code](https://golang.org/pkg/net/http/#pkg-constants).

  - 200 (Ok) means that everything should have worked just fine. If not, report an issue.
  - 400 (Bad Request) means you didn't send the 'X-Omnilog-Stream' header even after I told you to.
  - 403 (Forbidden) means (as of now) the token you sent in the Authorization header doesn't match what the server is
    looking for. I'm using the `Authorization: Bearer $token` style header like all the cool kids.
  - 405 (Method Not Allowed) means your HTTP method was something other than POST (*tsk tsk*).
  - 503 (Service Unavailable) means the system is shutting down.

## todo

  - time expiration (file close/reopen)
  - stream splitting (based on header)
