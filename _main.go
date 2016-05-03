package buffered

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	workerCount = 4
	capacity    = 32768
)

var (
	channel = make(chan []byte, 1024)
	workers = make([]*Worker, workerCount)
)

type Worker struct {
	fileRoot string
	buffer   []byte
	position int
}

func NewWorker(id int) (w *Worker) {
	return &Worker{
		//move the root path to some config or something
		fileRoot: strconv.Itoa(id) + "_",
		buffer:   make([]byte, capacity),
	}
}

func init() {
	for i := 0; i < workerCount; i++ {
		workers[i] = NewWorker(i)
		go workers[i].Work(channel)
	}
}

func Log(event []byte) {
	select {
	case channel <- event:
	case <-time.After(5 * time.Second):
		// throw away the message, so sad
	}
}

func (w *Worker) Work(channel chan []byte) {
	for {
		event := <-channel
		length := len(event)
		// we run with nginx's client_max_body_size set to 2K which makes this
		// unlikely to happen, but, just in case...
		if length > capacity {
			log.Println("message received was too large")
			continue
		}
		if (length + w.position) > capacity {
			w.Save()
		}
		copy(w.buffer[w.position:], event)
		w.position += length
	}
}

func (w *Worker) Save() {
	if w.position == 0 {
		return
	}
	f, _ := ioutil.TempFile("", "logs_")
	f.Write(w.buffer[0:w.position])
	f.Close()
	os.Rename(f.Name(), w.fileRoot+strconv.FormatInt(time.Now().UnixNano(), 10))
	w.position = 0
}
