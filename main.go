package main

import (
	"flag"
	"github.com/henderjon/omnilogger/counter"
	"github.com/henderjon/omnilogger/shutdown"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
	"path/filepath"
	"io/ioutil"
)

// the uint64 representation of a Kilobyte, Megabyte, Gigabyte as well as some program defaults
const (
	Kilobyte        = 1024
	Megabyte        = Kilobyte * Kilobyte
	Gigabyte        = Kilobyte * Megabyte
	defaultInterval = 10 * time.Minute
	defaultPrefix   = "log-omnilogs-"
	defaultPerms    = 0644
)

var (
	port          string                      // the port on which to listen
	pswd          string                      // a simple means of authentication
	requestBuffer int                         // the size of the incoming request buffer (channel)
	size          int                         // how many lines per log file
	scale         bool                        // how many bytes per log file
	numWorkers    int                         // how many bytes per log file
	logDir        string                      // the dir for the log file(s)
	help          bool                        // I forgot my options
	wg            sync.WaitGroup              // ensure that our goroutines finish before shut down
	methodPost    = "POST"                    // because net/http doesn't have this ...
	customHeader  = "X-Omnilogger-Stream"     // a custom header to validate intent
	closeInterval = defaultInterval           // how often to close our file and open a new one
	byteCount     = counter.NewCounter()      // keep track of how many bytes total have been received
	hitCount      = counter.NewCounter()      // keep track of how many bytes total have been received
	helpLogger    = log.New(os.Stderr, "", 0) // log to stderr without the timestamps
)

func init() {
	flag.StringVar(&port, "port", "8080", "The port used for the server.")
	flag.StringVar(&pswd, "auth", "", "If not empty, this is matched against the Authorization header (e.g. Authorization: Bearer my-password).")
	flag.IntVar(&requestBuffer, "buf", 500, "The size of the incoming request buffer. A zero (0) will disable buffering.")
	flag.IntVar(&numWorkers, "workers", 2, "The number of workers/buffers... most likely doesn't need changing and should probably not exceed the number of CPUs on the machine.")
	flag.IntVar(&size, "size", 64, "The size (in kilobytes) at which to split the log file(s).")
	flag.BoolVar(&scale, "m", false, "If set, -size will be in megabytes.")
	flag.StringVar(&logDir, "dir", "", "A dir to use for log files. The first arg after '--' is used as a filename prefix. (e.g. '% omnilogger -dir /path/to/log-dir -- file-prefix-')")
	flag.BoolVar(&help, "h", false, "Show this message.")
	flag.Parse()

	if help {
		helpLogger.Println("")
		helpLogger.Println("Omnilogger is an HTTP server that ingests log data from multiple sources to a common destination.")
		helpLogger.Println("")
		flag.PrintDefaults()
		helpLogger.Println("")
		os.Exit(0)
	}

	welcome()

}

func main() {
	inbound := make(chan []byte, requestBuffer)
	signal := make(shutdown.SignalChan)

	// if capacity == 0 -> stdout?
	for t := 0; t < numWorkers; t += 1 {
		go func(){
			for {
				b := <-inbound // pull data out of the channel
				name := filepath.Join(logDir, flag.Arg(0) + time.Now().Format(time.RFC3339Nano))
				ioutil.WriteFile(name, b, defaultPerms)
			}
		}()
	}

	go shutdown.Watch(signal, destructor)

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", http.StripPrefix("/", fs))
	http.Handle("/log", Adapt(parseRequest(inbound, &wg), checkHeader(customHeader), checkAuth(pswd), checkMethod(methodPost), checkShutdown(signal)))
	if e := http.ListenAndServe(":"+port, nil); e != nil {
		log.Fatal(e)
	}
}

// destructor is the func that gets called if we catch a shutdown signal. It waits
// for all goroutines to finish and then prints a final status message
func destructor() {
	wg.Wait()
	helpLogger.Printf(".collected %dm from %d hits in %s", byteCount.Current(Megabyte), hitCount.Current(0), byteCount.Since())
}
