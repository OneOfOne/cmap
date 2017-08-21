package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"go/scanner"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// if you fork this package, and it's not under OneOfOne, make sure you change this path
const CMapBase = `github.com/OneOfOne/cmap/internal/cmap`

var (
	pkgName   = flag.String("n", "", "package name")
	keyType   = flag.String("kt", "interface{}", "key type")
	valueType = flag.String("vt", "interface{}", "value type")
	hashFn    = flag.String("hfn", "cmap.DefaultKeyHasher", "hash func")
	pkgPath   = flag.String("p", "", "package path")
	lmapOnly  = flag.Bool("lmap", false, "generate LMap only")

	toStdout = flag.Bool("stdout", false, "write the output to stdout rather than creating a package")

	verbose        = flag.Bool("v", false, "verbose")
	debugReplacers = flag.Bool("dr", false, "debug replacers")
	isInternal     = flag.Bool("internal", false, "used to generate the main interface{}/interface{} package")

	cmapBase = execString("go", "list", "-f", "{{ .Dir }}", CMapBase)

	validPackageName = regexp.MustCompile(`^[\w\d]+$`)
)

func main() {
	log.SetFlags(log.Lshortfile)
	flag.Parse()

	if cmapBase == "" {
		log.Fatalf("couldn't find the base package: %s", CMapBase)
	}

	if *pkgName == "" {
		if *pkgPath != "" {
			*pkgName = filepath.Base(*pkgPath)
		} else {
			wd, _ := os.Getwd()
			*pkgName = filepath.Base(wd)

		}
	}

	if !validPackageName.MatchString(*pkgName) {
		log.Fatalf("invalid package name: %s", *pkgName)
	}

	if !*toStdout {
		if *pkgPath == "" {
			*pkgPath = "./" + *pkgName
		}

		if err := os.MkdirAll(*pkgPath, 0755); err != nil {
			log.Fatalf("error making package dir: %s", *pkgPath)
		}
	}

	files := []string{filepath.Join(cmapBase, "cmap.go"), filepath.Join(cmapBase, "lmap.go")}
	if *lmapOnly {
		files = files[1:]
	}

	r, closeReaders := multiFileReader(files...)
	var (
		filters = getLineFilters()
		buf     bytes.Buffer
	)

	fmt.Fprintf(&buf, `// AUTO-GENERATED by cmap-gen
// DO NOT EDIT"
// generated from %s

package %s

`, CMapBase, *pkgName)

	for sc := bufio.NewScanner(r); sc.Scan(); {
		if line := filters.apply(sc.Bytes()); len(line) > 0 {
			buf.Write(line)
			buf.WriteByte('\n')
		}
	}
	closeReaders()

	cmd := exec.Command("goimports", "-e")
	cmd.Stdin = bytes.NewReader(buf.Bytes())
	cmd.Stderr = os.Stderr

	src, err := cmd.Output()
	if err != nil {
		log.Fatalf("goimports error: %v\n%s", err, buf.Bytes())
	}

	if errs := verifyOutput(src); len(errs) > 0 {
		log.Fatalf("error compiling code: \n\t%s\n\n%s", strings.Join(errs, "\n\t"), src)
	}

	var fp string
	if *lmapOnly {
		fp = filepath.Join(*pkgPath, "lmap.go")
	} else {
		fp = filepath.Join(*pkgPath, "cmap.go")
	}

	if *verbose {
		log.Printf("writing typed cmap %s.CMap[%s][%s] (using HashFn: %s) to file %s...", *pkgName, *keyType, *valueType, *hashFn, fp)
	}

	if *toStdout {
		os.Stdout.Write(src)
	} else {
		if err = ioutil.WriteFile(fp, src, 0644); err != nil {
			log.Fatalf("error writing file (%s): %v", fp, err)
		}
	}
}

func verifyOutput(src []byte) (errs []string) {
	var s scanner.Scanner
	fset := token.NewFileSet()                      // positions are relative to fset
	file := fset.AddFile("", fset.Base(), len(src)) // register input "file"
	s.Init(file, src, func(pos token.Position, msg string) {
		errs = append(errs, msg)
	}, 0)

	for _, tok, _ := s.Scan(); tok != token.EOF; _, tok, _ = s.Scan() {
	}
	return
}

type filters []func([]byte) []byte

var (
	bPkg     = []byte("package")
	bImports = []byte("import")
)

func (f filters) apply(line []byte) (_ []byte) {
	if line = bytes.TrimSpace(line); len(line) == 0 {
		return
	}

	if bytes.HasPrefix(line, bPkg) {
		return
	}

	if bytes.HasPrefix(line, bImports) && bytes.ContainsRune(line, '"') {
		return
	}

	for _, fn := range f {
		if line = fn(line); len(line) == 0 {
			break
		}
	}

	return line
}

func getLineFilters() filters {
	out := filters{
		replacer(`KT`, *keyType),
		replacer(`VT`, *valueType),
	}

	if *hashFn != "cmap.DefaultKeyHasher" {
		out = append(out,
			replacer(`cm\.HashFn`, *hashFn),
			replacer(`.*HashFn.*|.*DefaultKeyHasher.*`, ""),
		)
	}

	if *isInternal {
		out = append(out,
			replacer(`cmap\.`, ""),
			replacer(`^(?:var||//) Break.*$`, ""),
		)
	}

	return out
}

func replacer(sre, r string) func([]byte) []byte {
	re := regexp.MustCompile(sre)
	bR := []byte(r)
	return func(in []byte) []byte {
		if *debugReplacers {
			log.Printf("re=%q in=%q repl=%q m=%q", re.String(), in, bR, re.ReplaceAll(in, bR))
		}
		return re.ReplaceAll(in, bR)
	}
}

func execString(name string, args ...string) string {
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		log.Printf("exec error: %v", err)
		return ""
	}
	return string(bytes.TrimSpace(out))
}

func multiFileReader(fnames ...string) (io.Reader, func()) {
	files := make([]io.Reader, 0, len(fnames))
	closer := func() {
		for _, f := range files {
			f.(io.Closer).Close()
		}
	}

	for _, fname := range fnames {
		f, err := os.Open(fname)
		if err != nil {
			closer()
			log.Fatalf("couldn't open %q: %v", fname, err)
		}
		files = append(files, f)
	}

	return io.MultiReader(files...), closer
}
