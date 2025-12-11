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
│   ├── root.go             # ルートコマンド定義
│   └── version.go          # バージョンコマンド
├── internal/
│   ├── api/
│   │   ├── client.go       # GraphQL クライアント
│   │   ├── queries.go      # GraphQL クエリ定義
│   │   └── types.go        # API レスポンス型定義
│   ├── tui/
│   │   ├── init.go         # TUI エントリーポイント
│   │   ├── model.go        # TUI メインモデル（phase orchestrator）
│   │   ├── update.go       # メッセージハンドリング
│   │   ├── view.go         # ビューレンダリング
│   │   ├── model/
│   │   │   └── model.go    # 共有データモデル
│   │   ├── phases/
│   │   │   ├── phase.go    # PhaseHandler インターフェース定義
│   │   │   ├── initial.go  # 設定読み込みフェーズ
│   │   │   ├── decide_question.go  # 問題選択フェーズ
│   │   │   ├── decide_language.go  # 言語選択フェーズ
│   │   │   ├── decide_path.go      # パス入力フェーズ
│   │   │   ├── overwrite_confirm.go # 上書き確認フェーズ
│   │   │   ├── generation.go       # ファイル生成フェーズ
│   │   │   └── done.go             # 完了フェーズ
│   │   └── styles/
│   │       └── styles.go   # スタイル定義
│   ├── generator/
│   │   ├── file.go         # ファイル生成ロジック
│   │   ├── content.go      # コンテンツ生成
│   │   ├── embed.go        # 埋め込みテンプレート
│   │   └── template/
│   │       ├── default.tmpl  # デフォルトテンプレート
│   │       ├── golang.tmpl   # Go用テンプレート
│   │       └── python.tmpl   # Python用テンプレート
│   ├── config/
│   │   ├── config.go       # 設定管理
│   │   ├── default.go      # デフォルト値
│   │   └── lang.go         # 言語設定マッピング
│   └── file/
│       ├── backup.go       # バックアップ機能
│       ├── exist.go        # ファイル存在確認
│       ├── generate.go     # ファイル生成
│       ├── read.go         # ファイル読み込み
│       └── yaml.go         # YAML パース
└── docs/
    └── design.md           # 本ドキュメント
```

---

## コマンド設計

### 基本使用法

```sh
leetcode-cli [keyword] [flags]
```

`--slug`オプション指定時は`[keyword]`は不要（無視される）

### フラグ

| フラグ        | 短縮 | 型     | 説明                                                                    |
| ------------- | ---- | ------ | ----------------------------------------------------------------------- |
| `--slug`      | -    | string | 検索をスキップし、titleSlug を直接指定                                  |
| `--lang`      | `-l` | string | 言語を指定（langSlug: `golang`, `python3` など）                        |
| `--path`      | `-p` | string | 出力ファイルパス（テンプレート対応）                                    |
| `--overwrite` | -    | string | 上書き動作を指定（`always`, `backup`, `never`, 省略時はプロンプト表示） |
| `--config`    | `-c` | string | 設定ファイルパス（デフォルト: `.leetcode-cli.yaml`）                    |

### 動作フロー

```
┌─────────────────────────────────────────────────────────────┐
│ 1. InitialPhase: 設定ファイル読み込み                        │
│    - CLI フラグ > 設定ファイル > デフォルト値の優先順位      │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│ 2. DecideQuestionPhase: 問題選択                             │
│    - --slug 指定時: このフェーズをスキップ                   │
│    - キーワード検索: problemsetPanelQuestionList API         │
│    - 表示: "150. Evaluate Reverse Polish Notation (MEDIUM)" │
│    - Premium問題は選択不可                                   │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│ 3. DecideLanguagePhase: 言語選択                             │
│    - questionDetail API で問題詳細を取得                    │
│    - --lang 指定時: 自動選択（無効な場合は警告を表示）       │
│    - codeSnippets から利用可能な言語を表示                    │
│    - h/left キーで前のフェーズに戻る                         │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│ 4. DecidePathPhase: 出力パス入力                             │
│    - デフォルトパス生成: ./{yyyy-mm-dd}/{num}.{slug}.{ext}  │
│    - --path 指定時: 自動で次のフェーズへ                     │
│    - テキスト入力でパス編集可能                              │
│    - h/left キーで前のフェーズに戻る                         │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│ 5. OverwriteConfirmPhase: 上書き確認（ファイル存在時のみ）   │
│    - --overwrite 指定時: 自動処理                            │
│      - always: 上書き                                        │
│      - backup: バックアップ作成後上書き                      │
│      - never: エラー終了                                     │
│    - 未指定時: ユーザーに選択を促す                          │
│      - o: 上書き                                             │
│      - b: バックアップ作成後上書き                           │
│      - r: パス入力に戻る                                     │
│      - q/esc: 終了                                           │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│ 6. GenerationPhase: ファイル生成                             │
│    - バックアップ処理（選択時）                              │
│    - ファイル生成（テンプレート適用）                        │
│    - ディレクトリ自動作成                                    │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│ 7. DonePhase: 完了                                           │
│    - 生成ファイルパスを表示                                  │
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
```

## TUI 設計 (Bubble Tea)

### アーキテクチャ: Phase Handler パターン

TUI は **Phase Handler パターン** を採用し、各フェーズを独立したハンドラーとして実装しています。

#### PhaseHandler インターフェース

```go
type PhaseHandler interface {
    // Enter: フェーズに遷移する際に呼び出される
    Enter(m *model.Model) tea.Cmd

    // Update: メッセージを処理し、次のフェーズを返す
    Update(m *model.Model, msg tea.Msg) (tea.Cmd, *PhaseType)

    // View: 現在のフェーズの画面を描画
    View(m *model.Model) string
}
```

#### フェーズ一覧

```go
const (
    InitialPhase PhaseType = iota           // 設定読み込み
    DecideQuestionPhase                     // 問題選択
    DecideLanguagePhase                     // 言語選択
    DecidePathPhase                         // パス入力
    OverwriteConfirmPhase                   // 上書き確認
    GenerationPhase                         // ファイル生成
    DonePhase                               // 完了
)
```

### モデル構造

#### メインモデル（orchestrator）

```go
// internal/tui/model.go
type Model struct {
    model.Model                             // 共有データモデル
    currentPhase   phases.PhaseType
    currentHandler phases.PhaseHandler
}
```

#### 共有データモデル

```go
// internal/tui/model/model.go
type Model struct {
    Err            error
    Spinner        spinner.Model
    Loading        bool

    Config         *config.Config

    Input          Input                    // CLI引数
    Questions      []api.Question
    SelectedQ      *api.Question
    QuestionDetail *api.QuestionDetail
    Languages      []api.CodeSnippet
    SelectedLang   *api.CodeSnippet
    Cursor         int
    TextInput      textinput.Model
    Client         *api.Client
    OutPath        string

    OverwriteChoice int                     // 0: overwrite, 1: backup, 2: return, 3: quit
}
```

### フェーズ遷移の仕組み

1. **Enter** メソッドがフェーズ開始時の処理を実行
2. **Update** メソッドがユーザー入力や非同期処理の結果を処理
3. 次のフェーズを返すと、メインモデルが自動的に遷移を実行
4. **View** メソッドが現在のフェーズの UI を描画

### キーバインド

#### 共通

| キー   | アクション |
| ------ | ---------- |
| Ctrl+C | 強制終了   |
| q/Esc  | 終了       |

#### リスト選択フェーズ（DecideQuestion, DecideLanguage）

| キー   | アクション         |
| ------ | ------------------ |
| ↑/k    | カーソル上移動     |
| ↓/j    | カーソル下移動     |
| Enter  | 選択確定           |
| h/left | 前のフェーズに戻る |

#### パス入力フェーズ（DecidePath）

| キー     | アクション         |
| -------- | ------------------ |
| Enter    | 入力確定           |
| h/left   | 前のフェーズに戻る |
| 通常入力 | パス編集           |

#### 上書き確認フェーズ（OverwriteConfirm）

| キー  | アクション                 |
| ----- | -------------------------- |
| ↑/k   | カーソル上移動             |
| ↓/j   | カーソル下移動             |
| Enter | 選択確定                   |
| o     | 上書き（即実行）           |
| b     | バックアップ作成（即実行） |
| r     | パス入力に戻る             |

### フェーズ実装例

```go
// internal/tui/phases/decide_question.go
type DecideQuestionHandler struct{}

func (h *DecideQuestionHandler) Enter(m *model.Model) tea.Cmd {
    m.Loading = true
    m.Cursor = 0
    return h.fetchQuestionList(m)
}

func (h *DecideQuestionHandler) Update(m *model.Model, msg tea.Msg) (tea.Cmd, *PhaseType) {
    switch msg := msg.(type) {
    case fetchQuestionListMsg:
        m.Loading = false
        m.Questions = msg.questionList
        return nil, nil
    case tea.KeyMsg:
        switch msg.String() {
        case "enter":
            m.SelectedQ = &m.Questions[m.Cursor]
            next := DecideLanguagePhase
            return nil, &next  // フェーズ遷移
        }
    }
    return nil, nil
}

func (h *DecideQuestionHandler) View(m *model.Model) string {
    // UI描画ロジック
}
```

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
