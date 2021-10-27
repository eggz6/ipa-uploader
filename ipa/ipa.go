package ipa

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/eggz6/ipa-uploader/plist"
	"github.com/eggz6/utils/format"
)

func OpenIPA(path string) (*zip.ReadCloser, error) {
	zr, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}

	return zr, nil
}

func ParseInfoPList(zr *zip.ReadCloser) (map[string]string, error) {
	var plistReader io.ReadCloser
	reg, _ := regexp.Compile(`(Payload\/.*\.app\/Info\.plist)$`)
	for _, f := range zr.File {
		if reg.Match([]byte(f.Name)) {
			r, err := f.Open()
			if err != nil {
				return nil, err
			}

			plistReader = r
			break
		}
	}

	if plistReader == nil {
		return nil, fmt.Errorf("no plist file in path. ")
	}

	defer plistReader.Close()

	p, err := plist.ReadFrom(plistReader)
	if err != nil {
		return nil, err
	}

	res := make(map[string]string, 0)

	bundleID, _ := p.GetString("CFBundleIdentifier")
	name, _ := p.GetString("CFBundleName")
	version, _ := p.GetString("CFBundleVersion")

	res["BUNDLE_ID"] = bundleID
	res["BUNDLE_VERSION"] = version
	res["BUNDLE_TITLE"] = name

	return res, nil
}

func ReadIconData(zr *zip.ReadCloser, iconName string) ([]byte, error) {
	var icon io.ReadCloser
	for _, f := range zr.File {
		if strings.Contains(f.Name, iconName) {
			r, err := f.Open()
			if err != nil {
				return nil, err
			}

			icon = r

			break
		}
	}

	if icon == nil {
		return nil, fmt.Errorf("no icon file")
	}
	defer icon.Close()

	data, err := ioutil.ReadAll(icon)

	return data, err
}

func SpellManifest(ipaURL, bid, bver, title, icon string) string {
	manifestStr := format.ReplaceArgs(manifest, map[string]string{
		"IPA_URL":        ipaURL,
		"BUNDLE_ID":      bid,
		"BUNDLE_VERSION": bver,
		"TITLE":          title,
		"ICON_URL":       icon,
	})

	return manifestStr
}

func SpellInstallHtml(manifestURL, iconURL, title string) string {
	htmlStr := format.ReplaceArgs(html, map[string]string{
		"TITLE":        title,
		"MANIFEST_URL": manifestURL,
		"ICON_URL":     iconURL,
	})

	return htmlStr
}
