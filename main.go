package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
)

var (
	ShowIndex = flag.Bool("index", false, "Print index on output")
	ShowKey   = flag.Bool("key", false, "Print key at GET/PUT command")
)

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		log.Fatalf("Filename requried")
	}
	args := flag.Args()
	path := args[0]
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		log.Fatalf("DB open failed: %s", err)
	}
	ldbsh := NewContext(db)
	defer ldbsh.Close()
	if flag.NArg() > 1 {
		err = ldbsh.StartBatch(args)
	} else {
		err = ldbsh.StartRepl()
	}
	if err != nil {
		log.Fatalf("%s", err)
	}
}
