# Tasks: markdown-viewer

## 1. プロジェクト基盤

- [x] 1.1 `go.mod` を初期化（module `github.com/sudabon/medit`、最新安定版 Go）し、依存（bubbletea / bubbles / glamour / lipgloss）を追加する
- [x] 1.2 パッケージ骨格（`cmd/medit`, `internal/loader`, `internal/render`, `internal/ui`）と main エントリポイントを作成し、引数検証（1 引数以外は usage + 非 0 終了）を実装する
- [x] 1.3 golangci-lint 設定と GitHub Actions ワークフロー（push / PR で test + vet + lint）を追加する

## 2. 読み込みとレンダリング（specs: file-loading / markdown-viewing）

- [x] 2.1 loader を実装する: UTF-8 読み込み、存在しないパスの新規ファイル扱い、ディレクトリ・権限エラーの報告、再読込 API（ユニットテスト含む）
- [x] 2.2 render を実装する: Glamour ラッパ（幅指定・ダーク/ライト自動スタイル・GFM 有効化）、内容/幅の変化時のみ再レンダリングする API（見出し・表・日本語折返しのスモークテスト含む）
- [x] 2.3 スクロール割合 → ソース行番号の近似変換を実装する（`vi +N` 用。境界のクランプ含む、ユニットテスト含む）

## 3. Viewer UI（specs: markdown-viewing）

- [x] 3.1 Bubble Tea の Model / Update / View 骨格と bubbles/viewport を統合し、`q` での終了と端末復元を実装する
- [x] 3.2 ナビゲーションキーを実装する: `j` `k` `Ctrl+d` `Ctrl+u` `f` `Space` `b` `gg` `G`、文書先頭・末尾での境界クランプ
- [x] 3.3 端末リサイズ対応を実装する: 新しい幅での再レンダリングとスクロール割合の維持
- [x] 3.4 ステータスバーを実装する: ファイル名・位置（%）・キーヒント・新規ファイル表示・一時メッセージ表示

## 4. 外部エディタ連携（specs: external-editing）

- [x] 4.1 `tea.ExecProcess` によるエディタ起動を実装する: `$VISUAL` → `$EDITOR` → `vi` の解決（フィールド分割）、`+N <path>` 引数、`e` / `i` キー割り当て
- [x] 4.2 復帰処理を実装する: 自動再読込・再レンダリング・スクロール位置の維持（クランプ）、非 0 終了時のメッセージ表示
- [x] 4.3 新規ファイル作成フローを実装する: 新規扱いからの `e` → 保存 → 復帰で内容表示と新規表示の解除、`r` 手動再読込とファイル消失時の表示維持 + エラー表示

## 5. 検証と仕上げ

- [x] 5.1 全ユニットテスト・`go vet`・golangci-lint がパスすることを確認する
- [x] 5.2 spec のシナリオを主要ターミナル（Terminal.app / iTerm2 / WezTerm）で手動確認する（vi 往復、`:q!` 時の非反映、日本語文書の折返し・表の表示、新規ファイル作成フロー、リサイズ）
- [x] 5.3 README を更新する: インストール方法、使い方、キー一覧、`$VISUAL` / `$EDITOR` の扱い
