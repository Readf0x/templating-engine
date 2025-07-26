package main

import (
	"fmt"
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

const help = `The best templating engine

Usage: te [FILE]...

Arguments:
  [FILE]...  File(s) to process. Must end with '.tet'

Options:
  -h, --help
	Show this text`

func main() {
	if len(os.Args) == 1 {
		fmt.Println("No arguments")
		os.Exit(2)
	} else {
		offset := 1
		switch os.Args[1] {
		case "--help":
			fallthrough
		case "-h":
			fmt.Println(help)
			offset++
		}
		for _, A := range os.Args[offset:] {
			if strings.HasSuffix(A, ".tet") {
				outname := A[:len(A)-4]
				fmt.Println(outname)
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
