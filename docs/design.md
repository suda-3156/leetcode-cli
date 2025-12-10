# LeetCode CLI 設計書

## 概要

キーワードで LeetCode の問題を検索し、選択した問題のコードスニペットを含むファイルを生成する CLI ツール。

## 技術スタック

| 項目               | 選定                                                     |
| ------------------ | -------------------------------------------------------- |
| 言語               | Go                                                       |
| CLI フレームワーク | [Cobra](https://github.com/spf13/cobra)                  |
| TUI                | [Bubble Tea](https://github.com/charmbracelet/bubbletea) |
| HTTP クライアント  | 標準ライブラリ `net/http`                                |
| GraphQL            | 手動でリクエスト構築（introspection が無効のため）       |

---

## ディレクトリ構成

```
leetcode-cli/
├── main.go                 # エントリーポイント
├── go.mod
├── go.sum
├── cmd/
│   └── root.go             # ルートコマンド定義
├── internal/
│   ├── api/
│   │   ├── client.go       # GraphQL クライアント
│   │   ├── queries.go      # GraphQL クエリ定義
│   │   └── types.go        # API レスポンス型定義
│   ├── tui/
│   │   ├── model.go        # Bubble Tea モデル
│   │   ├── question_list.go # 問題リスト選択画面
│   │   ├── language_select.go # 言語選択画面
│   │   └── path_input.go   # パス入力画面
│   ├── generator/
│   │   ├── file.go         # ファイル生成ロジック
│   │   └── template.go     # 言語別コメントテンプレート
│   └── config/
│       └── config.go       # 設定・定数管理
└── temp/                   # 開発用ドキュメント
```

---

## コマンド設計

### 基本使用法

```sh
leetcode-cli [keyword] [flags]
```

`--slug`オプション指定次は`[keyword]`は不要(無視される)

### フラグ

| フラグ   | 短縮 | 型     | 説明                                               |
| -------- | ---- | ------ | -------------------------------------------------- |
| `--slug` | -    | string | 検索をスキップし、titleSlug を直接指定             |
| `--lang` | `-l` | string | 言語を指定（langSlug: `golang`, `python3` など）   |
| `--path` | `-p` | string | 出力ファイルパス（`default` でデフォルトパス使用） |

### 動作フロー

```
┌─────────────────────────────────────────────────────────────┐
│ 1. キーワード引数 or --slug フラグ                           │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│ 2. problemsetPanelQuestionList API で問題検索               │
│    （--slug 指定時はスキップ）                               │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│ 3. TUI: 問題リストから選択                                   │
│    表示形式: "150. Evaluate Reverse Polish Notation (MEDIUM)"│
│    （--slug 指定時はスキップ）                               │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│ 4. questionDetail API で詳細取得                            │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│ 5. TUI: 言語選択                                            │
│    codeSnippets から利用可能な言語を表示                      │
│    （--lang 指定時はスキップ、無効な言語はエラー）            │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│ 6. TUI: 出力パス入力                                        │
│    デフォルト: ./{yyyy-mm-dd}/{num}.{slug}.{ext}            │
│    （--path 指定時はスキップ）                               │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│ 7. ファイル生成                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## 型定義

### API レスポンス型

```go
// internal/api/types.go

// ProblemsetResponse - 問題リスト検索レスポンス
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

// QuestionDetailResponse - 問題詳細レスポンス
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
```

---

## 生成ファイルフォーマット

### デフォルトパス

```
./{yyyy-mm-dd}/{questionFrontendId}.{titleSlug}.{extension}
```

例: `./2025-12-10/150.evaluate-reverse-polish-notation.go`

### 言語別拡張子マッピング

| langSlug         | 拡張子 | コメント記法 |
| ---------------- | ------ | ------------ |
| golang           | .go    | `//`         |
| python3 / python | .py    | `#`          |
| javascript       | .js    | `//`         |
| typescript       | .ts    | `//`         |
| java             | .java  | `//`         |
| cpp              | .cpp   | `//`         |
| c                | .c     | `//`         |
| csharp           | .cs    | `//`         |
| rust             | .rs    | `//`         |
| ruby             | .rb    | `#`          |
| swift            | .swift | `//`         |
| kotlin           | .kt    | `//`         |
| scala            | .scala | `//`         |
| php              | .php   | `//`         |
| dart             | .dart  | `//`         |
| sql              | .sql   | `--`         |

### ファイルテンプレート

```
{comment} {yyyy-mm-dd}
{comment} {questionFrontendId}. {title}
{comment} {language name}

{code snippet}
```

例（Go）:

```go
// 2025-12-10
// 150. Evaluate Reverse Polish Notation
// Go

func evalRPN(tokens []string) int {

}
```

---

## API クライアント設計

### エンドポイント

```
POST https://leetcode.com/graphql/
Content-Type: application/json
```

### リクエスト構造

```go
type GraphQLRequest struct {
    Query     string                 `json:"query"`
    Variables map[string]interface{} `json:"variables"`
}
```

### クエリ定義

`internal/api/queries.go` に GraphQL クエリを文字列定数として定義：

```go
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
    questionTitle
    codeSnippets {
      code
      lang
      langSlug
    }
  }
}
`
```

---

## TUI 設計 (Bubble Tea)

### 画面遷移

```
StateSearching → StateQuestionList → StateLanguageSelect → StatePathInput → StateDone
```

### モデル構造

```go
type Model struct {
    state          State
    keyword        string
    questions      []Question
    selectedQ      *Question
    questionDetail *QuestionDetail
    languages      []CodeSnippet
    selectedLang   *CodeSnippet
    outputPath     string
    cursor         int
    textInput      textinput.Model
    err            error
}
```

### キーバインド

| キー  | アクション     |
| ----- | -------------- |
| ↑/k   | カーソル上移動 |
| ↓/j   | カーソル下移動 |
| Enter | 選択確定       |
| q/Esc | 終了           |

---

## エラーハンドリング

| エラー種別                | 対応                                          |
| ------------------------- | --------------------------------------------- |
| ネットワークエラー        | エラーメッセージ表示して終了                  |
| 問題が見つからない        | "No questions found" 表示                     |
| 無効な言語指定 (`--lang`) | 利用可能な言語一覧を表示してエラー終了        |
| paidOnly 問題             | リストに `[Premium]` ラベル付与（選択は可能） |
| ファイル書き込みエラー    | エラーメッセージ表示して終了                  |

---

## 今後の拡張 (Future)

- [ ] 言語別の型定義テンプレート自動生成
- [ ] 設定ファイル対応 (`~/.leetcode-cli.yaml`)
