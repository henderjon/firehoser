package writesplitter

import (
	"testing"
	"bytes"
	"bufio"
)

type mockFs struct {}

type mockF struct{
	bytes.Buffer
}

func (mockFs) Create(name string) (file, error) {
	return &mockF{}, nil
}

func (m *mockF) Close() error {
	m.Reset()
	return nil
}

func TestWriteNoSplit(t *testing.T) {
	fs = mockFs{}
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
	total    := 0

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
