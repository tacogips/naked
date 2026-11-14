package naked

import (
	"fmt"
	"testing"
)

func TestRender(t *testing.T) {

	{
		tpl := Parse("http//wwww.com?aa=${AB_CD}&bb=${SOME}")
		r := tpl.Render(map[string]string{"AB_CD": "1234", "SOME": "ssss"})
		if "http//wwww.com?aa=1234&bb=ssss" != r {
			t.Error(fmt.Errorf("expected http//wwww.com?aa=1234&bb=ssss but %s", r))
		}
	}

	{
		tpl := Parse("http//wwww.com?aa=${  AB_CD  }&bb=${ SOME }")
		r := tpl.Render(map[string]string{"AB_CD": "1234", "SOME": "ssss"})
		if "http//wwww.com?aa=1234&bb=ssss" != r {
			t.Error(fmt.Errorf("expected http//wwww.com?aa=1234&bb=ssss but %s", r))
		}
	}

	{
		tpl := Parse(`http//wwww.com
		?aa=${  AB_CD  }&bb=${ SOME }`)
		r := tpl.Render(map[string]string{"AB_CD": "1234", "SOME": "ssss"})
		if `http//wwww.com
		?aa=1234&bb=ssss` != r {
			t.Error(fmt.Errorf("expected http//wwww.com?aa=1234&bb=ssss but %s", r))
		}
	}

	{
		tpl := Parse("http//wwww.com?aa=${AB_CD}&bb=${SOME}")
		r := tpl.Render(map[string]string{"SOME": "ssss", "NO_IN_TPL": "ssss"})
		if "http//wwww.com?aa=${AB_CD}&bb=ssss" != r {
			t.Error(fmt.Errorf("expected http//wwww.com?aa=${AB_CD}&bb=ssss but %s", r))
		}
	}

	{ //invalid tpl
		tpl := Parse("http//wwww.com?aa=${AB_CD&bb=${SOME}")
		r := tpl.Render(map[string]string{"AB_CD": "1234", "SOME": "ssss"})
		if "http//wwww.com?aa=${AB_CD&bb=${SOME}" != r {
			t.Error(fmt.Errorf("http//wwww.com?aa=${AB_CD&bb=${SOME} but %s", r))
		}

		r2 := tpl.Render(map[string]string{"AB_CD&bb=${SOME": "1234", "SOME": "ssss"})
		if "http//wwww.com?aa=1234" != r2 {
			t.Error(fmt.Errorf("http//wwww.com?aa=1234 but %s", r))
		}
	}

	{ // custom delimiter
		tpl := ParseWithDelims("http//wwww.com?aa={{AB_CD}}&bb={{SOME}}", "{{", "}}")
		r := tpl.Render(map[string]string{"AB_CD": "1234", "SOME": "ssss"})
		if "http//wwww.com?aa=1234&bb=ssss" != r {
			t.Error(fmt.Errorf("http//wwww.com?aa=1234&bb=ssss but %s", r))
		}
	}

	{ // complicating custom delimiter
		tpl := ParseWithDelims("http//wwww.com?aa={{{AB_CD}}}&bb={{{SOME}}}", "{{", "}}")
		r := tpl.Render(map[string]string{"{AB_CD": "1234", "{SOME": "ssss"})
		if "http//wwww.com?aa=1234}&bb=ssss}" != r {
			t.Error(fmt.Errorf("http//wwww.com?aa=1234}&bb=ssss} but %s", r))
		}
	}

}
func TestParse(t *testing.T) {

	{
		tpl := Parse("http//wwww.com?aa=${AB_CD}&bb=${SOME}")
		if tpl.openDelim != "${" {
			t.Error(fmt.Errorf("expected open tag ${", tpl.openDelim))
		}
		if tpl.closeDelim != "}" {
			t.Error(fmt.Errorf("expected open tag } but %s", tpl.closeDelim))
		}

		if len(tpl.elems) != 4 {
			t.Error(fmt.Errorf("expected aaaa but %s", len(tpl.elems)))
		}

		{
			first := tpl.elems[0]
			if te, ok := first.(textElem); !ok {
				t.Error(fmt.Errorf("expected text elem but %#v", first))
			} else {
				if te.text != "http//wwww.com?aa=" {
					t.Error(fmt.Errorf(`expected "http//wwww.com?aa=" but %s`, te.text))
				}
			}
		}

		{
			second := tpl.elems[1]
			if te, ok := second.(tagElem); !ok {
				t.Error(fmt.Errorf("expected text elem but %#v", second))
			} else {
				if te.tag != "AB_CD" {
					t.Error(fmt.Errorf("expected AB_CD but %s", te.tag))
				}
			}
		}

		{
			third := tpl.elems[2]
			if te, ok := third.(textElem); !ok {
				t.Error(fmt.Errorf("expected text elem but %#v", third))
			} else {
				if te.text != "&bb=" {
					t.Error(fmt.Errorf("expected &bb= but %s", te.text))
				}
			}
		}

		{
			fourth := tpl.elems[3]
			if te, ok := fourth.(tagElem); !ok {
				t.Error(fmt.Errorf("expected text elem but %#v", fourth))
			} else {
				if te.tag != "SOME" {
					t.Error(fmt.Errorf("expected &bb= but %s", te.tag))
				}
			}
		}

	}
}

func TestChomp(t *testing.T) {
	{
		f, s := chomp("aaaa{{bbb", "{{")
		if f != "aaaa" {
			t.Error(fmt.Errorf("expected aaaa but %s", f))
		}

		if s != "bbb" {
			t.Error(fmt.Errorf("expected bbb but %s", s))
		}
	}

	{
		f, s := chomp("aaaa{{", "{{")
		if f != "aaaa" {
			t.Error(fmt.Errorf("expected aaaa but %s", f))
		}

		if s != "" {
			t.Error(fmt.Errorf("expected empty but %s", s))
		}
	}

	{
		f, s := chomp("aaaa{{", "{")
		if f != "aaaa" {
			t.Error(fmt.Errorf("expected aaaa but %s", f))
		}

		if s != "{" {
			t.Error(fmt.Errorf("expected { but %s", s))
		}
	}

	{
		f, s := chomp("", "{")
		if f != "" {
			t.Error(fmt.Errorf("expected empty but %s", f))
		}

		if s != "" {
			t.Error(fmt.Errorf("expected empty but %s", s))
		}
	}

	{
		f, s := chomp("{", "{")
		if f != "" {
			t.Error(fmt.Errorf("expected empty but %s", f))
		}

		if s != "" {
			t.Error(fmt.Errorf("expected empty but %s", s))
		}
	}

	{
		f, s := chomp("{bb", "{")
		if f != "" {
			t.Error(fmt.Errorf("expected empty but %s", f))
		}

		if s != "bb" {
			t.Error(fmt.Errorf("expected bb but %s", s))
		}
	}

	{
		f, s := chomp("あいううう", "うう")
		if f != "あい" {
			t.Error(fmt.Errorf("expected あい but %s", f))
		}

		if s != "う" {
			t.Error(fmt.Errorf("expected う but %s", s))
		}
	}

}
