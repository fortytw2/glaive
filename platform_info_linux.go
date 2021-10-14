package glaive

import (
	"fmt"
	"strings"

	"github.com/fortytw2/lounge"
	"golang.org/x/sys/unix"
)

func printPlatformInformation(log lounge.Log) {
	utsname := unix.Utsname{}
	unix.Uname(&utsname)

	formattedString := fmt.Sprintf("linux uname: %s %s %s", string(utsname.Sysname[:]), string(utsname.Release[:]), string(utsname.Machine[:]))

	formattedString = strings.ReplaceAll(formattedString, "\u0000", "")

	log.Infof(formattedString)
}
