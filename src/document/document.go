package document

/*
Convenience powerful struct that parses raw data to provide consistency in how the
input data is parsed, and what some metadata about the file is.
*/
type Doc struct {
	Filename string
	Filetype Filetype
	Data     string
}

// Get the size this document takes up on disk in bytes
func (d *Doc) GetSizeInBytes() int {
	return len([]byte(d.Data))
}
