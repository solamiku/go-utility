/*
	屏蔽词测试
*/
package badwords

import (
	"regexp"
	"testing"
	"unicode/utf8"
)

func Benchmark_match(b *testing.B) {
	// str := "收人,关注微LｇL信公众号Lｇdanｇ领[星]2231[环]0101[互]0101[动]取传说装备套装Lgdang..测试L g d α n ｇ,,收人,关注微LｇL信公众号Lｇdanｇ领取传说装备套装Lgdang..测试L g d α n ｇ,,收人,关注微LｇL信公众号Lｇdanｇ领取传说装备套装Lgdang..测试L g d α n ｇ"
	str := "收人,关注微LｇL信公众号Lｇdanｇ领2[星]23101010101取传说装备套装Lgdang..测试L g d α n ｇ,,收人,关注微LｇL信公众号Lｇdanｇ领取传说装备套装Lgdang..测试L g d α n ｇ,,收人,关注微LｇL信公众号Lｇdanｇ领取传说装备套装Lgdang..测试L g d α n ｇ"
	// reg1 := regexp.MustCompile("星.*环.*互.*动")
	reg1 := regexp.MustCompile("星.{0,4}环.{0,4}互.{0,4}动")
	for i := 0; i < b.N; i++ {
		reg1.MatchString(str)
	}
}

func Test_CheckValidCharacter(t *testing.T) {
	check := func(str string) {
		f, ret := CheckValidCharacter(str)
		t.Log(f, utf8.RuneCountInString(str), len(ret), ret)
	}
	check("abcdefghijklmnopqrstuvwxyz")
	check("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	check("ａｂｃＤＥＦＧｈｉｊｋＬＭＮｏＰＱＲｓｔＵＶＷＸＹｚ")
	check("1234567890")
	check("１２３４５６７８９𝟎")
	check("+_-><=)(*&%#@!~|.,。，\"?/---")
	check("我是简体中文-开门，我是繁体中文-開門")
	check("϶؆⁒∁$×÷")
	check("❶①⑪⒈⑴⑾αΑνΝалЧāūēㄅㄑㄧ㉠ㄱあさアサマ㊤㈱╱╲●☼㏠㍙⒜⒨ヾω`")
}

func Test_removeReplace(t *testing.T) {
	t.Log(regRepalceRemoveSpace.ReplaceAllString(`你好呀
	请问你是 jack and 
	   sandeee 
	  ace    aa      吗`, " "))
}

func Test_adfilter(t *testing.T) {
	t.Log(processAlphabet("L g d α n ｇ"))
	t.Log(processAlphabet("abcABCａｂｃＡＢＣαβγΑΒΓ"))
	t.Log(processAlphabet("ΧHSHＥｑｕ"))
	t.Log(FilterByKeyWord(
		"微LｇL信Lｇdanｇ装Lgdang..测ΧHshequ,,收",
		"xhshequ",
	))
	t.Log(FilterByKeyWord(
		"微LｇL信LｇdAAnｇ装Lgdang..L g d α n ｇ,,收",
		"lgdang",
	))
	t.Log(FilterByKeyWord(
		"微LｇL信搜索·梦幻神奇·测L g d α 快 领 包 n ｇ,,收",
		"快领包",
	))
	t.Log(FilterByKeyWord(
		"测 Ζhoushen...",
		"zhoushen",
	))
	t.Log(FilterByKeyWord(
		"测 γΥf...",
		"yyf",
	))

}

func Test_badwordsCheck(t *testing.T) {
	err := InitBadwords("../data/badwords.txt", nil, []string{"。", "\""})
	if err != nil {
		t.Fatal(err)
	}
	tests := []string{
		"习大大",
		"几把",
		"鸡巴",
		"瘪三",
		"好的吧",
		"共青团",
	}
	for _, v := range tests {
		t.Logf("%s check result:%v", v, BadwordsCheck(v))
	}
	replace := map[string]string{
		"习大大你好呀":     "***你好呀",
		"好的吧":        "好的吧",
		`你好`:         "你好",
		"-_-\"":      "-_-*",
		"--_--。":     "--_--*",
		"一二三四五六七八九十": "一二三四五六七**十",
	}
	for k, v := range replace {
		if a := BadwordsReplace(k); a != v {
			t.Errorf("%s check failed. need:%s get:%s", k, v, a)
		}
	}
	t.Logf("test extra badwords not skip %v", BadwordsReplace("啊哈。."))
	t.Logf("test extra badwords skip %v", BadwordsReplace("啊哈。.", true))

	// if err := ddb.Load_db("../../data/design.db"); err == nil {
	// 	names := ddb.Table("randomName").RowsAll()
	// 	for _, v := range names {
	// 		name := v.Gets("text")
	// 		if BadwordsCheck(name) {
	// 			t.Logf("check name %s invalid.", name)
	// 		}
	// 	}
	// } else {
	// 	t.Logf("init ddb err:%v", err)
	// }

}

func Benchmark_adfilter(b *testing.B) {
	str := "收人,关注微LｇL信公众号Lｇdanｇ领取传说装备套装Lgdang..测试L g d α n ｇ,,收人,关注微LｇL信公众号Lｇdanｇ领取传说装备套装Lgdang..测试L g d α n ｇ,,收人,关注微LｇL信公众号Lｇdanｇ领取传说装备套装Lgdang..测试L g d α n ｇ"
	b.Log(FilterByKeyWord(str, "lgdang"))
	for i := 0; i < b.N; i++ {
		FilterByKeyWord(str, "lgdang")
	}
}
