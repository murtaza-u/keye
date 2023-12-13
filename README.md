<p align="center" width="100%">
  <strong>Keye</strong> (pronounce <em>"kai"</em>) is a key-value
  database with the ability to watch over keys
</p>

<p align="center" width="100%">
  <img width="50%" src="./artwork/logo.png">
</p>

---

<div align="center"><p>
  <a href="https://godoc.org/github.com/murtaza-u/keye">
    <img alt="GoDoc" src="https://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge&logo=github&color=30b976&logoColor=D9E0EE&labelColor=302D41"/>
  </a>

  <a href="https://github.com/murtaza-u/keye/pulse">
    <img alt="Last commit" src="https://img.shields.io/github/last-commit/murtaza-u/keye?style=for-the-badge&logo=github&color=8bd5ca&logoColor=D9E0EE&labelColor=302D41"/>
  </a>

  <a href="https://github.com/murtaza-u/keye/blob/main/LICENSE">
    <img alt="License" src="https://img.shields.io/github/license/murtaza-u/keye?style=for-the-badge&logo=github&color=ee999f&logoColor=D9E0EE&labelColor=302D41" />
  </a>

  <a href="https://github.com/murtaza-u/keye/stargazers">
    <img alt="Stars" src="https://img.shields.io/github/stars/murtaza-u/keye?style=for-the-badge&logo=github&color=c69ff5&logoColor=D9E0EE&labelColor=302D41" />
  </a>

  <a href="https://github.com/murtaza-u/keye/issues">
    <img alt="Issues" src="https://img.shields.io/github/issues/murtaza-u/keye?style=for-the-badge&logo=bilibili&color=F5E0DC&logoColor=D9E0EE&labelColor=302D41" />
  </a>

  <a href="https://github.com/murtaza-u/keye">
    <img alt="Repo Size" src="https://img.shields.io/github/repo-size/murtaza-u/keye?color=%23DDB6F2&label=SIZE&logo=codesandbox&style=for-the-badge&logoColor=D9E0EE&labelColor=302D41" />
  </a>

  <a href="https://twitter.com/intent/follow?screen_name=murtaza_u_">
    <img alt="Follow on Twitter" src="https://img.shields.io/twitter/follow/murtaza_u_?style=for-the-badge&logo=twitter&color=8aadf3&logoColor=D9E0EE&labelColor=302D41" />
  </a>
</p></div>

## Deploy database server

```sh
docker run -d \
    -p 23023:23023 \
    -v "$HOME/.local/share/keye:/data" \
    murtazau/keye:23.12
```

## Client library

```sh
go get -u github.com/murtaza-u/keye
```

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/murtaza-u/keye/client"
)

func main() {
	c, err := client.New(client.Config{
		Addr:    ":23023",
		Timeout: time.Second * 5,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	key := "foo"
	val := "bar"

	keys, err := c.Put(key, []byte(val))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("modified keys:")
	for _, k := range keys {
		fmt.Println(k)
	}
}
```

Full API reference: [GoDoc](https://godoc.org/github.com/murtaza-u/keye)
