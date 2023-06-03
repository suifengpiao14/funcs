package camelsnake

import (
	"bytes"
	"strings"
	"unicode"

	"github.com/suifengpiao14/camelsnake/expr"
)

// Casing exceptions
var ToLower = map[string]string{"OAuth": "oauth"}

// CamelCase produces the CamelCase version of the given string. It removes any
// non letter and non digit character.
//
// If firstUpper is true the first letter of the string is capitalized else
// the first letter is in lowercase.
//
// If acronym is true and a part of the string is a common acronym
// then it keeps the part capitalized (firstUpper = true)
// (e.g. APIVersion) or lowercase (firstUpper = false) (e.g. apiVersion).
func CamelCase(name string, firstUpper bool, acronym bool) string {
	if name == "" {
		return ""
	}

	runes := []rune(name)
	// remove trailing invalid identifiers (makes code below simpler)
	runes = removeTrailingInvalid(runes)

	// all characters are invalid
	if len(runes) == 0 {
		return ""
	}

	w, i := 0, 0 // index of start of word, scan
	for i+1 <= len(runes) {
		eow := false // whether we hit the end of a word

		// remove leading invalid identifiers
		runes = removeInvalidAtIndex(i, runes)

		if i+1 == len(runes) {
			eow = true
		} else if !validIdentifier(runes[i]) {
			// get rid of it
			runes = append(runes[:i], runes[i+1:]...)
		} else if runes[i+1] == '_' {
			// underscore; shift the remainder forward over any run of underscores
			eow = true
			n := 1
			for i+n+1 < len(runes) && runes[i+n+1] == '_' {
				n++
			}
			copy(runes[i+1:], runes[i+n+1:])
			runes = runes[:len(runes)-n]
		} else if isLower(runes[i]) && !isLower(runes[i+1]) {
			// lower->non-lower
			eow = true
		}
		i++
		if !eow {
			continue
		}

		// [w,i] is a word.
		word := string(runes[w:i])
		// is it one of our initialisms?
		if u := strings.ToUpper(word); CommonInitialisms[u] {
			switch {
			case firstUpper && acronym:
				// u is already in upper case. Nothing to do here.
			case firstUpper && !acronym:
				u = expr.Title(strings.ToLower(u))
			case w > 0 && !acronym:
				u = expr.Title(strings.ToLower(u))
			case w == 0:
				u = strings.ToLower(u)
			}

			// All the common initialisms are ASCII,
			// so we can replace the bytes exactly.
			copy(runes[w:], []rune(u))
		} else if w > 0 && strings.ToLower(word) == word {
			// already all lowercase, and not the first word, so uppercase the first character.
			runes[w] = unicode.ToUpper(runes[w])
		} else if w == 0 && strings.ToLower(word) == word && firstUpper {
			runes[w] = unicode.ToUpper(runes[w])
		}
		if w == 0 && !firstUpper {
			runes[w] = unicode.ToLower(runes[w])
		}
		//advance to next word
		w = i
	}

	return string(runes)
}

// SnakeCase produces the snake_case version of the given CamelCase string.
// News    => news
// OldNews => old_news
// CNNNews => cnn_news
func SnakeCase(name string) string {
	// Special handling for single "words" starting with multiple upper case letters
	for u, l := range ToLower {
		name = strings.Replace(name, u, l, -1)
	}

	// Remove leading and trailing blank spaces and replace any blank spaces in
	// between with a single underscore
	name = strings.Join(strings.Fields(name), "_")

	// Special handling for dashes to convert them into underscores
	name = strings.Replace(name, "-", "_", -1)

	var b bytes.Buffer
	ln := len(name)
	if ln == 0 {
		return ""
	}
	n := rune(name[0])
	b.WriteRune(unicode.ToLower(n))
	lastLower, isLower, lastUnder, isUnder := false, true, false, false
	for i := 1; i < ln; i++ {
		r := rune(name[i])
		isLower = unicode.IsLower(r) && unicode.IsLetter(r) || unicode.IsDigit(r)
		isUnder = r == '_'
		if !isLower && !isUnder {
			if lastLower && !lastUnder {
				b.WriteRune('_')
			} else if ln > i+1 {
				rn := rune(name[i+1])
				if unicode.IsLower(rn) && rn != '_' && !lastUnder {
					b.WriteRune('_')
				}
			}
		}
		b.WriteRune(unicode.ToLower(r))
		lastLower = isLower
		lastUnder = isUnder
	}
	return b.String()
}

// KebabCase produces the kebab-case version of the given CamelCase string.
func KebabCase(name string) string {
	name = SnakeCase(name)
	ln := len(name)
	if name[ln-1] == '_' {
		name = name[:ln-1]
	}
	return strings.Replace(name, "_", "-", -1)
}

// isLower returns true if the character is considered a lower case character
// when transforming word into CamelCase.
func isLower(r rune) bool {
	return unicode.IsDigit(r) || unicode.IsLower(r)
}

// validIdentifier returns true if the rune is a letter or number
func validIdentifier(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

// removeTrailingInvalid removes trailing invalid identifiers from runes.
func removeTrailingInvalid(runes []rune) []rune {
	valid := len(runes) - 1
	for ; valid >= 0 && !validIdentifier(runes[valid]); valid-- {
	}

	return runes[0 : valid+1]
}

// removeInvalidAtIndex removes consecutive invalid identifiers from runes starting at index i.
func removeInvalidAtIndex(i int, runes []rune) []rune {
	valid := i
	for ; valid < len(runes) && !validIdentifier(runes[valid]); valid++ {
	}

	return append(runes[:i], runes[valid:]...)
}

var (
	// common words who need to keep their
	CommonInitialisms = map[string]bool{
		// "API":   true,
		// "ASCII": true,
		// "CPU":   true,
		// "CSS":   true,
		// "DNS":   true,
		// "EOF":   true,
		// "GUID":  true,
		// "HTML":  true,
		// "HTTP":  true,
		// "HTTPS": true,
		// "ID":    true,
		// "IP":    true,
		// "JMES":  true,
		// "JSON":  true,
		// "JWT":   true,
		// "LHS":   true,
		// "OK":    true,
		// "QPS":   true,
		// "RAM":   true,
		// "RHS":   true,
		// "RPC":   true,
		// "SLA":   true,
		// "SMTP":  true,
		// "SQL":   true,
		// "SSH":   true,
		// "TCP":   true,
		// "TLS":   true,
		// "TTL":   true,
		// "UDP":   true,
		// "UI":    true,
		// "UID":   true,
		// "UUID":  true,
		// "URI":   true,
		// "URL":   true,
		// "UTF8":  true,
		// "VM":    true,
		// "XML":   true,
		// "XSRF":  true,
		// "XSS":   true,
	}
)
