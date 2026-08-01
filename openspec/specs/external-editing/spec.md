# external-editing Specification

## Purpose
TBD - created by archiving change markdown-viewer. Update Purpose after archive.
## Requirements
### Requirement: Launch external editor

システムは `e` または `i` キーで TUI を一時停止し、端末を明け渡して外部エディタを起動する（SHALL）。エディタは `$VISUAL` → `$EDITOR` → `vi` の順で解決する（SHALL）。起動引数には対象ファイルのパスと、現在の表示位置に対応する近似ソース行番号 `+N` を渡す（SHALL）。

#### Scenario: e キーで vi が起動する

- **WHEN** `$VISUAL` と `$EDITOR` が未設定の状態で `e` を押す
- **THEN** Viewer の描画が停止し、`vi +N <path>` が端末を占有して起動する

#### Scenario: $EDITOR が尊重される

- **WHEN** `$EDITOR=vim` の状態で `e` を押す
- **THEN** vim が起動する

#### Scenario: 表示位置が引き継がれる

- **WHEN** 文書の末尾付近を表示した状態で `e` を押す
- **THEN** エディタは文書末尾付近の行を初期位置として開く

### Requirement: Return and auto-reload

システムはエディタプロセスの終了を検知して Viewer へ復帰し、ファイルを自動で再読込・再レンダリングする（SHALL)。スクロール位置は復帰前の位置を維持し、内容が短くなった場合は末尾にクランプする（SHALL）。エディタが非 0 で終了した場合もViewer へ復帰し、ステータスバーにメッセージを表示したうえで再読込は実行する（SHALL）。

#### Scenario: 保存した変更が反映される

- **WHEN** `e` で vi を開き、内容を編集して `:wq` で終了する
- **THEN** Viewer に復帰し、編集後の内容がレンダリング表示される

#### Scenario: 保存しなければ表示は変わらない

- **WHEN** `e` で vi を開き、`:q!` で終了する
- **THEN** Viewer に復帰し、表示内容は編集前と同一である

#### Scenario: スクロール位置の維持

- **WHEN** 文書の中程を表示した状態で vi を往復する
- **THEN** 復帰後も概ね同じ箇所が表示されている

### Requirement: Create new file via editor

システムは新規ファイル扱いで開いている場合も `e` でエディタを起動し（SHALL）、エディタ内での保存によってファイルが作成された場合、復帰時の再読込でその内容を表示し新規ファイル表示を解除する（SHALL）。

#### Scenario: 新規ファイルを vi で作成する

- **WHEN** 存在しない `new.md` を開いた状態で `e` を押し、vi で内容を書いて `:wq` する
- **THEN** `new.md` がディスク上に作成され、復帰した Viewer にその内容がレンダリング表示される

