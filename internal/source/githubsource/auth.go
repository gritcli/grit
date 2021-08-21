package githubsource

// AppClientID is the Client ID for the "Grit CLI" GitHub app on github.com.
const AppClientID = "976b90cbef967ca64b7e"

// RequiredScopes is the set of scopes required to support all Grit
// functionality.
var RequiredScopes = []string{}

// diffScopes returns the scopes that are in a but not in b.
func diffScopes(a, b []string) []string {
	var diff []string

outer:
	for _, sa := range a {
		for _, sb := range b {
			if sa == sb {
				continue outer
			}
		}

		diff = append(diff, sa)
	}

	return diff
}
