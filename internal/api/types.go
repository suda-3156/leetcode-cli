package api

type ProblemsetResponse struct {
	Data struct {
		ProblemsetPanelQuestionList struct {
			Questions   []Question `json:"questions"`
			TotalLength int        `json:"totalLength"`
			HasMore     bool       `json:"hasMore"`
		} `json:"problemsetPanelQuestionList"`
	} `json:"data"`
}

type Question struct {
	ID                 int        `json:"id"`
	TitleSlug          string     `json:"titleSlug"`
	Title              string     `json:"title"`
	QuestionFrontendID string     `json:"questionFrontendId"`
	Difficulty         string     `json:"difficulty"`
	PaidOnly           bool       `json:"paidOnly"`
	TopicTags          []TopicTag `json:"topicTags"`
}

type TopicTag struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type QuestionDetailResponse struct {
	Data struct {
		Question QuestionDetail `json:"question"`
	} `json:"data"`
}

type QuestionDetail struct {
	Title              string        `json:"title"`
	TitleSlug          string        `json:"titleSlug"`
	QuestionID         string        `json:"questionId"`
	QuestionFrontendID string        `json:"questionFrontendId"`
	CodeSnippets       []CodeSnippet `json:"codeSnippets"`
}

type CodeSnippet struct {
	Code     string `json:"code"`
	Lang     string `json:"lang"`
	LangSlug string `json:"langSlug"`
}
