package main

import (
	"regexp"
	"strconv"
	"strings"
)

type WordDetails struct {
	word, color, rgba       string
	red, green, blue, alpha string
	weight                  float32
	regexp                  *regexp.Regexp
}

var (
	escapeRegex = regexp.MustCompile(`[\-\[\\\]\{\}\(\)\*\+\?\.\,\^\$\|\#]`)
)

func parseColorPart(part string) byte {
	val, _ := strconv.ParseUint(part, 16, 8)
	return byte(val)
}
func parseColor(color string) (r, g, b byte) {
	r = parseColorPart(color[1:2])
	g = parseColorPart(color[3:4])
	b = parseColorPart(color[5:6])
	return
}

func newWordDetails(word, color string, weight float32, rgba string) *WordDetails {
	escaped := string(escapeRegex.ReplaceAll([]byte(word), []byte(`\\$&`)))

	regexpString := strings.Join([]string{`\b`, escaped, `\b`}, ``)

	rgbaSlice := strings.Split(rgba, ",")

	return &WordDetails{
		word:   word,
		color:  color,
		rgba:   rgba,
		regexp: regexp.MustCompile(regexpString),
		weight: weight,
		red:    rgbaSlice[0],
		green:  rgbaSlice[1],
		blue:   rgbaSlice[2],
		alpha:  rgbaSlice[3],
	}
}

func (w *WordDetails) Word() string {
	return w.word
}

func (w *WordDetails) Color() string {
	return w.color
}

func (w *WordDetails) Rgba() string {
	return w.rgba
}

func (w *WordDetails) Regexp() *regexp.Regexp {
	return w.regexp
}

func (w *WordDetails) Weight() float32 {
	return w.weight
}

func (w *WordDetails) Red() string {
	return w.red
}

func (w *WordDetails) Green() string {
	return w.green
}

func (w *WordDetails) Blue() string {
	return w.blue
}

func (w *WordDetails) Alpha() string {
	return w.alpha
}

var WordIndex map[string]*WordDetails

var WordList = [...]*WordDetails{
	newWordDetails(`happy`, `#ffcd05`, 0.65, `255,205,5,30`),
	newWordDetails(`sad`, `#0e1030`, 1, `14,15,40,0`),
	newWordDetails(`cheeky`, `#ee3c96`, 1, `238,60,150,40`),
	newWordDetails(`angry`, `#ed2024`, 1, `237,32,5,56`),
	newWordDetails(`disappointed`, `#023334`, 1, `2,71,15,30`),
	newWordDetails(`lonely`, `#e6e7e8`, 1, `28,67,19,0`),
	newWordDetails(`depressed`, `#22205f`, 1, `22,30,94,0`),
	newWordDetails(`annoyed`, `#f05323`, 1, `241,84,6,255`),
	newWordDetails(`embarrassed`, `#eac2db`, 1, `181,194,79,35`),
	newWordDetails(`love`, `#cb2026`, 0.242, `203,19,2,15`),
	newWordDetails(`amazing`, `#f6eb14`, 1, `246,235,5,111`),
	newWordDetails(`glad`, `#359946`, 1, `53,175,10,0`),
	newWordDetails(`bored`, `#676766`, 1, `116,255,42,38`),
	newWordDetails(`tired`, `#cb9865`, 1, `203,152,12,0`),
	newWordDetails(`excited`, `#ee3a68`, 1, `238,58,58,180`),
	newWordDetails(`dead`, `#000000`, 1, `0,0,0,0`),
	newWordDetails(`confused`, `#0d9a48`, 1, `13,154,8,11`),
	newWordDetails(`grumpy`, `#0e1030`, 1, `10,19,33,0`),
	newWordDetails(`desperate`, `#663300`, 1, `38,21,0,0`),
	newWordDetails(`yay`, `#c673af`, 1, `198,115,121,85`),
	newWordDetails(`fun`, `#7c287d`, 1, `124,40,125,0`),
	newWordDetails(`ace`, `#478fcd`, 0.7, `124,143,255,0`),
	newWordDetails(`rage`, `#e63725`, 1, `230,55,4,0`),
	newWordDetails(`boss`, `#a4ce39`, 1, `164,206,6,0`),
	newWordDetails(`badass`, `#604d80`, 1, `96,77,128,0`),
	newWordDetails(`wicked`, `#d36127`, 1, `211,97,6,151`),
	newWordDetails(`sick`, `#cfdd2e`, 1, `207,211,14,0`),
	newWordDetails(`alright`, `#5691cd`, 1, `86,145,205,20`),
	newWordDetails(`awesome`, `#ef3c24`, 1, `239,60,2,26`),
	newWordDetails(`vibing`, `#ed1968`, 1, `237,25,104,255`),
	newWordDetails(`comfortable`, `#009e9e`, 1, `0,162,43,0`),
	newWordDetails(`rough`, `#343416`, 1, `53,121,13,0`),
	newWordDetails(`fantastic`, `#df3a42`, 1, `223,58,10,106`),
	newWordDetails(`hanging`, `#048248`, 1, `4,130,10,20`),
	newWordDetails(`meh`, `#a7a9ac`, 1, `116,229,140,0`),
	newWordDetails(`mardy`, `#023033`, 1, `2,48,9,8`),
	newWordDetails(`undesisive`, `#3953a4`, 1, `57,83,164,0`),
	newWordDetails(`uncomfortable`, `#caaaac`, 1, `202,170,172,0`),
	newWordDetails(`great`, `#fdc010`, 1, `253,192,3,64`),
	newWordDetails(`lol`, `#ee3a68`, 0.18, `238,58,104,100`),
	newWordDetails(`lmao`, `#d4344f`, 0.8, `146,52,29,0`),
	newWordDetails(`rubbish`, `#985a25`, 1, `152,90,4,22`),
	newWordDetails(`hate`, `#321110`, 0.83, `50,17,0,0`),
	newWordDetails(`enjoy`, `#cbdb2a`, 1, `203,219,5,25`),
	newWordDetails(`silly`, `#ed1568`, 1, `237,21,104,125`),
	newWordDetails(`:p`, `#ee3696`, 0.8, `143,54,104,108`),
	newWordDetails(`:)`, `#f3ec19`, 0.75, `243,236,6,0`),
	newWordDetails(`:(`, `#0e1130`, 0.7, `14,17,36,0`),
	newWordDetails(`:s`, `#5d8e49`, 0.7, `60,131,0,0`),
	newWordDetails(`:D`, `#993300`, 0.8, `130,58,0,45`),
}

func init() {
	WordIndex = make(map[string]*WordDetails)
	for _, word := range WordList {
		WordIndex[word.word] = word
	}
}
