# medit

ターミナル用の Markdown ビューアです。Glamour でレンダリングし、less / vi 風にスクロール、編集は本物のエディタに任せます。

## インストール

```bash
go install github.com/sudabon/medit/cmd/medit@latest
```

リポジトリを clone した場合:

```bash
go build -o medit ./cmd/medit
```

## 使い方

```bash
medit README.md
medit path/to/notes.md
```

存在しないパスを指定すると新規ファイルとして開きます（エディタで保存するまでディスク上には作成されません）。

## キー操作

| キー | 動作 |
|------|------|
| `j` / `k` | 1 行下 / 上へスクロール |
| `Ctrl+d` / `Ctrl+u` | 半画面下 / 上へスクロール |
| `f` / `Space` | 1 画面下へスクロール |
| `b` | 1 画面上へスクロール |
| `gg` | 先頭へ移動 |
| `G` | 末尾へ移動 |
| `e` / `i` | 外部エディタで開く |
| `r` | ファイルを再読込 |
| `t` | ライト / ダークテーマ切替 |
| `q` | 終了 |

## 外部エディタ

`e` または `i` で Viewer を一時停止し、エディタを起動します。引数は `+N <path>`（表示位置に対応する近似のソース行番号）です。

解決順:

1. `$VISUAL`
2. `$EDITOR`
3. `vi`

値は空白でフィールド分割します（シェルは経由しません）。例:

```bash
export VISUAL=nvim
export EDITOR="vim -N"
```

クォートを含む複雑なコマンドラインには対応していません。`$VISUAL` / `$EDITOR` はコマンドと単純なフラグにとどめてください。

エディタ終了後、medit はファイルを再読込し、スクロール位置を復元します（文書が短くなった場合は末尾にクランプ）。終了コードが非 0 でも Viewer に復帰して再読込し、ステータスバーにメッセージを表示します。

## 表示

テーマは次の順で自動選択されます。

1. `GLAMOUR_STYLE=dark` または `light`（明示指定）
2. OS の外観設定（macOS のライト / ダーク、GNOME の color-scheme）
3. 端末の背景色（`COLORFGBG` / 問い合わせ）
4. フォールバック: dark

自動判定が合わないときは `t` で切り替えできます。強制指定の例:

```bash
GLAMOUR_STYLE=light medit README.md
GLAMOUR_STYLE=dark medit README.md
```

本文は約 88 文字幅の固定カラムでレンダリングし、端末内で中央寄せします。フェンス付きコードブロックは角丸の枠線で囲んで表示します。

フォントサイズは medit からは変更できません。端末エミュレータ（または cmux）側の設定を使ってください。

## 動作環境

- macOS または Linux
- 上記のとおり利用できるターミナルエディタ（既定は vi）
