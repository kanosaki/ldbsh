package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/chzyer/readline"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/urfave/cli"
	"os"
	"flag"
)

var (
	ShowIndex = flag.Bool("index", false, "Print index on output")
	ShowKey = flag.Bool("key", false, "Print key at GET/PUT command")
)

func main() {
	app := cli.NewApp()
	app.Name = "ldbsh"
	app.Usage = "ldbsh -- LevelDB cli tool"
	app.Action = func(c *cli.Context) error {
		if c.NArg() < 1 {
			return errors.New("Filename requried")
		}
		path := c.Args()[0]
		rl, err := readline.New("> ")
		if err != nil {
			log.Fatalf("Readline error: %s", err)
		}
		db, err := leveldb.OpenFile(path, nil)
		if err != nil {
			log.Fatalf("DB open failed: %s", err)
		}
		ldbsh := NewContext(db, rl)
		defer ldbsh.Close()
		ldbsh.StartRepl()
		return nil
	}
	app.Run(os.Args)
}
