package writesplitter

import (
	"bufio"
	"bytes"
	"errors"
	"testing"
	"io"
)

type mockF struct {
	*bytes.Buffer
}

func (m *mockF) Close() error {
	m.Reset()
	return nil
}

func TestWriteNoSplit(t *testing.T) {
	createFile = func(name string) (io.WriteCloser, error) {
		return &mockF{&bytes.Buffer{}}, nil
	}
	ws := LineSplitter(5, "")

	mockD := bytes.NewBufferString(`Lorem ipsum dolor sit amet consectetur adipiscing elit
Cras in lacinia eros Aliquam aliquet sapien a
Ut mauris orci varius et cursus sed blandit
Mauris iaculis ac magna non tincidunt In rhoncus
Pellentesque quis erat quis ex aliquam porttitor Vestibulum
Pellentesque nec mollis nibh interdum eleifend nisl Donec
id commodo urna sed tempus mi Vestibulum facilisis
imperdiet dolor sed sollicitudin Proin in lectus sed`)

	expected := mockD.Len() - 7 // we do *not* expect the newlines
	total := 0

	scanner := bufio.NewScanner(mockD)
	for scanner.Scan() {

		n, _ := ws.Write(scanner.Bytes())
		total += n

		if err := scanner.Err(); err != nil {
			t.Error("scanner error", err)
		}
	}

	if expected != total {
		t.Error("len() mismatch: expected", expected, "actual", total)
	}

}

func TestWriteSplit(t *testing.T) {
	var b bytes.Buffer // pass in the buffer to allow for inspection

	createFile = func(name string) (io.WriteCloser, error) {
		return &mockF{&b}, nil
	}

	ws := ByteSplitter(255, "")

	mockD := bytes.NewBufferString(`Lorem ipsum dolor sit amet consectetur adipiscing elit
Cras in lacinia eros Aliquam aliquet sapien a
Ut mauris orci varius et cursus sed blandit
Mauris iaculis ac magna non tincidunt In rhoncus
Pellentesque quis erat quis ex aliquam porttitor Vestibulum
Pellentesque nec mollis nibh interdum eleifend nisl Donec
id commodo urna sed tempus mi Vestibulum facilisis
imperdiet dolor sed sollicitudin Proin in lectus sed`)

	expected := 102 // only the last three lines less two newlines (\n)
	total := 0

	scanner := bufio.NewScanner(mockD)
	for scanner.Scan() {

		n, _ := ws.Write(scanner.Bytes())
		total += n

		if err := scanner.Err(); err != nil {
			t.Error("scanner error", err)
		}
	}

	if expected != b.Len() {
		t.Error("len() mismatch: expected", expected, "actual", b.Len())
	}

}

func TestErrorOnCreate(t *testing.T) {
	createFile = func(name string) (io.WriteCloser, error) {
		return nil, errors.New("This is an error")
	}
	ws := LineSplitter(5, "")

	mockD := bytes.NewBufferString(`Lorem ipsum dolor sit amet consectetur adipiscing elit
Cras in lacinia eros Aliquam aliquet sapien a
Ut mauris orci varius et cursus sed blandit
Mauris iaculis ac magna non tincidunt In rhoncus
Pellentesque quis erat quis ex aliquam porttitor Vestibulum
Pellentesque nec mollis nibh interdum eleifend nisl Donec
id commodo urna sed tempus mi Vestibulum facilisis
imperdiet dolor sed sollicitudin Proin in lectus sed`)

	expected := mockD.Len() - 7 // we do *not* expect the newlines
	total := 0

	var n int
	var err error

	scanner := bufio.NewScanner(mockD)
	for scanner.Scan() {

		n, err = ws.Write(scanner.Bytes())
		if err == nil {
			t.Fail()
		}
		total += n

		if err = scanner.Err(); err != nil {
			t.Error("scanner error", err)
		}
	}

	if expected == total {
		t.Error("len() mismatch: expected", expected, "actual", total)
	}

}
