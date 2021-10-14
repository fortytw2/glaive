package buildinfo

var (
	DisplayVersion = "+development"
	GitVersion     = "GIT_VERSION_GOES_HERE"
)

const GlaiveVersion = "0.0.1"

func IsRelease() bool {
	return DisplayVersion != "+development"
}

var DevCommand string
var DevPort string
