package util

import (
	"github.com/k3a/html2text"
	"github.com/tidwall/gjson"
	"net/url"
	"regexp"
)

var (
	scriptRegex = regexp.MustCompile(`>AF_initDataCallback[\s\S]*?</script`)
	keyRegex    = regexp.MustCompile(`(ds:\d*?)'`)
	valueRegex  = regexp.MustCompile(`data:([\s\S]*?), sideChannel: {}}\);</`)
)

// AbsoluteURL return absolute url
func AbsoluteURL(base, path string) (string, error) {
	p, err := url.Parse(path)
	if err != nil {
		return "", err
	}
	b, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	return b.ResolveReference(p).String(), nil
}

// ExtractInitData from Google HTML
func ExtractInitData(html []byte) map[string]string {
	data := make(map[string]string)
	scripts := scriptRegex.FindAll(html, -1)
	for _, script := range scripts {
		key := keyRegex.FindSubmatch(script)
		value := valueRegex.FindSubmatch(script)
		if len(key) > 1 && len(value) > 1 {
			data[string(key[1])] = string(value[1])
		}
	}
	return data
}

// GetJSONArray by path
func GetJSONArray(data string, paths ...string) []gjson.Result {
	for _, path := range paths {
		value := gjson.Get(data, path)
		if value.Exists() && value.Type != gjson.Null {
			return value.Array()
		}
	}
	return nil
}

// GetJSONValue with multiple path
func GetJSONValue(data string, paths ...string) string {
	for _, path := range paths {
		value := gjson.Get(data, path)
		if value.Exists() && value.Type != gjson.Null {
			return value.String()
		}
	}
	return ""
}

// HTMLToText return plain text from HTML
func HTMLToText(html string) string {
	return html2text.HTML2TextWithOptions(html, html2text.WithUnixLineBreaks())
}
