package valid

import "regexp"

var userNameRegexp = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]*[._-]?[a-zA-Z0-9]+$`)

// Username validates username
func Username(username string) bool {
	return userNameRegexp.MatchString(username) && len(username) <= 20 && len(username) > 2
}
