package fwheader

type Overridable interface {
	Override(string, []string)
}

type Expandable interface {
	Add(string, string)
}
