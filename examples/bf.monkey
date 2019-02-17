#!./monkey-lang

fill := fn(x, i) {
  xs := []
  while (i > 0) {
    xs = push(xs, x)
    i = i - 1
  }
  return xs
}

buildJumpMap := fn(program) {
  stack := []
  map := {}

  n := 0
  while (n < len(program)) {
    if (program[n] == "[") {
      stack = push(stack, n)
    }
    if (program[n] == "]") {
      start := pop(stack)
      map[start] = n
      map[n] = start
    }
    n = n + 1
  }

  return map
}

ascii_table := "\x00\x01\x02\x03\x04\x05\x06\x07\x08\t\n\x0b\x0c\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~\x7f\x80\x81\x82\x83\x84\x85\x86\x87\x88\x89\x8a\x8b\x8c\x8d\x8e\x8f\x90\x91\x92\x93\x94\x95\x96\x97\x98\x99\x9a\x9b\x9c\x9d\x9e\x9f\xa0\xa1\xa2\xa3\xa4\xa5\xa6\xa7\xa8\xa9\xaa\xab\xac\xad\xae\xaf\xb0\xb1\xb2\xb3\xb4\xb5\xb6\xb7\xb8\xb9\xba\xbb\xbc\xbd\xbe\xbf\xc0\xc1\xc2\xc3\xc4\xc5\xc6\xc7\xc8\xc9\xca\xcb\xcc\xcd\xce\xcf\xd0\xd1\xd2\xd3\xd4\xd5\xd6\xd7\xd8\xd9\xda\xdb\xdc\xdd\xde\xdf\xe0\xe1\xe2\xe3\xe4\xe5\xe6\xe7\xe8\xe9\xea\xeb\xec\xed\xee\xef\xf0\xf1\xf2\xf3\xf4\xf5\xf6\xf7\xf8\xf9\xfa\xfb\xfc\xfd\xfe\xff"

ord := fn(s) {
  return ascii_table[s]
}

chr := fn(x) {
  if (x < 0) {
    return "??"
  }
  if (x > 127) {
    return "??"
  }
  return ascii_table[x]
}

read := fn() {
  buf := input()
  if (len(buf) > 0) {
    return ord(buf[0])
  }
  return 0
}

write := fn(x) {
  print(chr(x))
}

VM := fn(program) {
  jumps := buildJumpMap(program)

  ip := 0
  dp := 0
  memory := fill(0, 32)

  while (ip < len(program)) {
    op := program[ip]

    if (op == ">") {
      dp = dp + 1
    } else if (op == "<") {
      dp = dp - 1
    } else if (op == "+") {
      if (memory[dp] < 255) {
        memory[dp] = memory[dp] + 1
      } else {
        memory[dp] = 0
      }
    } else if (op == "-") {
      if (memory[dp] > 0) {
        memory[dp] = memory[dp] - 1
      } else {
        memory[dp] = 255
      }
    } else if (op == ".") {
      write(memory[dp])
    } else if (op == ",") {
      memory[dp] = read()
    } else if (op == "[") {
      if (memory[dp] == 0) {
        ip = jumps[ip]
      }
    } else if (op == "]") {
      if (memory[dp] != 0) {
        ip = jumps[ip]
      }
    }
    ip = ip + 1
  }

  print("memory:")
  print(memory)
  print("ip:")
  print(ip)
  print("dp:")
  print(dp)
}

// Hello World
program := "++++++++ [ >++++ [ >++ >+++ >+++ >+ <<<<- ] >+ >+ >- >>+ [<] <- ] >>.  >---.  +++++++..+++.  >>.  <-.  <.  +++.------.--------.  >>+.  >++."

// 2 + 5
// program := "++> +++++ [<+>-] ++++++++ [<++++++>-] < ."

VM(program)
