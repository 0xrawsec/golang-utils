package check

import "regexp"

var (
	MD5ValRegexp    = regexp.MustCompile("^(?i)([0-9a-h]{32})$")
	SHA1ValRegexp   = regexp.MustCompile("^(?i)([0-9a-h]{40})$")
	SHA256ValRegexp = regexp.MustCompile("^(?i)([0-9a-h]{64})$")
	SHA512ValRegexp = regexp.MustCompile("^(?i)([0-9a-h]{128})$")
)

func IsMD5(hash string) bool {
	return MD5ValRegexp.MatchString(hash)
}

func IsSHA1(hash string) bool {
	return SHA1ValRegexp.MatchString(hash)
}

func IsSHA256(hash string) bool {
	return SHA256ValRegexp.MatchString(hash)
}

func IsSHA512(hash string) bool {
	return SHA512ValRegexp.MatchString(hash)
}

func IsValidHash(hash string) bool {
	switch {
	case IsMD5(hash), IsSHA1(hash), IsSHA256(hash), IsSHA512(hash):
		return true
	default:
		return false
	}
}
