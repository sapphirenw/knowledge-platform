package document

type ParserTxt struct{}

func (p *ParserTxt) Parse(data string) (string, error) {
	return data, nil
}
