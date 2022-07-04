## tracktools gopro render

Renders image of GoPro GPS data

### Synopsis

Renders a image of GoPro GPS data.

```
tracktools gopro render [input mp4] [output image] [flags]
```

### Options

```
      --bearing float     override start bearing
      --distance float    override start distance
  -h, --help              help for render
      --latitude float    override start latitude
      --longitude float   override start longitude
      --min-dop float     override GPS Dilution of Precision filter
      --min-good int      override minimum good measurements
```

### Options inherited from parent commands

```
  -c, --config string   config file (Default .tracktools.toml)
  -v, --verbose count   verbose output
```

### SEE ALSO

* [tracktools gopro](tracktools_gopro.md)	 - Provides commands for manipulating GoPro videos

