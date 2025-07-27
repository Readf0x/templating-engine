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

const providedFunctions = ``+
//te:start
// Write
// Writes out to finalized file.
`func w(val any) {
	out.Write([]byte(fmt.Sprint(val)))
}`
//te:stop

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
				var imports []string
				if strings.HasPrefix(content, "!!import:") {
					newline := strings.IndexByte(content, '\n') + 1
					firstline := content[:newline]
					imports = strings.Split(firstline[9:], " ")
					content = content[newline:]
				}
				w(`//go:build ignore
package main

import (
	"os"
	"fmt"
	`+strings.Join(imports, "\n  ")+`
)

var out, _ = os.Create("`+outname+`")

`+providedFunctions+`

func main() {
	defer out.Close()
`, out)
				w("w(`", out)
				lft := false
				var prefix, suffix string
				var code bool
				for i := 0; i < len(content); i++ {
					s := content[i:]
					if s[0] == '`' {
						if !code {
							w("`+\"`\"+`", out)
							i++
						}
					}
					if len(s) > 1 {
						switch s[:2] {
						case "\\<":
							if s[2] == '|' {
								w("<|", out)
								i += 3
							}
						case "\\|":
							if s[2] == '>' {
								w("|>", out)
								i += 3
							}
						case "<|":
							w("`);\n", out)
							code = true
							if s[2] == ':' {
								lft = true
								prefix, suffix = pft(s[3])
								w(prefix, out)
								i += 2
							}
							i += 2
						case "|>":
							code = false
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
				cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
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
	return fmt.Sprintf("%c(", char), ");"
}
