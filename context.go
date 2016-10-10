package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/chzyer/readline"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"io"
)

type Context struct {
	db         *leveldb.DB
	rl         *readline.Instance
	dispatcher *Dispatcher
}

func NewContext(db *leveldb.DB, rl *readline.Instance) *Context {
	c := &Context{
		db: db,
		rl: rl,
	}
	c.SetDispatcher(NewDispatcher(c.NormalCommands()))
	return c
}

func (c *Context) NormalCommands() []*Command {
	return []*Command{
		{[]string{"get", "g"}, 1, c.handleGet, "Get key"},
		{[]string{"put", "p"}, 2, c.handlePut, "Put key value"},
		{[]string{"list", "ls"}, 0, c.handleListAll, "List all entires"},
		{[]string{"list", "ls"}, 1, c.handleListPrefix, "List entries with specified"},
	}
}

func (c *Context) Close() error {
	c.rl.Close()
	return c.db.Close()
}

func (c *Context) StartRepl() error {
	for {
		line, err := c.rl.Readline()
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

func (c *Context) SetDispatcher(d *Dispatcher) {
	c.dispatcher = d
}

func (c *Context) doPrint(index uint64, key, value []byte) {
	fmt.Printf("%d\t%s\t%s\n", index, string(key), string(value))
}

func (c *Context) get(key []byte) error {
	value, err := c.db.Get(key, nil)
	if err != nil {
		return err
	}
	c.doPrint(0, key, value)
	return nil
}

func (c *Context) put(key []byte, value []byte) error {
	return c.db.Put(key, value, nil)
}

func (c *Context) list(ran *util.Range) error {
	it := c.db.NewIterator(ran, nil)
	var index uint64
	for it.Next() {
		c.doPrint(index, it.Key(), it.Value())
		index++
	}
	if it.Error() != nil {
		return it.Error()
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
