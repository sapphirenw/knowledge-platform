package project

// for holding the list that is generated from the model
type projectIdeas struct {
	Ideas []*ProjectIdea `json:"ideas"`
}

type ProjectIdea struct {
	Title string `json:"title"`
}
