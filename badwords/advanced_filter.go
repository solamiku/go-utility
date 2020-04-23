package badwords

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

var rangeChangeSet [][3]int

type KeyWordMix struct {
	Word string
	Mix  int
	Reg  *regexp.Regexp
}

//set change map
//default:1.greek to latin 2.full width to half width 3.ⓐ->a 4.⒜->a
func setChangeRange(ranges [][3]int) {
	//default
	defaults := [][3]int{
		[3]int{65345, 26, 97}, //greek to latin
		[3]int{65313, 26, 97}, //full width to half width
		[3]int{9424, 26, 97},  //ⓐ->a
		[3]int{9372, 26, 97},  //⒜->a
	}
	rangeChangeSet = append(defaults, ranges...)
}

//fitler by key word.
//the Greek alphabet automatically changing to Latin alphabet before matching.
//@mixs set the mix char length.
func FilterByKeyWord(ori, keyword string, regs ...*regexp.Regexp) string {
	var reg *regexp.Regexp
	if len(regs) > 0 {
		reg = regs[0]
	} else {
		reg = genSpecialRegexp(keyword, 2)
	}

	text := processAlphabet(ori)
	if reg == nil {
		return ori
	}
	texts := reg.FindAllString(text, -1)
	// fmt.Println(reg, texts)
	if len(texts) > 0 {
		temp := make([]string, 0, len(text))
		for _, r := range ori {
			temp = append(temp, string(r))
		}
		for _, replace := range texts {
			text = strings.Replace(text, replace,
				strings.Repeat("*", utf8.RuneCountInString(replace)), 1)
		}
		after := make([]string, 0, len(temp))
		for _, nr := range text {
			after = append(after, string(nr))
		}
		if len(after) != len(temp) {
			return ori
		}
		for k, s := range after {
			if s != "*" {
				after[k] = temp[k]
			}
		}
		ori = strings.Join(after, "")
	}
	return ori
}

//generate special regexp including mixed char recognition.
func genSpecialRegexp(match string, sep int) *regexp.Regexp {
	strs := make([]string, 0, len(match))
	for _, str := range match {
		strs = append(strs, string(str))
	}
	matchreg := strings.Join(strs, fmt.Sprintf(".{0,%d}", sep))
	reg, err := regexp.Compile(matchreg)
	if err != nil {
		logger.Errorf("compile regexp %s err:%v", match, err)
		return nil
	}
	return reg
}

//process alphabet
//change to lower string.
//change greek to lating and full width to half width.
func processAlphabet(str string) string {
	getLetterChange := func(s rune) string {
		//greek to latin
		if to, ok := Greek2Latin[string(s)]; ok {
			return to
		}
		for _, set := range rangeChangeSet {
			if s >= rune(set[0]) && s <= rune(set[0])+rune(set[1]) {
				return string(s - rune(set[0]) + rune(set[2]))
			}
		}
		return string(s)
	}
	tmp := ""
	for _, s := range str {
		tmp += strings.ToLower(getLetterChange(s))
	}

	return tmp
}

var Greek2Latin = map[string]string{
	//lower
	"α": "a",
	"β": "b",
	"γ": "y",
	"δ": "o",
	"ε": "e",
	"ζ": "e",
	"η": "n",
	"θ": "o",
	"ι": "i",
	"κ": "k",
	"λ": "n",
	"μ": "u",
	"ν": "v",
	"ξ": "e",
	"ο": "o",
	"π": "n",
	"ρ": "p",
	"ς": "c",
	"σ": "o",
	"τ": "t",
	"υ": "u",
	"φ": "o",
	"χ": "x",
	"ψ": "w",
	"ω": "w",
	"ϊ": "i",
	//upper
	"Α": "A",
	"Β": "B",
	"Γ": "C",
	"Δ": "O",
	"Ε": "E",
	"Ζ": "Z",
	"Η": "H",
	"Θ": "O",
	"Ι": "I",
	"Κ": "K",
	"Λ": "A",
	"Μ": "M",
	"Ν": "N",
	"Ξ": "E",
	"Ο": "O",
	"Π": "N",
	"Ρ": "P",
	// "Σ": "E",
	"Σ": "E",
	"Τ": "T",
	"Υ": "Y",
	"Φ": "O",
	"Χ": "X",
	"Ψ": "W",
	"Ω": "M",
	"Ϊ": "I",
}
