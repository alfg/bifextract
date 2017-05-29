# bifextract [![Build Status](https://travis-ci.org/alfg/bifextract.svg?branch=master)](https://travis-ci.org/alfg/bifextract)
`bifextract` is a CLI utility for extracting images from a [BIF](https://sdkdocs.roku.com/display/sdkdoc/Trick+Mode+Support) file.

## Install from Source
```
go get github.com/alfg/bifextract
./bin/bifextract version
```

## Install from Homebrew
```
brew cask alfg/tap
brew install alfg/tap/bifextract
bifextract version
```

## Usage
`bifextract <file-path|url> <output-dir>`

## Example

`bifextract gladiator.bif gladiator`

This will parse `gladiator.bif`, create the `gladiator` directory and write the image frames in sequential order.

```
└─gladiator/
    ├─ frame_1.jpg
    ├─ frame_2.jpg
    ├─ ... 
    └─ frame_n-1.jpg
    └─ frame_n.jpg
```

## Develop

```
git clone git@github.com:alfg/bifextract.git
export GOPATH=$HOME/path/to/project
cd /to/project
go run main.go
```

## Resources
* https://sdkdocs.roku.com/display/sdkdoc/Trick+Mode+Support

## License
MIT
