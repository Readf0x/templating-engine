package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func w(str string, o *os.File) {
	o.Write([]byte(str))
}
func wb(b byte, o *os.File) {
	o.Write([]byte{b})
}

func main() {
	if len(os.Args) == 1 {
		println("No files")
		os.Exit(2)
	} else {
		for _, A := range os.Args[1:] {
			if strings.HasSuffix(A, ".tet") {
				outname := A[:len(A)-4]
				println(outname)
				f, _ := os.Open(A)
				defer f.Close()
				fname := "build-" + strings.ReplaceAll(outname, "/", "_") + ".go"
				out, _ := os.Create(fname)
				b, _ := io.ReadAll(f)
				content := string(b)
				w(`//go:build ignore
package main

import (
	"os"
	"fmt"
)

var out, _ = os.Create("`+outname+`")

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
								i += 2
							}
							i += 2
						case "|>":
							if lft {
								w(suffix, out)
								lft = false
							}
							w(";w(`", out)
							i += 2
						}
					}
					wb(content[i], out)
				}
				w("`);\n}\n", out)

				cmd := exec.Command("go", "run", fname)
				if err := cmd.Run(); err != nil {
					log.Fatal(err)
				}
				out.Close()
				os.Remove(fname)
			}
		}
	}
}

func pft(char byte) (string, string) {
	switch char {
	case 'w':
		return "w(", ");"
	}
	return "", ""
}
