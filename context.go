package main

import (
	"bufio"
	"bytes"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/chzyer/readline"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"io"
	"os"
)

var (
	TabBytes     = []byte("\t")
	NewlineBytes = []byte("\n")
)

type Context struct {
	db         *leveldb.DB
	dispatcher *Dispatcher
}

func NewContext(db *leveldb.DB) *Context {
	c := &Context{
		db: db,
	}
	c.SetDispatcher(c.NormalCommands())
	return c
}

func (c *Context) NormalCommands() []*Command {
	return []*Command{
		{[]string{"get", "g"}, 1, c.handleGet, "Get key"},
		{[]string{"put", "p"}, 2, c.handlePut, "Put key value"},
		{[]string{"list", "ls", "l"}, 0, c.handleListAll, "List all entires"},
		{[]string{"list", "ls", "l"}, 1, c.handleListPrefix, "List entries with specified"},
	}
}

func (c *Context) BatchCommands() []*Command {
	return []*Command{
		{[]string{"join"}, 1, c.handleJoinFilename, "Get key"},
		{[]string{"load"}, 1, c.handleLoadFilename, "Put key value"},
		{[]string{"join"}, 0, c.handleJoinStdin, "Get key"},
		{[]string{"load"}, 0, c.handleLoadStdin, "Put key value"},
		{[]string{"dump"}, 0, c.handleListAll, "List all entires"},
	}
}

func (c *Context) Close() error {
	return c.db.Close()
}

func (c *Context) StartRepl() error {
	rl, err := readline.New("> ")
	if err != nil {
		log.Fatalf("Readline error: %s", err)
	}
	for {
		line, err := rl.Readline()
		if err == io.EOF || err == readline.ErrInterrupt {
			break
		} else if err != nil {
			log.Fatalf("Readline error: %s", err)
		}
		err = c.dispatcher.Call(line)
		if err != nil {
			log.Errorf("%s", err)
		}
	}
	return nil
}

func (c *Context) StartBatch(args []string) error {
	c.SetDispatcher(
		append(
			c.NormalCommands(),
			c.BatchCommands()...,
		))
	com, params := c.dispatcher.Lookup(args[1:])
	if com == nil {
		log.Fatalf("Unknown command: %s/%d", args[1], len(args)-2)
	}
	err := com.Action(params)
	return err
}

func (c *Context) SetDispatcher(coms []*Command) {
	c.dispatcher = NewDispatcher(coms)
}

func (c *Context) doPrint(index uint64, key, value []byte, printKey bool) {
	if *ShowIndex {
		fmt.Printf("%d\t", index)
	}
	if printKey {
		fmt.Printf("%s\t", string(key))
	}
	fmt.Printf("%s\n", string(value))
}

func (c *Context) get(key []byte) error {
	value, err := c.db.Get(key, nil)
	if err != nil {
		return err
	}
	c.doPrint(0, key, value, *ShowKey)
	return nil
}

func (c *Context) put(key []byte, value []byte) error {
	return c.db.Put(key, value, nil)
}

func (c *Context) list(ran *util.Range) error {
	it := c.db.NewIterator(ran, nil)
	var index uint64
	for it.Next() {
		c.doPrint(index, it.Key(), it.Value(), true)
		index++
	}
	if it.Error() != nil {
		return it.Error()
	}
	return nil
}

func (c *Context) join(input io.Reader) error {
	printIndex := *ShowIndex
	scanner := bufio.NewScanner(input)
	defer os.Stdout.Sync()
	var counter uint64
	for scanner.Scan() {
		counter++
		key := scanner.Bytes()
		if len(key) == 0 {
			continue
		}
		value, err := c.db.Get(key, nil)
		if err != nil {
			return err
		}
		if printIndex {
			fmt.Printf("%d\t", counter)
		}
		os.Stdout.Write(key)
		os.Stdout.Write(TabBytes)
		os.Stdout.Write(value)
		os.Stdout.Write(NewlineBytes)
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}
	return nil
}

func (c *Context) load(input io.Reader) error {
	scanner := bufio.NewScanner(input)
	var counter uint64
	for scanner.Scan() {
		counter++
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		pivot := bytes.Index(line, TabBytes)
		if pivot < 0 {
			log.Errorf("FormatError: NO TAB (line: %d)", counter)
		} else {
			key := line[:pivot]
			value := line[pivot+1:]
			err := c.db.Put(key, value, nil)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Context) handleGet(params []string) error {
	key := []byte(params[0])
	return c.get(key)
}

func (c *Context) handlePut(params []string) error {
	key := []byte(params[0])
	value := []byte(params[1])
	return c.put(key, value)
}

func (c *Context) handleListAll(params []string) error {
	return c.list(nil)
}

func (c *Context) handleListPrefix(params []string) error {
	prefix := []byte(params[0])
	return c.list(util.BytesPrefix(prefix))
}

func (c *Context) handleJoinFilename(params []string) error {
	filename := params[0]
	if filename == "-" {
		return c.handleJoinStdin(params)
	}
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return c.join(f)
}

func (c *Context) handleLoadFilename(params []string) error {
	filename := params[0]
	if filename == "-" {
		return c.handleLoadStdin(params)
	}
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return c.load(f)
}

func (c *Context) handleJoinStdin(params []string) error {
	return c.join(os.Stdin)
}

func (c *Context) handleLoadStdin(params []string) error {
	return c.load(os.Stdin)
}
