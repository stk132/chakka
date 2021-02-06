# chakka

```
usage: chakka [<flags>] <command> [<args> ...]

fireworq queue and routing setting save, apply cli

Flags:
--help     Show context-sensitive help (also try --help-long and --help-man).
--version  Show application version.

Commands:
help [<command>...]
save --from=FROM [<flags>]
apply --to=TO [<flags>]
```

## usage

### save setting

```shell
chakka save --from=http://localhost:8888 --queue-file=./queues.json --routing-file=./routings.json
```

### apply setting

```shell
chakka apply --to=http://localhost:8888 --queue-file=./queues.json --routing-file=./routings.json
```