# tracktools

[![Reference](https://pkg.go.dev/badge/github.com/stevenh/tracktools.svg)](https://pkg.go.dev/github.com/stevenh/tracktools) [![License](https://img.shields.io/badge/License-BSD_3--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause) [![Go Report Card](https://goreportcard.com/badge/github.com/stevenh/tracktools)](https://goreportcard.com/report/github.com/stevenh/tracktools)

tracktools is a [Golang](https://go.dev/) set of packages and command line tool that support manipulating data for a number of popular track apps.

## Features
* Automated joining of multi chapter [GoPro](https://gopro.com/) videos and format conversion.
* Convert between [HP Tuners TrackAddict](https://racerender.com/TrackAddict/) and [Harry's LapTimer](https://www.gps-laptimer.de/) data formats.

## Installing

### Command line tool
To install the [command line tool](cmd/tracktools):
```shell
go install github.com/stevenh/tracktools/cmd/tracktools@latest
```
This tool uses [Cobra](https://github.com/spf13/cobra) so supports full interactive command line help which can also be found in [markdown](cmd/tracktools/docs/tracktools.md).

If you want to use the [GoPro](https://gopro.com/) functionality you will also need a working install of [FFmpeg](https://ffmpeg.org/) which can be [downloaded here](https://ffmpeg.org/download.html). Once installed if it's not in your path you can configure the `Binary` in [.tracktools.toml](cmd/tracktools/cmd/.tracktools.toml#L2) which can be placed in your home directory.

### Golang packages
To use the packages:
```shell
go get -u github.com/stevenh/tracktools@latest
```

## Credits
A big thank you to Harry of [Harry's LapTimer](https://www.gps-laptimer.de/) for his support in answering all my questions while creating this on the [LapTimer forums](http://forum.gps-laptimer.de/).
