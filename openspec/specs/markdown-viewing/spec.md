# markdown-viewing Specification

## Purpose
TBD - created by archiving change markdown-viewer. Update Purpose after archive.
## Requirements
### Requirement: Rendered display

システムは Markdown を CommonMark + GFM（表・取り消し線）としてレンダリングし、見出し・リスト・引用・コードブロック（シンタックスハイライト付き）・強調・リンク・表を装飾表示する（SHALL）。本文は端末幅に合わせて折り返す（SHALL）。端末のダーク / ライトに応じたスタイルを自動選択する（SHALL）。

#### Scenario: 見出しと表の整形表示

- **WHEN** 見出しとテーブル記法を含む Markdown を開く
- **THEN** 見出しは装飾され、テーブルは罫線付きで整形表示される

#### Scenario: 長い段落の折返し

- **WHEN** 端末幅を超える長さの段落を含む Markdown を開く
- **THEN** 段落は端末幅で折り返されて全文が表示され、水平スクロールは発生しない

### Requirement: Scroll navigation

システムは以下のキーでスクロールする（SHALL）: `j` / `k`（1 行下 / 上）、`Ctrl+d` / `Ctrl+u`（半画面下 / 上）、`f` / `Space`（1 画面下）、`b`（1 画面上）、`gg`（先頭）、`G`（末尾）。文書の先頭・末尾を超えるスクロールは行わない（SHALL）。`q` で終了し、端末を元の状態に復元する（SHALL）。

#### Scenario: 行・画面単位のスクロール

- **WHEN** 画面の高さを超える文書で `j` `Ctrl+d` `f` をそれぞれ押す
- **THEN** 表示が 1 行 / 半画面 / 1 画面ずつ下へスクロールする

#### Scenario: 先頭・末尾への移動と境界

- **WHEN** `G` を押し、続けて `j` を押す
- **THEN** 文書末尾が表示され、末尾を超えてスクロールしない
- **WHEN** `gg` を押す
- **THEN** 文書先頭が表示される

#### Scenario: q で終了する

- **WHEN** `q` を押す
- **THEN** Viewer が終了し、端末が元の状態に復元される

### Requirement: Text search

システムは `s` で検索プロンプトをステータスバーに表示し、入力した文字列をレンダリング後の表示テキストから検索する（SHALL）。`Enter` で確定すると現在の表示位置以降の最初の一致（無ければ先頭からの最初の一致）へスクロールし、全一致をハイライトし、現在の一致を区別して表示する（SHALL）。`n` / `p` で次 / 前の一致へ循環移動し（SHALL）、`Esc` で入力の中止またはハイライトの解除を行う（SHALL）。検索語に大文字が含まれない場合は大文字小文字を区別しない（SHALL）。

#### Scenario: 検索して一致へ移動する

- **WHEN** 画面外に一致を含む文書で `s` を押し、検索語を入力して `Enter` を押す
- **THEN** 一致する行が画面内に表示され、一致箇所がハイライトされ、ステータスバーに `[1/総数]` が表示される

#### Scenario: 次 / 前の一致へ循環する

- **WHEN** 検索確定後に `n` を繰り返し押す
- **THEN** 一致を順に移動し、最後の一致の次は最初の一致に戻る

#### Scenario: 一致しない

- **WHEN** 文書に存在しない文字列を検索する
- **THEN** ステータスバーに `not found` が表示され、表示位置は変わらない

### Requirement: Resize handling

システムは端末サイズの変更を検知し、新しい幅で再レンダリングして表示位置（スクロール割合）を概ね維持する（SHALL）。

#### Scenario: 端末リサイズで再レンダリングされる

- **WHEN** 文書の中程を表示中に端末の幅を変更する
- **THEN** 新しい幅で折返しが再計算され、表示位置は変更前と概ね同じ箇所を示す

### Requirement: CJK display

システムは日本語等の全角文字を含む文書を、桁崩れ・折返し不良なく表示する（SHALL）。

#### Scenario: 日本語段落の表示

- **WHEN** 日本語の段落・見出し・表を含む Markdown を開く
- **THEN** 全角文字が正しい幅で描画され、折返し・罫線が崩れない

### Requirement: Status bar

システムは画面最下部にステータスバーを常時表示し、ファイル名・スクロール位置（%）・主要キーのヒント（`s:search e:edit r:reload t:theme q:quit`）を表示する（SHALL）。エラーや通知のメッセージは同領域に一時表示する（SHALL）。

#### Scenario: ステータスバーの内容

- **WHEN** `README.md` を開いて文書の中程までスクロールする
- **THEN** ステータスバーにファイル名 `README.md`、位置（例 `42%`）、キーヒントが表示される

