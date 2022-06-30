package gpmf

func parseMetadata(e *Element) error {
	e.parent.Metadata[e.friendlyName()] = e.Data

	return nil
}

func parseHasMetadata(e *Element) error {
	e.metadata()

	return nil
}
