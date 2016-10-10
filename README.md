# ldbsh

A simple LevelDB cli utility.

## Installation

```
go get github.com/kanosaki/ldbsh
```

## Interactive mode

Following command will be open specified LevelDB dir, or create it if not exists.

```
$ ldbsh <LevelDB dir>
```

### Basic commands

```
get <key>         :: Get value from DB
put <key> <value> :: Put key-value pair to DB
list              :: (no arguments) Print all entries in DB
list <prefix>     :: List entries in DB with given prefix
load <filename>   :: Load tsv file (see Batch mode)
join <filename>   :: Join key list file (see Batch mode)
```


### Examples

```
> put somekey somevalue
> get somekey
somevalue
> put anotherkey anothervalue
> list
somekey somevalue
anotherkey anothervalue
```

## Batch mode

You can give command follow the <LevelDB dir> argument.

```
$ ldbsh <LevelDB dir> <command>
```

### Batch commands

You can use basic commands (see above) and following additional commands.

```
load :: (no arguments) Load tsv file (see Batch mode)
join :: (no arguments) Join key list file (see Batch mode)
dump :: (no arguments) Alias to 'list'
```

### Examples
`data.tsv` (A tab separated file with 2 colums)
```tsv
a 1
b 2
c 3
d 4
```

`keys.txt`
```
a
c
```

```
# You can load key-value pair from tsv.
$ ldbsh <LevelDB dir> load data.tsv
# Check the result.
$ ldbsh <LevelDB dir> list
a 1
b 2
c 3
d 4
$ ldbsh <LevelDB dir> join keys.txt
a 1
c 3
```

## License
MIT
