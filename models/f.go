package models

import (
	"bytes"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
)

func MarkdownToPlainText(mdContent string) (string, error) {
	// change markdown to html
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(mdContent), &buf); err != nil {
		return "", err
	}

	// remove tag html
	p := bluemonday.StrictPolicy()
	plainText := p.Sanitize(buf.String())
	return plainText, nil
}
