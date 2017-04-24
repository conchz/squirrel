/*
	Package boxes is a wrapper around a set of go.rice boxes created by a go.rice Config
	defined to first look for files in the package directory, and then in the binary.
*/
package boxes

import rice "github.com/GeertJohan/go.rice"

//go:generate rice embed-go

// There are some odd things in this package due to how go.rice's rice tool
// works.  In order to generate the appropriate go files, the tool must go
// through this package and find all calls to FindBox (and MustFindBox) to
// figure out which directories need to be put into "Boxes".

// go.rice idiosyncrasy #1 - The package name variable must shadow the Config
// (if a rice.Config is used) in order for the rice tool to pick up calls to FindBox.

// go.rice idiosyncrasy #2 - Calls to FindBox must be done with string literals.

var (
	favicon   *rice.Box
	images    *rice.Box
	css       *rice.Box
	js        *rice.Box
	templates *rice.Box

	_rice = rice.Config{
		LocateOrder: []rice.LocateMethod{
			rice.LocateFS,
			rice.LocateWorkingDirectory,
			rice.LocateEmbedded,
			rice.LocateAppended,
		},
	}
)

func Favicon() *rice.Box {
	if favicon != nil {
		return favicon
	}

	favicon, err := _rice.FindBox("../static")
	if err != nil {
		panic(err)
	}

	return favicon
}

func Images() *rice.Box {
	if images != nil {
		return images
	}

	images, err := _rice.FindBox("../static/img")
	if err != nil {
		panic(err)
	}

	return images
}

func CSS() *rice.Box {
	if css != nil {
		return css
	}

	css, err := _rice.FindBox("../static/css")
	if err != nil {
		panic(err)
	}

	return css
}

func JS() *rice.Box {
	if js != nil {
		return js
	}

	js, err := _rice.FindBox("../static/js")
	if err != nil {
		panic(err)
	}

	return js
}

func Templates() *rice.Box {
	if templates != nil {
		return templates
	}

	templates, err := _rice.FindBox("../static/views")
	if err != nil {
		panic(err)
	}

	return templates
}
