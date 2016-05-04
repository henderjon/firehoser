## log all the things

[![License: BSD-3](https://img.shields.io/badge/license-BSD--3-blue.svg)](https://img.shields.io/badge/license-BSD--3-blue.svg)
[![GoDoc](https://godoc.org/github.com/henderjon/omnilogger?status.svg)](https://godoc.org/github.com/henderjon/omnilogger)
[![Build Status](https://travis-ci.org/henderjon/omnilogger.svg?branch=dev)](https://travis-ci.org/henderjon/omnilogger)
[![Go Report Card](https://goreportcard.com/badge/github.com/henderjon/omnilogger)](https://goreportcard.com/report/github.com/henderjon/omnilogger)

Omnilogger is an HTTP server that ingests log data from multiple sources to a
common destination. Each worker (default 2) has a buffer in memory (default 64k).
When a buffer is filled, it's written to disk. Given the number of cores on the
machine you're using, you'll need to play with the number and size of the workers.
There is also a buffer (default 500) for incoming requests that feeds all four workers.

Use `-h` to view the available options

## faq

### Why?

The intended functionality was to quickly ingest line based (CSV/TSV)
log data from many different EC2 instances being auto-scaled up and down.

### Won't I fill up my mom's 250GB hard drive really fast?

Potentially, yes. I'd recommend a cron job that rotates logs to a long-term
storage facility--something like AWS S3.

### Is there anything else I ought to know?

  - All HTTP requests must be a POST, but the body is not parsed (e.g.
    form-encoded data will get logged as is)
  - All HTTP requests have to send a custom header ('X-Omnilogger-Stream'). As of
    now it only checks to see if it's there. In the future it **might** use
    this header to divert data to separate destinations or other purposes.
  - After 10 minutes of inactivity, the currently open file is closed. Another is
    automatically opened on the next write.

### You're HTTP errors are kinda cryptic.

Yeah, I'm lazy. Instead of making super cool error messages, I'm depending on
the default text related to a given [HTTP status code](https://golang.org/pkg/net/http/#pkg-constants).

  - 200 (Ok) means that everything should have worked just fine. If not,
    report an issue.
  - 400 (Bad Request) means you didn't send the 'X-Omnilogger-Stream' header even
    after I told you to.
  - 403 (Forbidden) means (as of now) the token you sent in the Authorization
    header doesn't match what the server is looking for. I'm using the
    `Authorization: Bearer $token` style header like all the cool kids.
  - 405 (Method Not Allowed) means your HTTP method was something other than
    POST (*tsk tsk*).
  - 503 (Service Unavailable) means the system is shutting down.

### Why would you ever want ALL your various log data in one stream?

The goal was to collect data from a variable number of servers as quickly as
possible. To this end, by convention, all the data is sent as interlaced csv
rows--the last value of each line is the name of the stream. Part of being
*quick*, is to do as little as possible with the data. Therefore, *at this
time*, there isn't a need for that feature because a simple one-liner in AWK
will do this after the fact when speed and time are less of an issue (e.g.
`cat file.log | awk -c '$(NR) == "stream_name"{print}'`). If keeping streams
separate is important, multiple instances running on different ports can be
used to accomplish the same thing.


## todo

  - stream splitting (based on header)
