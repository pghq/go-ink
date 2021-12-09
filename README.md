# go-ink
Golang text translation library.

## Installation

go-ink may be installed using the go get command:

```
go get github.com/pghq/go-ink
```
## Usage

```
import (
    "github.com/pghq/go-ink"
    "github.com/pghq/go-ink/lang
)
```

To create a new client:

```
lin := ink.NewLinguist("your-api-key")
text, err := lin.Translate(context.TODO(), "Hello, world!", lang.German)
if err := radar.Error(); err != nil{
    panic(err)
}
```

## Powered by
* DeepL Pro - https://www.deepl.com/pro 
