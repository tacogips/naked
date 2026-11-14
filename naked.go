package naked

import (
	"bytes"
	"strings"
)

type textElem struct {
	text string
}

type tagElem struct {
	tag string
}

type Template struct {
	openDelim  string
	closeDelim string
	elems      []interface{}
}

func Parse(tpl string) *Template {
	return ParseWithDelims(tpl, "${", "}")
}

func ParseWithDelims(tpl, openDelim, closeDelim string) *Template {
	t := &Template{
		openDelim:  openDelim,
		closeDelim: closeDelim,
	}

	t.parseElem(tpl)

	return t
}

func (t *Template) parseElem(str string) {
	var elems []interface{}
	var text string
	isInTag := false

	for {
		if len(str) == 0 {
			break
		} else if isInTag {
			maybeTag, afterCloseTag := chomp(str, t.closeDelim)
			if str == maybeTag {
				// if no close tag
				elems = append(elems, textElem{text: maybeTag})
				break
			} else {
				elems = append(elems, tagElem{tag: strings.TrimSpace(maybeTag)})
				str = afterCloseTag
				isInTag = false
			}
		} else {
			text, str = chomp(str, t.openDelim)
			elems = append(elems, textElem{text: text})
			isInTag = true
		}
	}
	t.elems = elems
}

func chomp(s string, delim string) (string, string) {
	ss := strings.SplitN(s, delim, 2)
	if len(ss) == 1 {
		return ss[0], ""
	}
	return ss[0], ss[1]
}

func (t Template) Render(vars map[string]string) string {
	var b bytes.Buffer

	for _, elem := range t.elems {
		switch typedElem := elem.(type) {
		case textElem:
			b.WriteString(typedElem.text)
		case tagElem:
			if v, ok := vars[typedElem.tag]; ok {
				b.WriteString(v)
			} else {
				b.WriteString(t.openDelim)
				b.WriteString(typedElem.tag)
				b.WriteString(t.closeDelim)
			}
		}
	}

	return b.String()
}
