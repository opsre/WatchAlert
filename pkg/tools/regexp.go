package tools

import (
	"fmt"
	"regexp"
	"sync"
)

type regexpCacheEntry struct {
	expression *regexp.Regexp
	err        error
}

var regexpCache sync.Map

// CompileRegexp compiles and caches a regular expression. Compilation errors
// are cached as well so an invalid configuration cannot cause repeated work.
func CompileRegexp(pattern string) (*regexp.Regexp, error) {
	if cached, ok := regexpCache.Load(pattern); ok {
		entry := cached.(regexpCacheEntry)
		return entry.expression, entry.err
	}

	expression, err := regexp.Compile(pattern)
	entry := regexpCacheEntry{expression: expression, err: err}
	actual, _ := regexpCache.LoadOrStore(pattern, entry)
	cached := actual.(regexpCacheEntry)
	return cached.expression, cached.err
}

// MatchRegexp matches value against a cached regular expression. An invalid
// expression never matches and is returned to the caller for validation or
// diagnostics.
func MatchRegexp(pattern, value string) (bool, error) {
	expression, err := CompileRegexp(pattern)
	if err != nil {
		return false, fmt.Errorf("invalid regular expression %q: %w", pattern, err)
	}
	return expression.MatchString(value), nil
}

// ValidateRegexp validates a regular expression without exposing the cache
// implementation to configuration services.
func ValidateRegexp(pattern string) error {
	_, err := CompileRegexp(pattern)
	return err
}
