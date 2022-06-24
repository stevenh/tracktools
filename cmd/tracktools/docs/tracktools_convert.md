## tracktools convert

Convert between track app formats.

### Synopsis

Convert between different track app logging formats

```
tracktools convert input-file output-file [flags]
```

### Options

```
      --compress           Override Compress option for output
      --decoder string     Override Decoder for the input
      --encoder string     Override Encoder for the output
  -h, --help               help for convert
      --note string        Override Note for the output
      --start-date date    Override StartDate option for output (format YYYY-MM-DD) (default 0001-01-01)
      --tags stringArray   Override Tags for the output
      --track string       Override Track for the output
      --vehicle string     Override Vehicle for the output
```

### Options inherited from parent commands

```
  -c, --config string   config file (Default .tracktools.toml)
  -v, --verbose count   verbose output
```

### SEE ALSO

* [tracktools](tracktools.md)	 - A set of tools for creating track videos

