package middlewares

import (
	"regexp"
	"strings"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/logger"
)

var rmTagsRegexp = regexp.MustCompile(`<script>.*</script>|<.*?>`)
var rmEmptyLinesRegexp = regexp.MustCompile(`(?m)^\s+$`)
var rmUnnecessaryNewLinesRegexp = regexp.MustCompile(`(\r\n?|\n){2,}`)

var DefaultApplicableRegexps = []*ReplacableRegexp{
	{rmTagsRegexp, ""},
	{rmEmptyLinesRegexp, ""},
	{rmUnnecessaryNewLinesRegexp, ""},
}

type ReplacableRegexp struct {
	Regexp            *regexp.Regexp
	ReplacementString string
}

type RemoveHTMLTags struct {
	ApplicableRegexps []*ReplacableRegexp
	Hook              logger.Hook
}

func NewRemoveHTMLTags(hook logger.Hook, applicableRegexps []*ReplacableRegexp) *RemoveHTMLTags {
	return &RemoveHTMLTags{
		ApplicableRegexps: applicableRegexps,
		Hook:              hook,
	}
}

func (t *RemoveHTMLTags) Log(l *logger.ExtLog) {
	if strings.Contains(l.Error, "<html>") {
		l.Error = applyReplacableRegexpsToString(l.Error, t.ApplicableRegexps)
		l.Error = strings.TrimSpace(l.Error)
	}
	t.Hook.Log(l)
}

func applyReplacableRegexpsToString(src string, regexps []*ReplacableRegexp) string {
	for _, rg := range regexps {
		src = rg.Regexp.ReplaceAllString(src, rg.ReplacementString)
	}
	return src
}
