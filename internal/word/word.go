package word

import (
	"regexp"
	"strings"
	"unicode"
)

func ToUpper(s string) string {
	return strings.ToUpper(s)
}

func ToLower(s string) string {
	return strings.ToLower(s)
}

func UnderscoreToUpperCamelCase(s string) string {
	s = strings.Replace(s, "_", " ", -1)
	s = strings.Title(s)
	return strings.Replace(s, " ", "", -1)
}

func UnderscoreToLowerCamelCase(s string) string {
	s = UnderscoreToUpperCamelCase(s)
	return string(unicode.ToLower(rune(s[0]))) + s[1:]
}

func CenterStr25(str string) string {
	return CenterStr(str, 25)

}

func CenterStr40(str string) string {
	return CenterStr(str, 40)

}

func CenterStr90(str string) string {
	return CenterStr(str, 80)

}

func CenterStr(str string, lenS int) string {
	var hzRegexp = regexp.MustCompile("^[\u4e00-\u9fa5]$")
	// s := fmt.Sprintf("%25s", str)
	var totalZW int
	for _, v := range str {

		//如果是中文
		if hzRegexp.MatchString(string(v)) {
			totalZW += 1
		}
	}

	spaceLen := lenS - len(str) + totalZW*2
	var sResult = "嘿"

	for i := 0; i < spaceLen/2; i++ {
		sResult += " "
	}
	sResult += str
	for i := 0; i < spaceLen/2; i++ {
		sResult += " "
	}
	sResult += "嘿"
	return sResult

}

func CenterStr18(str string) string {
	return CenterStr(str, 18)

}

func DescTow(s int) int {

	return s - 2
}

func DescThree(s int) int {

	return s - 3
}

func DescFour(s int) int {

	return s - 4
}

func CamelCaseToUnderscore(s string) string {
	var output []rune
	for i, r := range s {
		if i == 0 {
			output = append(output, unicode.ToLower(r))
			continue
		}
		if unicode.IsUpper(r) {
			output = append(output, '_')
		}
		output = append(output, unicode.ToLower(r))
	}
	return string(output)

}
