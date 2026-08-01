# Proposal: markdown-viewer

## Why

ターミナル上で Markdown を読みやすく閲覧し、そのまま vi で編集して閲覧に戻る、という往復を 1 コマンドで完結させたい。当初構想（`vi-style-markdown-editor`）は vi 風エディタの自前実装だったが、vi 互換性の妥協・IME 対応・表示と編集の対応付けという制約が大きいため、編集は本物の vi（外部コマンド）へ委譲し、medit はレンダリング表示に特化した Viewer とする方針に転換する。

## What Changes

- Go 製 TUI の Markdown Viewer `medit` を新規実装する
- `medit <file.md>` で Glamour によるフルレンダリング表示（見出し・リスト・表・コードブロック・強調・端末幅での折返し）
- vi / less 風のスクロールナビゲーション（`j` `k` `Ctrl+d` `Ctrl+u` `f` `b` `Space` `gg` `G`、`q` で終了）
- `e` キーで TUI を一時停止して外部エディタ（`$VISUAL` → `$EDITOR` → `vi`）を起動する。表示中の位置に対応する行番号を `+N` 引数で引き継ぐ
- エディタ終了で自動的に Viewer へ復帰し、ファイルを再読込して最新内容を表示する（スクロール位置は維持）
- `r` で手動再読込
- 存在しないパスは新規ファイル扱いで空表示し、`e` からの vi での保存によりファイルが作成される
- ステータスバー（ファイル名・スクロール位置・キーヒント・メッセージ）

## Capabilities

### New Capabilities

- `file-loading`: Markdown ファイルの読み込み・新規ファイル扱い・再読込・エラー処理
- `markdown-viewing`: Glamour によるレンダリング表示、スクロールナビゲーション、リサイズ追従、ステータスバー
- `external-editing`: 外部エディタ（vi）の起動・終了後の復帰・自動再読込・表示位置の引き継ぎ

### Modified Capabilities

（なし — 新規プロジェクトのため既存仕様はない）

## Impact

- 新規 Go モジュール `github.com/sudabon/medit`（リポジトリ初のアプリケーションコード）
- 主要依存: Bubble Tea（TUI）、bubbles/viewport（スクロール）、Glamour（Markdown レンダリング）、Lip Gloss（スタイリング）
- 開発基盤: `go test` / golangci-lint / GitHub Actions を最小限整備する
- 既存の変更 `openspec/changes/vi-style-markdown-editor/` は本変更に置き換えられる（未実装のまま残置。不要になった時点で削除可）
- 実行環境は vi が利用できる Unix 系（macOS / Linux）を対象とする

## Open Questions

回答が得られるまでは各項目の「想定」をデフォルトとして設計・実装を進める。

1. **編集キー**: `e` に加えて `i` でも vi を起動するか（想定: `e` と `i` の両方。「編集に入る」という vi のメタファーに合わせる）
2. **エディタの選択**: `$VISUAL` → `$EDITOR` → `vi` の順で解決する想定でよいか（常に `vi` 固定にすることも可能）
3. **Viewer 内のテキスト検索（`/`）・見出しジャンプ（TOC）**: MVP 対象外とし、次の変更で検討する想定でよいか
4. **ファイル変更の自動検知（fsnotify）**: MVP 対象外（エディタ復帰時と `r` のみで再読込）とする想定でよいか
