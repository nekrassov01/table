# API

This document is the user guide to the public API.

## Table of contents

- [API](#api)
  - [Table of contents](#table-of-contents)
  - [Core API](#core-api)
    - [Table](#table)
    - [Stream](#stream)
    - [Output differences](#output-differences)
  - [Common interfaces](#common-interfaces)
  - [Row adapters](#row-adapters)
  - [Output specifications](#output-specifications)
  - [Errors](#errors)
    - [Error structure](#error-structure)
    - [Sentinel errors](#sentinel-errors)
    - [Output after an error](#output-after-an-error)
  - [Concurrent use](#concurrent-use)
  - [Options](#options)
    - [Feature matrix](#feature-matrix)
    - [Defaults](#defaults)
    - [Application order](#application-order)
    - [Argument ownership](#argument-ownership)
  - [Option reference](#option-reference)
    - [text](#text)
    - [html](#html)
    - [markdown](#markdown)
    - [backlog](#backlog)
    - [csv](#csv)
  - [Column resolution and settings](#column-resolution-and-settings)
    - [Column count](#column-count)
    - [Column selection](#column-selection)
    - [Scopes](#scopes)
  - [Value display](#value-display)
    - [String conversion](#string-conversion)
    - [Transformers](#transformers)
    - [Placeholders](#placeholders)
  - [Dynamic footers](#dynamic-footers)
  - [Cell spans](#cell-spans)
  - [Alignment](#alignment)
  - [Indexes](#indexes)
  - [Captions](#captions)

## Core API

Every output package provides `Table` and `Stream`. Constructors apply each `Option` in the order supplied. Settings cannot be changed after construction.

```go
func NewTable(w io.Writer, opts ...Option) *Table
func (t *Table) Render(rows [][]any) error

func NewStream(w io.Writer, opts ...Option) *Stream
func (s *Stream) Render(row []any) error
func (s *Stream) Close() error
```

### Table

`Table` accepts the complete body and writes an independent table for each call to `Render`.

- `Render` can be called repeatedly on the same `Table`.
- `Render(nil)` represents an empty body. Configured headers or footers can still establish a table and include its caption where supported.
- Constructor settings and the destination `io.Writer` are reused.

### Stream

`Stream` writes one table a row at a time.

1. Call `Render` zero or more times.
2. Call `Close` after the input ends.
3. Discard the stream after `Close`.

Repeated calls to a successful `Close` return `nil`, and a subsequent call to `Render` returns `table.ErrClosed`. If a write error or an error during `Close` has been recorded, later calls return that same error. `Close` writes any footer, closing tags, or bottom border and releases internal workspaces. Call it even when no footer is configured.

### Output differences

`Table` and `Stream` produce the same logical content when lookahead does not affect the result. The following differences are part of the contract.

| Format     | Aspect         | `Table`                                                    | `Stream`                                                                        |
| ---------- | -------------- | ---------------------------------------------------------- | ------------------------------------------------------------------------------- |
| `text`     | Width          | Derives widths from the complete header, body, and footer. | Freezes widths when output starts, then wraps or truncates later values.        |
| `text`     | Index          | Uses the width required by the complete body row count.    | Reserves three digits by default; `WithIndexWidth` can override the minimum.    |
| `html`     | Vertical spans | Emits `rowspan`.                                           | Cannot know future run lengths and emits continuation positions as empty cells. |
| `markdown` | Padding        | Pads to the greatest display width among all rows.         | Does not pad to a fixed width.                                                  |
| `backlog`  | Padding        | Pads to the greatest display width among all rows.         | Does not pad to a fixed width.                                                  |

Whitespace added by `markdown` and `backlog` for alignment is not part of the parsed cell value.

## Common interfaces

```go
type Tabular interface {
    Render(rows [][]any) error
}

type Streamer interface {
    Render(row []any) error
    Close() error
}
```

`Tabular` abstracts batch output and `Streamer` abstracts row-at-a-time output. Each format's `Table` implements `Tabular`, and each `Stream` implements `Streamer`.

## Row adapters

The root package provides helpers that convert typed data into rows accepted by `Render`.

```go
func TableOf[T any](values []T, fn func(T) []any) [][]any
func StreamOf[T any](values iter.Seq2[T, error], fn func(T) []any) iter.Seq2[[]any, error]
```

- `TableOf` converts a typed slice to `[][]any`.
- `StreamOf` converts each iterator value to `[]any` and stops at the first error.

## Output specifications

The following table identifies the primary consumer, reference specification, and conformance status of each format.

| Format     | Primary consumer                  | Reference                                                                                                                                                                                                        | Conformance                                                                                                                  |
| ---------- | --------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| `text`     | Terminal or text viewer           | Terminal display width<br/>ECMA-48 SGR                                                                                                                                                                           | Custom format                                                                                                                |
| `html`     | Browser                           | [HTML Standard: Tables](https://html.spec.whatwg.org/multipage/tables.html#tables)                                                                                                                               | Substantially conformant. `colspan` is limited to 64 columns, while `rowspan` may exceed the standard maximum of 65,534.     |
| `markdown` | GitHub or a GFM-compatible parser | [GitHub Flavored Markdown: Tables extension](https://github.github.com/gfm/#tables-extension-)                                                                                                                   | Substantially conformant. A body row wider than the header returns `table.ErrColumnCount` instead of ignoring extra cells.   |
| `backlog`  | Backlog notation parser           | [Backlog Notation](https://help-center.backlog.com/%E3%83%86%E3%82%AD%E3%82%B9%E3%83%88%E6%95%B4%E5%BD%A2%E3%83%AB%E3%83%BC%E3%83%AB%EF%BC%88Backlog%E8%A8%98%E6%B3%95%EF%BC%89-6a1d4d7f3abb3ada78c5658b#index7) | Substantially conformant. Footers use header-cell notation, and span continuations use empty cells.                          |
| `csv`      | `encoding/csv`-compatible reader  | [encoding/csv](https://pkg.go.dev/encoding/csv)<br/>[RFC 4180](https://www.rfc-editor.org/rfc/rfc4180)                                                                                                           | `encoding/csv.Writer`-compatible quoting, validation, and CRLF handling. Tab/LF by default; options select RFC 4180.         |

## Errors

Errors are structured for standard Go error handling.

### Error structure

Every non-nil error returned by `Render` or `Close` has an outer `*table.Error` containing the output package name and cause. Input errors unwrap to a shared sentinel, with an intermediate cause when the error includes details. A write error retains both `table.ErrWriteFailed` and the error returned by the destination.

```go
var tableErr *table.Error
if errors.As(err, &tableErr) {
    log.Printf("package: %s", tableErr.Pkg)
}
```

### Sentinel errors

Use `errors.Is` to test the following sentinels.

| Error                     | Condition                                                | Recovery                                                                             |
| ------------------------- | -------------------------------------------------------- | ------------------------------------------------------------------------------------ |
| `table.ErrWriteFailed`    | The destination returned an error.                       | Inspect the destination and construct a new `Table` or `Stream` if necessary.        |
| `table.ErrClosed`         | `Stream.Render` was called after a successful `Close`.   | Construct a new `Stream`.                                                            |
| `table.ErrColumnCount`    | A body row or footer exceeded the resolved column count. | Correct the input. After `Stream.Render`, a later valid row can continue the stream. |
| `table.ErrHeaderRequired` | A Markdown header had zero cells.                        | Configure a non-empty header with `WithHeader`.                                      |
| `table.ErrDelimiter`      | CSV received an invalid delimiter.                       | Supply a valid `rune`.                                                               |

### Output after an error

Column-count errors follow these rules:

- `Table` validates every row before writing the body and therefore writes nothing.
- `Stream.Render` does not write an invalid row or advance its row index. A later row within the resolved column count can still be written.
- If a dynamic footer exceeds a stream's column count, `Close` returns an error and omits the footer. Any header and body already written remain in the destination.
- When indexing is enabled, `got` and `want` in the error message are logical column counts that include the generated index column.

Write errors follow these rules:

- `errors.Is` and `errors.As` traverse the causes inside `table.Error`.
- `Stream` retains the first write error and returns it from later calls.
- Successful writes before an error are not rolled back. Both `Table` and `Stream` may leave partial output.
- For atomic output, render to a `bytes.Buffer` and copy it to the final destination only after rendering succeeds.

## Concurrent use

Individual `Table` and `Stream` instances do not synchronize method calls. Do not call methods on the same instance concurrently from multiple goroutines.

Separately constructed instances can be used concurrently. When instances share an `io.Writer`, borrowed settings, or state captured by footer and transformer closures, the caller must provide the required synchronization.

## Options

Each output package defines its own closed set of options, while all formats follow the same configuration model.

### Feature matrix

| Feature                  | `text`                                          | `html`                                          | `markdown`                         | `backlog`                                       | `csv`                              |
| ------------------------ | ----------------------------------------------- | ----------------------------------------------- | ---------------------------------- | ----------------------------------------------- | ---------------------------------- |
| Static header            | `WithHeader`<br/>Any number of rows             | `WithHeader`<br/>Any number of rows             | `WithHeader`<br/>Required, one row | `WithHeader`<br/>Any number of rows             | `WithHeader`<br/>Optional, one row |
| Dynamic footer           | `WithFooter`                                    | `WithFooter`                                    | -                                  | `WithFooter`                                    | `WithFooter`                       |
| Index column             | `WithIndex`<br/>`WithIndexWidth`                | `WithIndex`                                     | `WithIndex`                        | `WithIndex`                                     | `WithIndex`                        |
| Placeholder              | `WithPlaceholder`                               | `WithPlaceholder`                               | `WithPlaceholder`                  | `WithPlaceholder`                               | `WithPlaceholder`                  |
| Value transformation     | `WithTransformer`                               | `WithTransformer`                               | `WithTransformer`                  | `WithTransformer`                               | `WithTransformer`                  |
| Alignment                | `WithAlign`                                     | `WithAlign`                                     | `WithAlign`                        | -                                               | -                                  |
| Vertical cell spanning   | `WithRowspan`                                   | `WithRowspan`                                   | `WithRowspan`                      | `WithRowspan`                                   | -                                  |
| Horizontal cell spanning | `WithColspan`                                   | `WithColspan`                                   | `WithColspan`                      | `WithColspan`                                   | -                                  |
| Target sections          | `ScopeHeader`<br/>`ScopeBody`<br/>`ScopeFooter` | `ScopeHeader`<br/>`ScopeBody`<br/>`ScopeFooter` | `ScopeHeader`<br/>`ScopeBody`      | `ScopeHeader`<br/>`ScopeBody`<br/>`ScopeFooter` | -                                  |
| Color                    | `WithAttr`                                      | `WithColor`                                     | `WithColor`                        | `WithColor`                                     | -                                  |
| Decoration               | `WithAttr`                                      | `WithDecoration`                                | `WithDecoration`                   | `WithDecoration`                                | -                                  |
| Caption                  | `WithCaption`                                   | `WithCaption`                                   | -                                  | -                                               | -                                  |
| Fixed column width       | `WithWidth`                                     | `WithTableAttr`<br/>`WithCellAttr`              | -                                  | -                                               | -                                  |
| In-cell wrapping         | `WithWidth`                                     | `WithTableAttr`<br/>`WithCellAttr`              | -                                  | -                                               | -                                  |
| In-cell truncation       | `WithTruncate`                                  | `WithTableAttr`<br/>`WithCellAttr`              | -                                  | -                                               | -                                  |
| Fit to terminal width    | `WithAutoFit`                                   | -                                               | -                                  | -                                               | -                                  |
| Padding                  | `WithPadding`                                   | `WithTableAttr`<br/>`WithCellAttr`              | -                                  | -                                               | -                                  |
| Omit inter-row borders   | `WithCompact`                                   | `WithTableAttr`<br/>`WithCellAttr`              | -                                  | -                                               | -                                  |
| Borders                  | `WithStyle`                                     | `WithTableAttr`<br/>`WithCellAttr`              | -                                  | -                                               | -                                  |
| Delimiter                | -                                               | -                                               | -                                  | -                                               | `WithDelimiter`                    |
| Record ending            | -                                               | -                                               | -                                  | -                                               | `WithCRLF`                         |

A dash means that the output format has no corresponding feature.

### Defaults

| Setting             | `text`                                | `html`                                     | `markdown`      | `backlog`       | `csv`        |
| ------------------- | ------------------------------------- | ------------------------------------------ | --------------- | --------------- | ------------ |
| Placeholder         | One ASCII space                       | Empty string                               | One ASCII space | One ASCII space | Empty string |
| Header alignment    | Center                                | CSS default                                | GFM default     | -               | -            |
| Body alignment      | Left                                  | CSS default                                | GFM default     | -               | -            |
| Footer alignment    | Left                                  | CSS default                                | No footer       | -               | -            |
| Index alignment     | Center in the header, right elsewhere | CSS default in the header, right elsewhere | Right           | -               | -            |
| Border or delimiter | `StyleLight`                          | HTML table elements                        | `\|`            | `\|`            | Tab          |
| Line ending         | LF                                    | LF                                         | LF              | LF              | LF           |
| Caption position    | Bottom                                | CSS default                                | -               | -               | -            |

### Application order

Options are applied in the order supplied. A later global setting replaces an earlier one. A later column setting replaces the same setting only for the selected columns and `Scope` values.

- `WithIndex`, `text.WithCompact`, `text.WithAutoFit`, and `csv.WithCRLF` only enable a feature and cannot disable it.
- `text.WithTruncate`, `WithRowspan`, and `WithColspan` accumulate selected columns and, where accepted, scopes.
- `text.WithIndexWidth` enables indexing. A positive value replaces the width, while zero or a negative value leaves an existing width unchanged.

```go
text.WithAlign(text.ScopeBody, text.AllColumns(), text.AlignLeft),
text.WithAlign(text.ScopeBody, text.Columns(2), text.AlignRight),
```

This example aligns every body column to the left, then aligns input column 2 to the right.

### Argument ownership

The public API treats slices, pointers, and functions as follows.

| Input                                                               | Treatment                                                                       |
| ------------------------------------------------------------------- | ------------------------------------------------------------------------------- |
| Column indexes passed to `Columns`                                  | Cloned at construction and owned by the returned `ColumnSelector`.              |
| Rows passed to `WithHeader`                                         | Borrowed read-only by `Table` or `Stream`.                                      |
| `text.Style`                                                        | The struct is copied; internal pointers and byte slices are borrowed read-only. |
| `text.WithAttr`                                                     | The `Attr` pointer and its byte slices are borrowed read-only.                  |
| Colors and decorations in `html`, `markdown`, and `backlog`         | Referenced markup is borrowed read-only.                                        |
| HTML table, section, and cell attributes                            | Copied by value and normalized when the `Option` is constructed.                |
| `WithFooter` and `WithTransformer`                                  | The function is retained; the caller owns any state captured by its closure.    |
| Body rows and values returned by a footer function                  | Borrowed read-only only for the corresponding call to `Render` or `Close`.      |

Do not mutate borrowed values or closure state while the associated `Table` or `Stream` is in use. Defensive copies of headers and markup are not part of the API contract. References to body and footer rows are discarded before the corresponding call returns.

Settings captured by value, such as strings, numbers, booleans, and runes, are independent of the caller. A constructed `Option` can be reused by multiple `Table` and `Stream` instances.

## Option reference

### text

```go
func WithHeader(rows ...[]string) Option
func WithFooter(fn func() [][]string) Option
func WithCaption(s string, side CaptionSide) Option
func WithStyle(style Style) Option
func WithCompact() Option
func WithIndex() Option
func WithIndexWidth(n int) Option
func WithAutoFit() Option
func WithPlaceholder(s string) Option
func WithAlign(scopes Scope, columns ColumnSelector, align AlignSide) Option
func WithWidth(columns ColumnSelector, width int) Option
func WithTruncate(columns ColumnSelector) Option
func WithPadding(columns ColumnSelector, left, right int) Option
func WithRowspan(scopes Scope, columns ColumnSelector) Option
func WithColspan(scopes Scope, columns ColumnSelector) Option
func WithAttr(scopes Scope, columns ColumnSelector, attr *Attr) Option
func WithTransformer(columns ColumnSelector, fn func(any) (string, *Attr)) Option
```

- `WithWidth` sets the display-width boundary of cell content, excluding padding. A value of zero or less removes the boundary.
- `WithTruncate` replaces wrapping with `...`. In a column without `WithWidth`, it applies to values that exceed the initial column width after a stream has begun.
- `WithPadding` sets left and right space widths. Negative values become zero, and padding contributes to the total table width.
- `WithAutoFit` reduces column widths to fit terminal output within the terminal width. It has no effect for a non-terminal destination or if any column uses `WithWidth` or `WithTruncate`.
- `WithIndexWidth` sets the minimum width of the index column. A stream reserves at least three digits.
- `WithCompact` omits horizontal borders between body rows. At a vertically spanned cell boundary, it retains the horizontal segments for cells that are not spanned.
- Cell widths use Unicode terminal display widths. Wrapping preserves grapheme clusters; a cluster wider than the boundary remains intact and makes that physical output line wider than the column. When a stream begins rendering body rows, it reserves one display cell for a column whose initial content has zero display width.
- LF, CR, and CRLF are in-cell line breaks; output lines end with LF.
- Invalid UTF-8 bytes are preserved rather than replaced, and ANSI sequences embedded in values are not parsed.

`Attr` represents ANSI SGR color and decoration. `WithAttr` overrides `Style.Content` per column, and an `Attr` returned by a transformer overrides it per cell. ANSI attributes are omitted when the destination is not a terminal.

`NewAttr` combines multiple `Code` values into one SGR sequence. It returns `nil` when called without arguments.

Built-in borders are `StyleASCII`, `StyleLight`, `StyleRounded`, `StyleHeavy`, and `StyleDouble`, together with colored variants of the latter four. Set individual `Style` fields to define custom borders and attributes.

### html

```go
func WithHeader(rows ...[]string) Option
func WithFooter(fn func() [][]string) Option
func WithCaption(s string, side CaptionSide) Option
func WithTableAttr(attr TableAttr) Option
func WithIndex() Option
func WithPlaceholder(s string) Option
func WithAlign(scopes Scope, columns ColumnSelector, align AlignSide) Option
func WithRowspan(scopes Scope, columns ColumnSelector) Option
func WithColspan(scopes Scope, columns ColumnSelector) Option
func WithColor(scopes Scope, columns ColumnSelector, color *Color) Option
func WithDecoration(scopes Scope, columns ColumnSelector, decoration *Decoration) Option
func WithCellAttr(scopes Scope, columns ColumnSelector, attr Attr) Option
func WithTransformer(columns ColumnSelector, fn func(any) (string, *Color, *Decoration)) Option
```

- Output contains a `table` element only when at least one column exists. It includes `caption`, `thead`, `tbody`, and `tfoot` only when the corresponding section exists.
- Headers use `thead` and `th`, the body uses `tbody` and `td`, and footers use `tfoot` and `td`. Empty sections are omitted.
- `Attr` holds the classes and inline style for one element. `TableAttr` groups attributes for the table and its sections, while `SectionAttr` groups section, row, and cell attributes.
- `WithCellAttr` appends per-column cell attributes. Classes are joined with an ASCII space and styles with a semicolon; alignment from `WithAlign` is appended last.
- Displayed values are escaped as HTML text. CR, LF, and CRLF become `<br>`. C0 controls other than tab, CR, and LF, together with DEL and invalid UTF-8, become U+FFFD.
- Values in the `Class` and `Style` fields, and color strings passed to `NewColor`, are escaped as HTML attributes. NUL, C0 controls other than ASCII whitespace, DEL, C1 controls, Unicode noncharacters, and invalid UTF-8 become U+FFFD. CR and CRLF are normalized to LF, while tab, LF, and form feed are preserved.
- Decoration presets are `DecorationBold`, `DecorationUnderline`, `DecorationItalic`, `DecorationStrikethrough`, `DecorationCode`, and `DecorationPreformatted`.
- When decoration and color are combined, the decoration element is outside the color `span`.
- `NewDecoration` writes the supplied markup without escaping. Pass only trusted HTML. It returns `nil` when `prefix` is empty.
- `NewColor` returns `nil` when both foreground and background are empty.

### markdown

```go
func WithHeader(header []string) Option
func WithIndex() Option
func WithPlaceholder(s string) Option
func WithAlign(columns ColumnSelector, align AlignSide) Option
func WithRowspan(columns ColumnSelector) Option
func WithColspan(columns ColumnSelector) Option
func WithColor(scopes Scope, columns ColumnSelector, color *Color) Option
func WithDecoration(scopes Scope, columns ColumnSelector, decoration *Decoration) Option
func WithTransformer(columns ColumnSelector, fn func(any) (string, *Color, *Decoration)) Option
```

- A GFM table requires one header row. Omitting `WithHeader` or supplying an empty header produces `table.ErrHeaderRequired`.
- `WithAlign` sets alignment markers on the GFM delimiter row.
- Backslashes, vertical bars, backticks, emphasis markers, square brackets, angle brackets, and ampersands in displayed values are escaped. NUL and invalid UTF-8 become U+FFFD.
- Strings resembling URLs or email addresses may be autolinked by a GFM implementation.
- Color uses an HTML `span`. Other decorations surround the color span; with `DecorationCode`, the color span surrounds the code span.
- `DecorationCode` follows the GFM code-span rules. Its fence is longer than any backtick run in the value, and LF, CR, and CRLF become spaces. When the normalized content begins and ends with spaces but is not entirely spaces, the emitted span adds one space at each end so GFM parsing preserves them.
- `DecorationPreformatted` preserves whitespace with `<pre>`.
- `NewColor` escapes a CSS color as an HTML attribute. Vertical bars become character references so they cannot split a GFM table row. Invalid characters follow HTML attribute replacement rules, and line breaks become spaces. CSS validity is not checked.
- `NewDecoration` writes the supplied delimiters without escaping. Pass only trusted markup. It returns `nil` when `prefix` is empty.

### backlog

```go
func WithHeader(rows ...[]string) Option
func WithFooter(fn func() [][]string) Option
func WithIndex() Option
func WithPlaceholder(s string) Option
func WithRowspan(scopes Scope, columns ColumnSelector) Option
func WithColspan(scopes Scope, columns ColumnSelector) Option
func WithColor(scopes Scope, columns ColumnSelector, color *Color) Option
func WithDecoration(scopes Scope, columns ColumnSelector, decoration *Decoration) Option
func WithTransformer(columns ColumnSelector, fn func(any) (string, *Color, *Decoration)) Option
```

- Header and footer cells use header-cell notation beginning with `~`. A footer is a library-level section; Backlog itself does not distinguish it.
- Backslashes, vertical bars, CR, and LF are escaped for Backlog notation. Invalid UTF-8 bytes are preserved rather than replaced.
- The header-cell `~` immediately follows the opening vertical bar, and padding follows the value.
- When color is combined with any decoration other than `DecorationCode`, the color notation surrounds the decoration.
- Backlog notation cannot represent `DecorationCode` and color simultaneously, so code decoration is retained and color is omitted.
- `NewColor` combines foreground and background into Backlog color notation. It returns `nil` for empty values, reserved characters, or line breaks.
- `NewDecoration` writes the supplied delimiters without escaping. Pass only trusted Backlog notation. It returns `nil` when `prefix` is empty.

### csv

```go
func WithHeader(header []string) Option
func WithFooter(fn func() [][]string) Option
func WithDelimiter(delimiter rune) Option
func WithCRLF() Option
func WithIndex() Option
func WithPlaceholder(s string) Option
func WithTransformer(columns ColumnSelector, fn func(any) string) Option
```

- The default delimiter is a tab. `WithDelimiter` changes it, and an invalid rune produces `table.ErrDelimiter`.
- A field is quoted when it contains the delimiter, a double quote, CR, or LF; begins with Unicode whitespace; or is exactly `\.`. Double quotes within a quoted field are doubled. These rules match `encoding/csv.Writer`.
- Records end with LF by default. `WithCRLF` changes record endings according to the same rules as `encoding/csv.Writer.UseCRLF`.
- Combining `WithDelimiter(',')` and `WithCRLF()` selects the delimiter and record ending specified by RFC 4180.
- A header has one row. Footers are emitted as ordinary records.
- Invalid UTF-8 bytes are preserved rather than replaced.

## Column resolution and settings

The following rules determine the table's column count, the input columns selected by a column option, and the sections to which a setting applies.

### Column count

The input column count, excluding an index, is determined by the format and input.

- `markdown` uses the width of its required header.
- Other formats use the widest non-empty header row.
- Without a header, `Table` uses the greater of the first non-empty body row and the widest footer row.
- Without a header, `Stream` uses its first non-empty body row. If `Close` is called without emitting a body, it uses the widest footer row.

The following rules then apply:

- Zero-column body rows received before the column count is established are not emitted.
- Outside Markdown, empty header or footer rows produce no output when no other section establishes a column.
- Short rows are extended with missing cells.
- A body row or footer wider than the resolved column count produces `table.ErrColumnCount`. A wider later body row never expands the column count.

### Column selection

Column-specific options select input columns with `ColumnSelector`.

```go
func Columns(indexes ...int) ColumnSelector
func AllColumns() ColumnSelector
```

- `Columns` selects zero-based input column indexes and ignores negative indexes.
- `AllColumns` selects every input column, including columns resolved later.
- Selecting a nonexistent column does not add an output column. The generated index column is never selectable.

Enabling `WithIndex` does not shift input column indexes.

### Scopes

`Scope` selects the section to which a column setting applies.

```go
ScopeHeader
ScopeBody
ScopeFooter
```

Scopes can be combined, as in `ScopeHeader | ScopeBody`. Zero and undefined bits are ignored. Markdown defines only `ScopeHeader` and `ScopeBody`. Spans never cross section boundaries.

Transformers and text width, padding, and truncation apply to complete body columns and therefore do not accept a `Scope`.

## Value display

A body cell selects its displayed value from a transformer result, default string conversion, and the placeholder, in that order. Spans are then determined from the displayed value before format-specific escaping and markup are applied.

This order defines which value wins, not the exact evaluation order. Do not rely on the order in which side effects from `String()` methods or transformers occur.

### String conversion

When a transformer does not supply a displayed value, an `any` value is converted as follows.

| Value                                                  | Displayed string                                                |
| ------------------------------------------------------ | --------------------------------------------------------------- |
| `string`, `bool`, integers, and floating-point numbers | Standard string representation                                  |
| `[]byte` and named byte slices                         | Bytes interpreted as a string                                   |
| `fmt.Stringer`                                         | Result of `String()`                                            |
| `error`                                                | Result of `Error()`                                             |
| Other slices and arrays                                | `[list N item(s)]`                                              |
| Maps                                                   | `{map N key(s)}`                                                |
| Structs                                                | `{struct N field(s)}`                                           |
| Pointers and interfaces                                | `nil`, or recursively apply the same rule to the concrete value |
| Other values                                           | `fmt.Sprint`                                                    |

Floating-point values use the shortest representation produced by `strconv`. Compound values are not expanded recursively. Convert them before rendering or use a transformer when a representation such as JSON is required.

### Transformers

`WithTransformer` receives the raw value of an existing body cell. A non-empty returned string replaces the displayed value; an empty string selects default conversion. A transformer cannot explicitly produce an empty displayed value.

A `nil` color, decoration, or `Attr` preserves the corresponding column setting.

```go
html.WithTransformer(html.Columns(2), func(v any) (string, *html.Color, *html.Decoration) {
    if n, ok := v.(int); ok && n < 0 {
        return strconv.Itoa(n), html.ColorFgRed, html.DecorationBold
    }
    return "", nil, nil
})
```

### Placeholders

`WithPlaceholder` sets the displayed value for a missing body cell. The following values are missing:

- A `nil` interface.
- A typed `nil`, including a pointer, slice, or map.
- An empty string.
- An empty slice or array.
- A trailing cell absent from a short row.
- A value for which both transformer and default conversion produce an empty string.

The default is one ASCII space, except in `html` and `csv`, where it is an empty string. Placeholders do not apply to explicitly empty header or footer labels.

Column and transformer colors, decorations, and `text.Attr` values do not apply to missing cells. Actual data equal to the placeholder string is not missing and retains its markup.

## Dynamic footers

Every format except Markdown configures a footer as follows.

```go
func WithFooter(fn func() [][]string) Option
```

`WithFooter` accepts a closure that produces footer rows.

| API      | Invocation               | Order relative to body transformers |
| -------- | ------------------------ | ----------------------------------- |
| `Table`  | At the start of `Render` | Before body transformers            |
| `Stream` | At the start of `Close`  | After every body transformation     |

A `nil` function or return value means no footer. Returned rows are emitted in order from top to bottom. Repeated calls to `Render` on the same `Table` invoke the function once per call.

For `Table`, the footer contributes to both the column count and text column widths. After a stream has emitted output, its logical column count cannot expand. A text stream also retains its initial column widths. If the column count remains unresolved, the footer can establish it.

```go
var total int
s := text.NewStream(w,
    text.WithHeader([]string{"Name", "Count"}),
    text.WithFooter(func() [][]string {
        return [][]string{{"Total", strconv.Itoa(total)}}
    }),
)
```

## Cell spans

Every format except CSV provides `WithRowspan` and `WithColspan`. `Scope` selects their target sections, except in Markdown, where spans apply only to the body.

- `WithRowspan` spans vertically adjacent equal values in selected columns. When a selected column to the left changes, selected columns to its right begin new spans.
- `WithColspan` spans horizontally adjacent equal values when both columns are selected. A vertical continuation is not eligible.
- Body cells compare strings after placeholder application; headers and footers compare configured strings. Escaping, color, and decoration are excluded, and spans never cross section boundaries.
- Spanning depends only on consecutive displayed strings. Adjacent logical groups with equal displayed strings will also span; select different columns or return distinct transformer strings when a boundary is required.
- Spans are limited to the first 64 output columns. Columns 65 and later never span. A generated index counts as one output column.
- HTML spans horizontally only across cells with equal `rowspan` values, preserving rectangles.
- `html.Table` does not split a vertical run longer than 65,534 rows and may therefore emit a `rowspan` above the HTML Standard limit.

Formats represent spans as follows.

| Format     | Vertical                                                                    | Horizontal                                             |
| ---------- | --------------------------------------------------------------------------- | ------------------------------------------------------ |
| `text`     | Connect borders and render one cell.                                        | Connect borders and render one cell.                   |
| `html`     | Use `rowspan`; streaming bodies emit empty cells.                           | Use `colspan`.                                         |
| `markdown` | Retain the first cell and empty subsequent cells; applies only to the body. | Retain the leftmost cell and empty cells to its right. |
| `backlog`  | Retain the first cell and empty subsequent cells.                           | Retain the leftmost cell and empty cells to its right. |

## Alignment

`text`, `html`, and `markdown` define `AlignDefault`, `AlignLeft`, `AlignRight`, and `AlignCenter`.

| Format     | Default                                           | Representation                               |
| ---------- | ------------------------------------------------- | -------------------------------------------- |
| `text`     | Center headers; align body and footer cells left. | Padding based on display width               |
| `html`     | Emit no `text-align`.                             | `text-align` in the cell's `style` attribute |
| `markdown` | Emit no colons and use the GFM default.           | Colons on the GFM delimiter row              |

## Indexes

`WithIndex` adds a leading column containing one-based row numbers. Its header is `#`, and its footer is empty. An index column is not created when there are no input columns.

## Captions

`text` and `html` define `CaptionDefault`, `CaptionTop`, and `CaptionBottom`. The default is below the table in text and the CSS default in HTML. HTML emits `caption` as the first child of `table` and selects its position with `caption-side`.
