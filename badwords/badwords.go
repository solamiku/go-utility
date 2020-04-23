/*
	屏蔽词函数
*/

package badwords

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	slogger "github.com/solamiku/go-utility/logger"

	simplejson "github.com/bitly/go-simplejson"
)

var logger slogger.Logger = slogger.DefaultLogger

func SetLogger(l slogger.Logger) {
	logger = l
}

// 敏感词库敏感词
var badWords []string
var BadwordsReplacer *strings.Replacer

// 扩充的global.xml特殊符号敏感词,
// 目的是为了在屏蔽普通敏感词时，可以根据参数决定是否屏蔽特殊符号敏感词
// 比如名字不允许出现句号。(最初为影响某些上报BI系统出现符号处理异常）
// 但是正常聊天中可以出现句号。
var extraBadWords []string
var ExtraBadwordsReplacer *strings.Replacer

// db读取的高级敏感词扩充
var dbExtraBadWords []string
var dbExtraRegs []*regexp.Regexp
var dbExtraRegsGrp map[int][]*regexp.Regexp

// db读取的，高级替换（正则和字符集映射）敏感词扩充
var specialBadWords []*KeyWordMix

// emoji的编码转换-用于db保存等
var regEmoji *regexp.Regexp
var regEmojiExtract *regexp.Regexp

// 有效字符集检测相关
var validSet []*unicode.RangeTable

// 仅提取中文，英文，数字的正则
var regIgnoreOther *regexp.Regexp

// 去除多余空格，制表符，换行的正则
var regRepalceRemoveSpace *regexp.Regexp

func init() {
	// emoji表情的数据表达式
	regEmoji = regexp.MustCompile("\\[[\\\\u0-9a-zA-Z]+\\]")
	// 提取emoji数据表达式
	regEmojiExtract = regexp.MustCompile("\\[\\\\u|]")
	// 提取中文，数字，英文
	regIgnoreOther = regexp.MustCompile("[^\u4e00-\u9fa5|0-9a-zA-Z]")
	var err error
	regRepalceRemoveSpace, err = regexp.Compile("\\s{2,}|\t|\n")
	if err != nil {
		log.Fatalf("compile replace remove space err:%v", err)
	}
	// 英文字母
	enLatin := unicode.RangeTable{
		R16: []unicode.Range16{
			{0x0041, 0x005a, 1}, //A-Z
			{0x0061, 0x007a, 1}, //a-z
		},
		LatinOffset: 6,
	}
	// 普通数学符号
	normalMath := unicode.RangeTable{
		R16: []unicode.Range16{
			{0x002b, 0x002b, 1}, //+
			{0x003c, 0x003e, 1}, //<=>
			{0x007c, 0x007e, 2}, //|~
		},
	}
	// 普通数字
	normalNumber := unicode.RangeTable{
		R16: []unicode.Range16{
			{0x0030, 0x0039, 1}, //0-9
		},
	}
	validSet = []*unicode.RangeTable{
		unicode.Han, unicode.P, unicode.White_Space,
		&normalMath, &enLatin, &normalNumber,
	}
	// 设置整体转换表范围
	setChangeRange(nil)
}

// 初始化敏感词

// @path - 敏感词文件txt路径
// @jdata - 一般为配置在数据库的json格式，涉及正则，特殊敏感词等
// @extraPunctuation - 一般为配置在global.xml 为了特别区分可以skip掉检测的一些敏感词
func InitBadwords(path string, jdata *simplejson.Json, extraPunctuation []string) error {
	badwordsLocal, err := parseBadWordsFile(path)
	if err != nil {
		return errors.New(fmt.Sprintf("load badwords from file err:%v", err))
	}
	specials := make([]*KeyWordMix, 0, 10)
	dbWords := make([]string, 0, 1)
	if jdata == nil {
		// 测试说明用的json文件
		// 1.ex 支持三种格式：
		// 	@1 - `reg:xxx` 正则模式
		//  @2 - `reg[\d+]:xxx` 正则组模式
		//  @3 - `xxx` 普通敏感词模式 (与badwords.txt内容合并)
		// 2.sp：
		//  [xxxx, 2] 为一个完整词汇构造支持中间插入任意字符的正则表达式->[x.{0,2}x.{0,2}x....]
		// 3.set:
		//  [起始字符码，范围，起始转换] 即当某个字符为大于等于起始字符吗，并且在范围以内，则将其转换为起始转换对应的字符
		//  eg:[9424, 26, 97] 该set即为将圈英文转换为小写英文:ⓐ->a
		testJdat, err := simplejson.NewJson([]byte(`{
			"ex": [
				"reg:装.{0,3}备.{0,3}私.{0,3}聊",
				"reg1:[loO0-9一二三四五六七八九零]{8,}",
				"烈士"
			],
			"sp":[
				["星环互动", 2]
			],
			"set":[]
		}`))
		jdata = testJdat
		if jdata == nil {
			return errors.New(fmt.Sprintf("json data is nil. err:%v", err))
		}
	}
	exKey := "ex"
	spKey := "sp"
	setKey := "set"
	for k := range jdata.Get(exKey).MustArray() {
		dbWords = append(dbWords, jdata.Get(exKey).GetIndex(k).MustString())
	}
	for k := range jdata.Get(spKey).MustArray() {
		specials = append(specials, &KeyWordMix{
			Word: jdata.Get(spKey).GetIndex(k).GetIndex(0).MustString(),
			Mix:  jdata.Get(spKey).GetIndex(k).GetIndex(1).MustInt(),
		})
	}
	sets := make([][3]int, 0, 10)
	for k := range jdata.Get(setKey).MustArray() {
		sets = append(sets, [3]int{
			jdata.Get(setKey).GetIndex(k).GetIndex(0).MustInt(),
			jdata.Get(setKey).GetIndex(k).GetIndex(1).MustInt(),
			jdata.Get(setKey).GetIndex(k).GetIndex(2).MustInt(),
		})
	}

	setChangeRange(sets)
	badwordsLocal = append(badwordsLocal, dbWords...)
	setBadWrods(badwordsLocal, extraPunctuation)
	setSpecialBadWords(specials)
	setDBBadWrodsExtra(dbWords)
	return nil
}

// 格式化屏蔽词文件
func parseBadWordsFile(path string) ([]string, error) {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// 查找第一行，第一行全是分隔符定义
	idx_r := -1
	for i := 0; i < len(bs); i++ {
		if bs[i] == '\r' || bs[i] == '\n' {
			idx_r = i
			break
		}
	}
	if idx_r <= 0 {
		return nil, fmt.Errorf("failed to find separator.")
	}
	// 分隔符
	m_sep := make(map[rune]int, 10)
	s_sep := string(bs[:idx_r])
	for _, r := range s_sep {
		m_sep[r] = 1
	}
	// split
	f := func(c rune) bool {
		if _, ok := m_sep[c]; ok {
			return true
		}
		if c == ' ' || c == '　' {
			return false
		}
		return unicode.IsSpace(c)
	}
	badWords := strings.FieldsFunc(string(bs[idx_r:]), f)
	for idx := len(badWords) - 1; idx >= 0; idx-- {
		badWords[idx] = strings.TrimSpace(badWords[idx])
		if len(badWords[idx]) == 0 {
			badWords = append(badWords[:idx], badWords[idx+1:]...)
		}
	}
	return badWords, nil
}

// 初始化基本屏蔽词,扩展的标点
func setBadWrods(aBadWords, aExtraBadWords []string) {
	badWords = aBadWords
	extraBadWords = aExtraBadWords

	bads := make([]string, 0, len(badWords)*2)
	exbads := make([]string, 0, len(badWords)*2)
	for _, v := range aBadWords {
		bads = append(bads, v, strings.Repeat("*", utf8.RuneCountInString(v)))
	}
	for _, v := range aExtraBadWords {
		exbads = append(exbads, v, strings.Repeat("*", utf8.RuneCountInString(v)))
	}

	BadwordsReplacer = strings.NewReplacer(bads...)
	ExtraBadwordsReplacer = strings.NewReplacer(exbads...)
}

//初始化特殊的替换屏蔽词
func setSpecialBadWords(aSepcialBadWords []*KeyWordMix) {
	tempSp := make([]*KeyWordMix, 0, len(aSepcialBadWords))
	for k, v := range aSepcialBadWords {
		reg := genSpecialRegexp(v.Word, v.Mix)
		if reg == nil {
			continue
		}
		v.Reg = reg
		tempSp = append(tempSp, aSepcialBadWords[k])
	}
	specialBadWords = tempSp
}

//初始化db加载的特殊屏蔽词（用于去除标点，空格后检测是否存在屏蔽词）
func setDBBadWrodsExtra(words []string) {
	normals := make([]string, 0, 1)
	regs := make([]*regexp.Regexp, 0, 1)
	grpRegs := make(map[int][]*regexp.Regexp, 1)
	regMatch := regexp.MustCompile("reg([0-9]*)")
	for _, word := range words {
		strs := strings.Split(word, ":")
		if len(strs) >= 2 && strings.Contains(strs[0], "reg") {
			reg, err := regexp.Compile(strs[1])
			if err != nil {
				logger.Errorf("compile reg:%s err:%v", strs[1], err)
				continue
			}

			matchs := regMatch.FindStringSubmatch(strs[0])
			grp := 0
			if len(matchs) > 1 {
				grp, _ = strconv.Atoi(matchs[1])
			}

			if grp == 0 {
				regs = append(regs, reg)
			} else {
				if _, ok := grpRegs[grp]; !ok {
					grpRegs[grp] = make([]*regexp.Regexp, 0, 1)
				}
				grpRegs[grp] = append(grpRegs[grp], reg)
			}
		} else {
			normals = append(normals, word)
		}
	}
	dbExtraBadWords = normals
	dbExtraRegs = regs
	dbExtraRegsGrp = grpRegs
}

/*----------------------------------------------------------------------------*/
// 敏感词检测各大接口简要说明:
// 1.普通敏感词检测用  BadwordsCheck(string, ...boool)， 它包含:
//   @1 badwords.txt中的敏感词
//   @2 db中的ex铭感次
//   @3 global中的特殊标点敏感词
//	可选参数可以选择是否跳过@3的检测，比如名字不允许出现global中定义的特殊标点，但是聊天中可以

// 2.高级敏感词检测， BadwordsCheckAdvanced(str, bool, ...int) 通常用于聊天中内容检测
//   step1.用regIgnoreOther生成只包含中文，数字，英文的字符串after
//   step2.检测是否存在db中的敏感词，注意此处并不会检测属于badwords.txt的敏感词
//   step3.根据db正则匹配是否有敏感内容（大多为广告）
//   step4.根据正则组对特殊正则组进行敏感内容检测
//   step5.根据参数决定是否检测global中的特殊标点敏感词

// 3.普通敏感词替换: BadwordsReplace()
//   step1.先进行badwords.txt和db中的ex敏感词替换成对应长度的**
//   step2.根据参数进行global的ex标点进行替换*
//   step3.根据db中的sp对整段包含敏感内容的信息进行替换等长度*

// 总结:
// 1.取名之类的需要严格处理的，执行前先调用BadwordsCheck，然后决定是否能进行
// 2.聊天时：
//	 聊天信息因为可以用很多特殊的数字和英文甚至是中文编码来尝试绕过广告的正则检测，
//	 原本使用setChangeRange中规则来进行编码转换,但是架不住编码集的庞大，所以对聊天信息直接先进行合法集合检测
//   用CheckValidCharacter接口约束聊天内容只能包含简单的中文，数字，英文，简化敏感内容检测。
//   然后用BadwordsCheckAdvanced对信息中的广告信息加以匹配，进行约束,因为这里只检测正则，
//   所以广告约束就算只是一个词，都要以正则形式配置，这样才能与合并到badwords.txt里面的内容区分
//   最后再用BadwordsReplace对聊天信息中的剩余敏感词转为"*"号，允许其发出来。

// -------------------------------------
// 是否包含屏蔽词，返回true表示有
// skip表示跳过额外敏感词检测
func BadwordsCheck(str string, skip ...bool) bool {
	for _, s := range badWords {
		if strings.Index(str, s) >= 0 {
			return true
		}
	}
	// skip表示跳过额外敏感词检测
	bSkip := false
	if len(skip) > 0 {
		bSkip = skip[0]
	}
	if !bSkip {
		for _, s := range extraBadWords {
			if strings.Index(str, s) >= 0 {
				return true
			}
		}
	}
	return false
}

//高级敏感词检测
//@str需要先满足规定字符集，只取中英文数字进行敏感词匹配
func BadwordsCheckAdvanced(str string, extra bool, grps ...int) int {
	//去除特殊符号后进行dbExtraBadWords的敏感词匹配
	if len(dbExtraRegs) > 0 || len(dbExtraRegsGrp) > 0 { // || len(dbExtraBadWords) > 0
		after := regIgnoreOther.ReplaceAllString(str, "")
		after = strings.ToLower(after)
		// for _, s := range dbExtraBadWords {
		// 	if strings.Index(after, s) >= 0 {
		// 		return 1
		// 	}
		// }
		for _, reg := range dbExtraRegs {
			if reg.MatchString(after) {
				return 2
			}
		}

		for _, grp := range grps {
			for _, reg := range dbExtraRegsGrp[grp] {
				if reg.MatchString(after) {
					return 3
				}
			}
		}

	}
	if extra {
		for _, s := range extraBadWords {
			if strings.Index(str, s) >= 0 {
				return 4
			}
		}
	}
	return 0
}

// 屏蔽词替换
// skip表示跳过额外敏感词检测
func BadwordsReplace(str string, skip ...bool) string {
	bSkip := false
	if len(skip) > 0 {
		bSkip = skip[0]
	}
	//内容替换，将换行制表替换成空格，多空格只保留一个
	str = regRepalceRemoveSpace.ReplaceAllString(str, " ")
	str = BadwordsReplacer.Replace(str)
	if !bSkip {
		str = ExtraBadwordsReplacer.Replace(str)
	}
	for _, keyw := range specialBadWords {
		str = FilterByKeyWord(str, keyw.Word, keyw.Reg)
	}
	return str
}

// 检测是否包含标点符号
func ExistPunctuation(str string) bool {
	r := regexp.MustCompile("[\\p{P}+~$`^=|<>～｀＄＾＋＝｜＜＞￥×]")
	rets := r.FindAllStringSubmatch(str, -1)
	if len(rets) > 0 {
		return true
	}
	return false
}

//判断字符串是否超长,中文字符长度算2，其余算1。
func CheckNameLen(name string, max int) bool {
	l := 0
	for _, v := range name {
		if unicode.Is(unicode.Scripts["Han"], v) {
			l += 2
		} else {
			l++
		}
	}
	if l > max {
		return true
	}
	return false
}

//判断是否包含预设之外的字符
func CheckValidCharacter(str string) (bool, []string) {
	invalids := []string{}
	f := true
	for _, s := range str {
		if !unicode.IsOneOf(validSet, s) {
			f = false
			invalids = append(invalids, string(s))
		}
	}
	return f, invalids
}

func UnicodeEmojiDecode(s string) string {
	src := regEmoji.FindAllString(s, -1)
	for i := 0; i < len(src); i++ {
		e := regEmojiExtract.ReplaceAllString(src[i], "")
		p, err := strconv.ParseInt(e, 16, 32)
		if err == nil {
			s = strings.Replace(s, src[i], string(rune(p)), -1)
		}
	}
	return s
}

func UnicodeEmojiCode(s string, clear ...bool) string {
	cflag := false
	if len(clear) > 0 {
		cflag = true
	}
	ret := ""
	rs := []rune(s)
	for i := 0; i < len(rs); i++ {
		if len(string(rs[i])) == 4 {
			if cflag {
				continue
			}
			u := `[\u` + strconv.FormatInt(int64(rs[i]), 16) + `]`
			ret += u
		} else {
			ret += string(rs[i])
		}
	}
	return ret
}

//-----------------------------------------

// 后台调用检测接口，返回一般用法
func AdminCheckOutput(str string, grp int, w http.ResponseWriter) {
	outputLine := func() {
		fmt.Fprintf(w, "====================================\n")
	}
	fmt.Fprintf(w, "origin:%v|\n", str)
	outputLine()
	fmt.Fprintf(w, "BadwordsCheck -- 普通敏感词检测:\n")
	fmt.Fprintf(w, "result:%v\n", BadwordsCheck(str))
	outputLine()
	fmt.Fprintf(w, "BadwordsCheck-skip ex -- 普通敏感词检测,跳过global的ex标点等:\n")
	fmt.Fprintf(w, "result:%v\n", BadwordsCheck(str, true))
	outputLine()
	fmt.Fprintf(w, "BadwordsCheckAdvanced-true -- db正则和敏感词检测,检测global的ex标点:\n")
	fmt.Fprintf(w, "result:%v\n", BadwordsCheckAdvanced(str, true))
	outputLine()
	fmt.Fprintf(w, "BadwordsCheckAdvanced-false -- db正则和敏感词检测,不检测global的ex标点:\n")
	fmt.Fprintf(w, "result:%v\n", BadwordsCheckAdvanced(str, false))
	outputLine()
	fmt.Fprintf(w, "BadwordsCheckAdvanced-false,grp:%d -- db正则和敏感词检测,不检测global的ex标点，对指定正则组进行检测:\n", grp)
	fmt.Fprintf(w, "result:%v\n", BadwordsCheckAdvanced(str, false, grp))
	outputLine()
	fmt.Fprintf(w, "BadwordsReplace -- 普通敏感词替换:\n")
	fmt.Fprintf(w, "result:%v\n", BadwordsReplace(str))
	outputLine()
	fmt.Fprintf(w, "BadwordsReplace skip ex -- 普通敏感词替换，跳过global的ex标点:\n")
	fmt.Fprintf(w, "result:%v\n", BadwordsReplace(str, true))
	outputLine()
	fmt.Fprintf(w, "CheckValidCharacter valid -- 检测是否在有效字符集中:\n")
	checkf, checks := CheckValidCharacter(str)
	fmt.Fprintf(w, "result:%v 非法的字符集:%v\n", checkf, checks)
}
