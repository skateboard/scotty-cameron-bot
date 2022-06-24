package module

import (
	"encoding/base64"
	"regexp"
	"strings"
	"testing"
)

func TestPixel(t *testing.T) {
	url := "https://www.scottycameron.com/akam/13/pixel_6ed0dae?a=dD1lMTgzZDg0Yzc5ZDdhZmNkZDUxNzEyOTRkZTAzNjUxMTljMTMwZDg4JmpzPW9mZg=="
	//String str2 = str.split("\\?")[1];
	//        String substring = str2.substring(str2.indexOf("=") + 1);
	//        if (substring.contains("&")) {
	//            substring = substring.split("&")[0];
	//        }
	//        Matcher matcher = T_PATTERN.matcher(new String(Base64.getDecoder().decode(substring)));
	//        if (matcher.find()) {
	//            return matcher.group(1);
	//        }
	str2 := strings.Split(url, "?")[1]
	substring := str2[strings.Index(str2, "=")+1:]
	if strings.Contains(substring, "&") {
		substring = substring[:strings.Index(substring, "&")]
	}
	b, _ := base64.StdEncoding.DecodeString(substring)
	re := regexp.MustCompile("t=([0-z]*)")
	matcher := re.FindStringSubmatch(string(b))
	if len(matcher) == 0 {
		t.Error("Error Parsing Pixel")
		return
	}
	t.Log(matcher[1])
}
