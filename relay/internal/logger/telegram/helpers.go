package telegram

import "strings"

const (
	maxMessageLength = 4096
	separator        = " "
)

// safeSplitText splits text into slices of maxLength by separator
func safeSplitText(text string, maxLength int, separator string) []string {
	tempText := text
	var slices []string

	for len(tempText) > 0 {

		if len(tempText) > maxLength {
			splitPosition := strings.LastIndex(tempText[:maxLength], separator)

			// if there is no separator in the last 1/4 (of maxLength) symbols then split at maxLength
			if splitPosition < maxLength*3/4 {
				splitPosition = maxLength
			}

			slices = append(slices, tempText[:splitPosition])

			tempText = tempText[splitPosition:]
			tempText = strings.TrimLeft(tempText, "\n\t\v\f\r ")

		} else {
			slices = append(slices, tempText)
			break
		}

	}
	return slices
}
