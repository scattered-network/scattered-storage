package validators

import (
	"regexp"

	"github.com/rs/zerolog/log"
)

// ValidateInput exists to output debug data on err.
func ValidateInput(check *regexp.Regexp, input string) bool {
	if check.MatchString(input) {
		return true
	}

	log.Error().Str("Expression", check.String()).Str("Input", input).Msg(ErrInvalidRegex.Error())

	return false
}

// ValidateRegex ensures that the regex expression is valid.
func ValidateRegex(expression string) *regexp.Regexp {
	if check, err := regexp.Compile(expression); err == nil {
		return check
	}

	log.Error().Str("Expression", expression).Msg(ErrInvalidRegex.Error())

	return nil
}
