## log all the things

Omnilogger is an HTTP or TCP server that coalesces log data (line by line) from
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


