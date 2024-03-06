package document

type Parser interface {
	Parse(data string) (string, error)
}
