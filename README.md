# ipa-uploader

# Desc 

upload local ipa to remote server

# Usage

- implement ipa.Uploader interface

- call ipa.Do

# For example use local file system

```
package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/eggz6/ipa-uploader/ipa"
	"github.com/eggz6/ipa-uploader/uploader"
	"github.com/gin-gonic/gin"
)

var _ipa = flag.String("ipa", "./Runner.ipa", "source ipa file path")

func main() {
	source := *_ipa
	up := uploader.NewFS()

	res, err := ipa.Do(context.TODO(), source, "AppIcon60x60@3x.png", up)
	if err != nil {
		log.Fatalf("ipa do failed. err=%v", err)
	}

	log.Println(fmt.Sprintf("success, res=%+v", res))
	router := gin.Default()
	router.Static("output", "output")

	router.Run(":8080")
}

```


