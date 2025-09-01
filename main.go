package main

import (
	_ "embed"
	"flag"
	"log"
	"strings"

	"github.com/connordoman/pos/internal/escpos"
)

//go:embed templates/mr_worldwide.txt
var mrWorldwideLyrics string

func main() {
	testFlag := flag.Bool("test", false, "Run test print")
	mrWorldwideFlag := flag.Bool("mrworldwide", false, "Print Mr. Worldwide lyrics")

	flag.Parse()

	p, err := escpos.InitPrinter()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := p.Close(); err != nil {
			log.Printf("close error: %v", err)
		}
	}()

	p.Init()

	if *mrWorldwideFlag {
		p.WriteString(mrWorldwideLyrics)
		return
	}

	if *testFlag {
		// p.TestPrint()

		// if _, err := p.Flush(); err != nil {
		// 	log.Printf("flush error: %v", err)
		// }

		// time.Sleep(1 * time.Second)

		// p.Feed(10)
		// p.Init()
		// p.WriteString("Hello, world!\n")
		p.SelectFont(escpos.FontBName)
		p.Log("Hello, world!")

		p.SimpleLine()

		longLine := strings.Repeat("long-line-", 30)
		p.WriteString(longLine + "\n")

		p.SimpleLine()

		p.FeedAndCut(254)
		// p.Cut()

		n, err := p.Flush()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("wrote %d bytes", n)
		return
	}

	// ESC/POS commands
	// text := []byte("Hello, world!\n")

	// cut := []byte{0x1d, 0x56, 0x00}

	p.Log("Hello, world!\n")
	p.SimpleLine()
	// Feed and cut.
	p.FeedAndCut(10)

	n, err := p.Flush()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Wrote %d bytes", n)

}
