package main

import (
	"io"
	"os"
)

func w(str string, o *os.File) {
	o.Write([]byte(str))
}
func wb(b byte, o *os.File) {
	o.Write([]byte{b})
}

func main() {
	// name := os.Args[1]
	name := "index.html.tet"
	f, _ := os.Open(name)
	defer f.Close()
	out, _ := os.Create("build.go")
	defer out.Close()
	b, _ := io.ReadAll(f)
	content := string(b)
	w(`//go:build ignore
package main

import (
	"os"
	"fmt"
)

var out, _ = os.Create("`+name[:len(name)-4]+`")

func w(str string) {
	out.Write([]byte(fmt.Sprint(str)))
}

func main() {
	defer out.Close()
`, out)
	w("w(`", out)
	lft := false
	var prefix, suffix string
	for i := 0; i < len(content); i++ {
		s := content[i:]
		if len(s) > 1 {
			switch s[:2] {
			case "<|":
				w("`);\n", out)
				if s[2] == ':' {
					lft = true
					prefix, suffix = pft(s[3])
					w(prefix, out)
					i+=2
				}
				i+=2
			case "|>":
				if lft {
					w(suffix, out)
					lft = false
				}
				w(";w(`", out)
				i+=2
			}
			wb(content[i], out)
		}
	}
	w("`);\n}\n", out)
}

func pft(char byte) (string, string) {
	switch char {
	case 'w':
		return "w(", ");"
	}
	return "", ""
}

