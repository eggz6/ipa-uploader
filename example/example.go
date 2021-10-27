package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/eggz6/ipa-uploader/ipa"
	"github.com/eggz6/ipa-uploader/uploader"
)

var _ipa = flag.String("ipa", "./Runner.ipa", "source ipa file path")

func main() {
	source := *_ipa
	up, err := uploader.NewOss("bucket", "endpoint", "id", "key")
	if err != nil {
		log.Fatalf("new oss upload failed. err=%v", err)
	}

	res, err := ipa.Do(context.TODO(), source, "Payload/Runner.app/AppIcon60x60@3x.png", up)
	if err != nil {
		log.Fatalf("ipa do failed. err=%v", err)
	}

	log.Println(fmt.Sprintf("success, res=%+v", res))
}
