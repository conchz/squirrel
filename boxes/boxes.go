/*
	Package boxes is a wrapper around a set of go.rice boxes created by a go.rice Config
	defined to first look for files in the package directory, and then in the binary.
*/
package boxes

import rice "github.com/GeertJohan/go.rice"

//go:generate rice clean
//go:generate rice embed-go

// reference: https://github.com/yext/revere/blob/master/boxes/boxes.go

// There are some odd things in this package due to how go.rice's rice tool
// works.  In order to generate the appropriate go files, the tool must go
// through this package and find all calls to FindBox (and MustFindBox) to
// figure out which directories need to be put into "Boxes".

// go.rice idiosyncrasy #1 - The package name variable must shadow the Config
// (if a rice.Config is used) in order for the rice tool to pick up calls to FindBox.

// go.rice idiosyncrasy #2 - Calls to FindBox must be done with string literals.

var assets *rice.Box

func Assets() *rice.Box {
	if assets != nil {
		return assets
	}

	assets := rice.MustFindBox("../static")

	return assets
}
