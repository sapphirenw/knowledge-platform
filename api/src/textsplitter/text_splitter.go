package textsplitter

// CODE PULLED FROM: https://github.com/tmc/langchaingo/blob/main/textsplitter/
// Would install as a dependency, but do not need the entire Langchain code.

// TextSplitter is the standard interface for splitting texts.
type TextSplitter interface {
	SplitText(text string) ([]string, error)
}

const (
	_defaultTokenChunkSize    = 512
	_defaultTokenChunkOverlap = 100
)
