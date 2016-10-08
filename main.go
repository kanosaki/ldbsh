package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/chzyer/readline"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/urfave/cli"
	"io"
	"os"
	"strings"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var (
	DBFile string
	DB     *leveldb.DB
	Rl     *readline.Instance
)

type CommandPair struct {
	Name   string
	Action func(string)
}

var (
	Commands = []CommandPair{
		{Name: "get ", Action: get},
		{Name: "g ", Action: get},
		{Name: "put ", Action: put},
		{Name: "p ", Action: put},
		{Name: "range ", Action: lsRange},
		{Name: "prefix ", Action: lsPrefix},
		{Name: "all", Action: lsAll},
	}
)

func openDB(path string, createIfAbsent bool) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		log.Fatalf("DB open failed: %s", err)
	}
	DB = db
}

func get(param string) {
	res, err := DB.Get([]byte(param), nil)
	if err == leveldb.ErrNotFound {
		// not found
	} else if err != nil {
		log.Errorf("GET error: %s", err)
	} else {
		fmt.Println(string(res))
		_, err := Rl.Write(res)
		if err != nil {
			log.Errorf("error: %s", err)
		}
	}
}

func put(param string) {
	pair := strings.Split(param, " ")
	if len(pair) != 2 {
		log.Errorf("put <key> <value>")
	}
	err := DB.Put([]byte(pair[0]), []byte(pair[1]), nil)
	if err != nil {
		log.Errorf("PUT error: %s", err)
	}
}

func lsPrefix(param string) {
	execIterate(util.BytesPrefix([]byte(param)))
}

func lsRange(param string) {
	pair := strings.Split(param, " ")
	if len(pair) != 2 {
		log.Errorf("range <from> <to>")
	}
	execIterate(&util.Range{Start: []byte(pair[0]), Limit: []byte(pair[1])})
}

func lsAll(param string) {
	execIterate(nil)
}

func execIterate(ran *util.Range) {
	it := DB.NewIterator(ran, nil)
	for it.Next() {
		fmt.Printf("%s\t%s", string(it.Key()), string(it.Value()))
	}
	it.Release()
	err := it.Error()
	if err != nil {
		log.Errorf("%s", err)
	}
}

func help() {
	log.Infof("HELP")
}

func interpret(line string) {
	for _, c := range Commands {
		if strings.HasPrefix(line, c.Name) {
			c.Action(line[len(c.Name):])
			return
		}
	}
	log.Errorf("Unknown command: %s", line)
}

func startRepl() {
	rl, err := readline.New("> ")
	if err != nil {
		log.Fatalf("Readline error: %s", err)
	}
	Rl = rl
	for {
		line, err := rl.Readline()
		if err == io.EOF || err == readline.ErrInterrupt {
			break
		} else if err != nil {
			log.Fatalf("Readline error: %s", err)
		}
		if line == "h" || line == "help" {
			help()
		} else {
			interpret(strings.TrimSpace(line))
		}
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "ldbsh"
	app.Usage = "ldbsh -- LevelDB cli tool"
	app.Action = func(c *cli.Context) error {
		if c.NArg() < 1 {
			return errors.New("Filename requried")
		}
		DBFile = c.Args()[0]
		openDB(DBFile, true)
		defer DB.Close()
		startRepl()
		return nil
	}
	app.Run(os.Args)
}
