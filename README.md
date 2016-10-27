## naked
Deadly simple template engine in Go

- NO snippet
- NO block
- NO function
- NO struct mapping
- NO method mapping
- NO loading template from file
- NO escaping
- NO validation
- NO any other rich feature
- SO very fast

## example
```
// example
tpl := Parse("http//wwww.com?aa=${AB_CD}&bb=${SOME}")

r := tpl.Render(map[string]string{"AB_CD": "1234", "SOME": "ssss"})
// r => "http//wwww.com?aa=1234&bb=ssss"
```


## Feature
- [ ] only render with "map[string]string"
- [ ] change delimiter

## TODO
- [ ] delimter can contains rune
- [ ] fill test case
- [ ] benchmarking

