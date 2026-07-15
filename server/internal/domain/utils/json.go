package utils

import (
	"bytes"
	"encoding/json"
	"html/template"
	"regexp"
)

func FormatJSON(rawJSON string) string {
	var prettyJSON bytes.Buffer

	err := json.Indent(&prettyJSON, []byte(rawJSON), "", "  ")
	if err != nil {
		return rawJSON
	}

	return prettyJSON.String()
}

// It must be included to html template.
// .json-container .j-key { color: #f92672; font-weight: bold; }    /* Pink keys */
// .json-container .j-string { color: #e6db74; }                 /* Yellow keys */
// .json-container .j-value { color: #ae81ff; }                  /* Purple numbers/bool */
func PrettyColoredJson(jsonString string) template.HTML {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, []byte(jsonString), "", "  ")
	if err != nil {
		return template.HTML(jsonString)
	}

	jsonStr := prettyJSON.String()

	var keyRegex = regexp.MustCompile(`(".*?")\s*:`)
	jsonStr = keyRegex.ReplaceAllString(jsonStr, `<span class="j-key">$1</span>:`)

	var stringRegex = regexp.MustCompile(`(:\s*)(".*?")`)
	jsonStr = stringRegex.ReplaceAllString(jsonStr, `$1<span class="j-string">$2</span>`)

	var numberBoolRegex = regexp.MustCompile(`(:\s*)(true|false|null|-?\d+(?:\.\d+)?)`)
	jsonStr = numberBoolRegex.ReplaceAllString(jsonStr, `$1<span class="j-value">$2</span>`)

	return template.HTML(jsonStr)
}
