package core

type NavigationItems []*NavigationItem
type NavigationItem struct {
	Label  string
	Link   string
	Active bool
}
