//go:generate go run semver/semver.go -f semver.tpl -o core.gen.go

package core

var (
	// SemVer is the semantic version.
	SemVer = "unknown"
)
