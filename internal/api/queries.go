package api

const ProblemsetPanelQuestionListQuery = `
query problemsetPanelQuestionList($searchKeyword: String, $limit: Int, $skip: Int) {
  problemsetPanelQuestionList(
    searchKeyword: $searchKeyword
    limit: $limit
    skip: $skip
  ) {
    questions {
      id
      titleSlug
      title
      questionFrontendId
      paidOnly
      difficulty
    }
    totalLength
    hasMore
  }
}
`

const QuestionDetailQuery = `
query questionDetail($titleSlug: String!) {
  question(titleSlug: $titleSlug) {
    title
    titleSlug
    questionId
    questionFrontendId
    codeSnippets {
      code
      lang
      langSlug
    }
  }
}
`
