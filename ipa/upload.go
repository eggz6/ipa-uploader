package ipa

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/skip2/go-qrcode"
)

const entityKey = "entity_key"

type Result struct {
	InstallURL string
	IpaURL     string
	QrURL      string
}

type entity struct {
	bid            string
	bver           string
	btitle         string
	ipa            string
	pl             map[string]string
	zr             *zip.ReadCloser
	iconURL        string
	ipaURL         string
	manifestURL    string
	installHTMLURL string
	qrURL          string
	uploader       Uploader
}

type Uploader interface {
	Upload(ctx context.Context, r io.Reader, path string) (url string, err error)
}

// Do parse ipa and upload to remote storage.
func Do(ctx context.Context, ipa string, icon string, up Uploader) (*Result, error) {
	e := &entity{ipa: ipa, uploader: up}
	ctx = context.WithValue(ctx, entityKey, e)

	chain := []exec{
		parsePList(ctx),
		readIconData(ctx, "Payload/Runner.app/"+icon),
		uploadIPAFile(ctx),
		uploadManifest(ctx),
		uploadInstallHTML(ctx),
		uploadQRCode(ctx),
	}

	res := &Result{}

	for _, exe := range chain {
		err := decorator(ctx, exe)
		if err != nil {
			return nil, err
		}
	}

	res.InstallURL = e.installHTMLURL
	res.IpaURL = e.ipaURL
	res.QrURL = e.qrURL

	return res, nil
}

func decorator(ctx context.Context, f exec) error {
	val := ctx.Value(entityKey)
	e, ok := val.(*entity)
	if !ok {
		return errors.New("no entity in the ctx")
	}

	return f(e)
}

type exec func(e *entity) error

func parsePList(ctx context.Context) exec {
	return func(e *entity) error {
		zr, err := OpenIPA(e.ipa)
		if err != nil {
			return fmt.Errorf("open ipa failed. entity=%+v, err=%v", e, err)
		}

		e.zr = zr

		pl, err := ParseInfoPList(zr)
		if err != nil {
			return fmt.Errorf("parse info plist failed. entity=%+v, err=%v", e, err)
		}

		e.pl = pl
		e.bid, _ = pl["BUNDLE_ID"]
		e.bver, _ = pl["BUNDLE_VERSION"]
		e.btitle, _ = pl["BUNDLE_TITLE"]

		return nil
	}
}

func spellPath(e *entity) string {
	return fmt.Sprintf("/%s/%s", e.bid, e.bver)
}

func upload(ctx context.Context, e *entity, r io.Reader, fileName string) (string, error) {
	path := spellPath(e) + "/" + fileName
	url, err := e.uploader.Upload(ctx, r, path)

	if err != nil {
		return "", err
	}

	return url, nil
}

func readIconData(ctx context.Context, iconName string) exec {
	return func(e *entity) error {
		iconBytes, err := ReadIconData(e.zr, iconName)
		if err != nil {
			return fmt.Errorf("read icon data failed. iconName=%s, entity=%+v, err=%v", iconName, e, err)
		}

		url, err := upload(ctx, e, bytes.NewBuffer(iconBytes), iconName)
		if err != nil {
			return fmt.Errorf("upload icon data failed. iconName=%s, entity=%+v, err=%v", iconName, e, err)
		}

		e.iconURL = url
		e.zr.Close()

		return nil
	}
}

func uploadIPAFile(ctx context.Context) exec {
	return func(e *entity) error {
		f, err := os.Open(e.ipa)
		if err != nil {
			return fmt.Errorf("upload ipa open file failed. entity=%+v,err=%v", err)
		}

		defer f.Close()

		_, name := filepath.Split(e.ipa)

		url, err := upload(ctx, e, f, name)
		if err != nil {
			return fmt.Errorf("upload ipa file failed. entity=%+v,err=%v", err)
		}

		e.ipaURL = url

		return nil
	}
}

func uploadManifest(ctx context.Context) exec {
	return func(e *entity) error {
		content := SpellManifest(e.ipaURL, e.bid, e.bver, e.btitle, e.iconURL)

		url, err := upload(ctx, e, bytes.NewBuffer([]byte(content)), "manifest.plist")
		if err != nil {
			return fmt.Errorf("upload manifest file failed. entity=%+v,err=%v", err)
		}

		e.manifestURL = url

		return nil
	}
}

func uploadInstallHTML(ctx context.Context) exec {
	return func(e *entity) error {
		name, _ := e.pl["BUNDLE_TITLE"]
		content := SpellInstallHtml(e.manifestURL, e.iconURL, name)

		url, err := upload(ctx, e, bytes.NewBuffer([]byte(content)), "install.html")
		if err != nil {
			return fmt.Errorf("upload install html failed. entity=%+v,err=%v", err)
		}

		e.installHTMLURL = url

		return nil
	}
}

func uploadQRCode(ctx context.Context) exec {
	return func(e *entity) error {
		png, err := qrcode.Encode(e.installHTMLURL, qrcode.Medium, 256)
		if err != nil {
			return fmt.Errorf("encode install html to qrcode failed. entity=%+v,err=%v", err)
		}

		url, err := upload(ctx, e, bytes.NewBuffer(png), "qrcode.png")
		if err != nil {
			return fmt.Errorf("upload install html failed. entity=%+v,err=%v", err)
		}

		e.qrURL = url

		return nil
	}
}

//png, err := qrcode.Encode("https://example.org", qrcode.Medium, 256)
