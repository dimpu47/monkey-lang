package code

import (
	"testing"
)

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		{LoadConstant, []int{65534}, []byte{byte(LoadConstant), 255, 254}},
		{Add, []int{}, []byte{byte(Add)}},
		{LoadLocal, []int{255}, []byte{byte(LoadLocal), 255}},
		{MakeClosure, []int{65534, 255}, []byte{byte(MakeClosure), 255, 254, 255}},
	}

	for _, tt := range tests {
		instruction := Make(tt.op, tt.operands...)

		if len(instruction) != len(tt.expected) {
			t.Errorf("instruction has wrong length. want=%d, got=%d",
				len(tt.expected), len(instruction))
		}

		for i, b := range tt.expected {
			if instruction[i] != tt.expected[i] {
				t.Errorf("wrong byte at pos %d. want=%d, got=%d",
					i, b, instruction[i])
			}
		}
	}
}

func TestInstructionsString(t *testing.T) {
	instructions := []Instructions{
		Make(Add),
		Make(LoadLocal, 1),
		Make(LoadConstant, 2),
		Make(LoadConstant, 65535),
		Make(MakeClosure, 65535, 255),
	}

	expected := `0000 Add
0001 LoadLocal 1
0003 LoadConstant 2
0006 LoadConstant 65535
0009 MakeClosure 65535 255
`
	concatted := Instructions{}
	for _, ins := range instructions {
		concatted = append(concatted, ins...)
	}

	if concatted.String() != expected {
		t.Errorf("instructions wrongly formatted.\nwant=%q\ngot=%q",
			expected, concatted.String())
	}
}

func TestReadOperands(t *testing.T) {
	tests := []struct {
		op        Opcode
		operands  []int
		bytesRead int
	}{
		{LoadConstant, []int{65535}, 2},
		{LoadLocal, []int{255}, 1},
		{MakeClosure, []int{65535, 255}, 3},
	}

	for _, tt := range tests {
		instruction := Make(tt.op, tt.operands...)

		def, err := Lookup(byte(tt.op))
		if err != nil {
			t.Fatalf("definition not found: %q\n", err)
		}

		operandsRead, n := ReadOperands(def, instruction[1:])
		if n != tt.bytesRead {
			t.Fatalf("n wrong. want=%d, got=%d", tt.bytesRead, n)
		}

		for i, want := range tt.operands {
			if operandsRead[i] != want {
				t.Errorf("operand wrong. want=%d, got=%d", want, operandsRead[i])
			}
		}
	}
}
