package fwc

import (
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/timdrysdale/hub"
)

func TestInstantiateHub(t *testing.T) {

	mh := hub.New()

	h := New(mh)

	if reflect.TypeOf(h.Clients) != reflect.TypeOf(make(map[string]*Client)) {
		t.Error("Hub.Clients map of wrong type")
	}
	if reflect.TypeOf(h.Add) != reflect.TypeOf(make(chan Rule)) {
		t.Error("Hub.Add channel of wrong type")
	}

	if reflect.TypeOf(h.Delete) != reflect.TypeOf(make(chan string)) {
		t.Errorf("Hub.Delete channel of wrong type wanted/got %v %v", reflect.TypeOf(""), reflect.TypeOf(h.Delete))
	}

	if reflect.TypeOf(h.Rules) != reflect.TypeOf(make(map[string]Rule)) {
		t.Error("Hub.Broadcast channel of wrong type")
	}

}

func TestAddRule(t *testing.T) {

	mh := hub.New()
	h := New(mh)
	closed := make(chan struct{})
	go h.Run(closed)

	id := "rule0"
	stream := "/stream/large"
	filename := "./large.ts"

	r := &Rule{Id: id,
		Stream:   stream,
		Filename: filename}

	h.Add <- *r

	time.Sleep(time.Millisecond)

	if _, ok := h.Rules[id]; !ok {
		t.Error("Rule not registered in Rules")

	} else {

		if h.Rules[id].Filename != filename {
			t.Errorf("Rule has incorrect filename wanted/got %v %v\n", filename, h.Rules[id].Filename)
		}
		if h.Rules[id].Stream != stream {
			t.Errorf("Rule has incorrect stream wanted/got %v %v\n", stream, h.Rules[id].Stream)
		}
	}
	close(closed)
}

func TestAddRules(t *testing.T) {

	closed := make(chan struct{})

	mh := hub.New()
	go mh.Run(closed)

	h := New(mh)
	go h.Run(closed)

	id := "rule0"
	stream := "/stream/large"
	filename := "./large.ts"

	r := &Rule{Id: id,
		Stream:   stream,
		Filename: filename}

	h.Add <- *r

	id2 := "rule1"
	stream2 := "/stream/large2"
	filename2 := "./large2.ts"

	r2 := &Rule{Id: id2,
		Stream:   stream2,
		Filename: filename2}

	h.Add <- *r2

	time.Sleep(time.Millisecond)

	if _, ok := h.Rules[id]; !ok {
		t.Error("Rule not registered in Rules")

	} else {

		if h.Rules[id].Filename != filename {
			t.Errorf("Rule has incorrect filename wanted/got %v %v\n", filename, h.Rules[id].Filename)
		}
		if h.Rules[id].Stream != stream {
			t.Errorf("Rule has incorrect stream wanted/got %v %v\n", stream, h.Rules[id].Stream)
		}
	}

	if _, ok := h.Rules[id2]; !ok {
		t.Error("Rule not registered in Rules")

	} else {

		if h.Rules[id2].Filename != filename2 {
			t.Errorf("Rule has incorrect filename wanted/got %v %v\n", filename2, h.Rules[id2].Filename)
		}
		if h.Rules[id2].Stream != stream2 {
			t.Errorf("Rule has incorrect stream wanted/got %v %v\n", stream2, h.Rules[id2].Stream)
		}
	}
	err := os.Remove(filename2)
	if err != nil {
		t.Errorf("Error deleting test file %s\n", filename2)
	}
	close(closed)
}

func TestAddDupeRule(t *testing.T) {

	closed := make(chan struct{})

	mh := hub.New()
	go mh.Run(closed)

	h := New(mh)
	go h.Run(closed)

	id := "rule0"
	stream := "/stream/large"
	filename := "./large.ts"

	r := &Rule{Id: id,
		Stream:   stream,
		Filename: filename}

	h.Add <- *r

	//no change to id
	stream2 := "/stream/large2"
	filename2 := "./large2.ts"

	r2 := &Rule{Id: id, //same id as first rule
		Stream:   stream2,
		Filename: filename2}

	h.Add <- *r2

	time.Sleep(time.Millisecond)

	if _, ok := h.Rules[id]; !ok {
		t.Error("Rule not registered in Rules")

	} else {

		if h.Rules[id].Filename != filename2 {
			t.Errorf("Rule has incorrect filename wanted/got %v %v\n", filename2, h.Rules[id].Filename)
		}
		if h.Rules[id].Stream != stream2 {
			t.Errorf("Rule has incorrect stream wanted/got %v %v\n", stream2, h.Rules[id].Stream)
		}
	}
	err := os.Remove(filename2)
	if err != nil {
		t.Errorf("Error deleting test file %s\n", filename2)
	}
	close(closed)
}

func TestDeleteRule(t *testing.T) {

	closed := make(chan struct{})

	mh := hub.New()
	go mh.Run(closed)

	h := New(mh)
	go h.Run(closed)

	id := "rule0"
	stream := "/stream/large"
	filename := "./large.ts"

	r := &Rule{Id: id,
		Stream:   stream,
		Filename: filename}

	h.Add <- *r

	id2 := "rule1"
	stream2 := "/stream/large2"
	filename2 := "./large2.ts"

	r2 := &Rule{Id: id2,
		Stream:   stream2,
		Filename: filename2}

	h.Add <- *r2

	time.Sleep(time.Millisecond)

	if _, ok := h.Rules[id]; !ok {
		t.Error("Rule not registered in Rules")

	} else {

		if h.Rules[id].Filename != filename {
			t.Errorf("Rule has incorrect filename wanted/got %v %v\n", filename, h.Rules[id].Filename)
		}
		if h.Rules[id].Stream != stream {
			t.Errorf("Rule has incorrect stream wanted/got %v %v\n", stream, h.Rules[id].Stream)
		}
	}

	if _, ok := h.Rules[id2]; !ok {
		t.Error("Rule not registered in Rules")

	} else {

		if h.Rules[id2].Filename != filename2 {
			t.Errorf("Rule has incorrect filename wanted/got %v %v\n", filename2, h.Rules[id2].Filename)
		}
		if h.Rules[id2].Stream != stream2 {
			t.Errorf("Rule has incorrect stream wanted/got %v %v\n", stream2, h.Rules[id2].Stream)
		}
	}

	h.Delete <- r.Id

	time.Sleep(time.Millisecond)

	if _, ok := h.Rules[id]; ok {
		t.Error("Deleted rule registered in Rules")

	}

	err := os.Remove(filename2)
	if err != nil {
		t.Errorf("Error deleting test file %s\n", filename2)
	}

	close(closed)
}

func TestWriteMessage(t *testing.T) {

	closed := make(chan struct{})

	mh := hub.New()
	go mh.Run(closed)

	h := New(mh)
	go h.Run(closed)

	id := "rule0"
	stream := "/stream/large"
	filename := "./large.ts"

	r := &Rule{Id: id,
		Stream:   stream,
		Filename: filename}

	h.Add <- *r

	reply := make(chan hub.Message)

	c := &hub.Client{Hub: mh, Name: "testing", Topic: stream, Send: reply}

	mh.Register <- c

	time.Sleep(time.Millisecond)

	lines := []string{"a\n", "test\n", "message\n"}

	for _, line := range lines {

		payload := []byte(line)

		mh.Broadcast <- hub.Message{Data: payload, Type: websocket.TextMessage, Sender: *c, Sent: time.Now()}

		time.Sleep(time.Millisecond)

	}

	time.Sleep(time.Millisecond)

	h.Delete <- r.Id

	time.Sleep(time.Millisecond)

	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Error("Error reading file", err)
	}

	tokens := strings.Fields(string(dat))

	if len(lines) != len(tokens) {
		t.Errorf("Incorrect number of lines in file; got %d, wanted %d\n", len(tokens), len(lines))
	} else {

		for i, line := range lines {
			if strings.TrimSpace(line) != strings.TrimSpace(tokens[i]) {
				t.Errorf("lines in file did not match messages %s/%s", line, tokens[i])
			}
		}
	}

	err = os.Remove(filename)
	if err != nil {
		t.Errorf("Error deleting test file %s\n", filename)
	}

}

func TestWriteMessageToChangingDestination(t *testing.T) {

	closed := make(chan struct{})

	mh := hub.New()
	go mh.Run(closed)

	h := New(mh)
	go h.Run(closed)

	id := "rule0"
	stream := "/stream/large"
	filename := "./large.ts"

	r := &Rule{Id: id,
		Stream:   stream,
		Filename: filename}

	h.Add <- *r

	reply := make(chan hub.Message)

	c := &hub.Client{Hub: mh, Name: "testing", Topic: stream, Send: reply}

	mh.Register <- c

	time.Sleep(time.Millisecond)

	lines := []string{"test\n", "message\n"}

	payload := []byte(lines[0])

	mh.Broadcast <- hub.Message{Data: payload, Type: websocket.TextMessage, Sender: *c, Sent: time.Now()}

	time.Sleep(time.Millisecond)

	filename2 := "./large2.ts"

	r = &Rule{Id: id,
		Stream:   stream,
		Filename: filename2}

	h.Add <- *r

	time.Sleep(time.Millisecond)

	payload = []byte(lines[1])

	mh.Broadcast <- hub.Message{Data: payload, Type: websocket.TextMessage, Sender: *c, Sent: time.Now()}

	time.Sleep(time.Millisecond)

	h.Delete <- r.Id

	time.Sleep(time.Millisecond)

	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Error("Error reading file", err)
	}

	tokens := strings.Fields(string(dat))

	if len(tokens) != 1 {
		t.Errorf("Incorrect number of lines in file; got %d, wanted %d\n", len(tokens), 1)
	} else {
		if strings.TrimSpace(lines[0]) != strings.TrimSpace(tokens[0]) {
			t.Errorf("lines in file did not match messages %s/%s", lines[0], tokens[0])
		}
	}

	dat, err = ioutil.ReadFile(filename2)
	if err != nil {
		t.Error("Error reading file", err)
	}

	tokens = strings.Fields(string(dat))

	if len(tokens) != 1 {
		t.Errorf("Incorrect number of lines in file; got %d, wanted %d\n", len(tokens), 1)
	} else {
		if strings.TrimSpace(lines[1]) != strings.TrimSpace(tokens[0]) {
			t.Errorf("lines in file did not match messages %s/%s", lines[1], tokens[0])
		}
	}

	err = os.Remove(filename)
	if err != nil {
		t.Errorf("Error deleting test file %s\n", filename)
	}

	err = os.Remove(filename2)
	if err != nil {
		t.Errorf("Error deleting test file %s\n", filename)
	}
}
