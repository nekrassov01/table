# Design

This document describes the design principles and decisions for maintainers.

## Table of contents

- [Design](#design)
  - [Table of contents](#table-of-contents)
  - [Design principles](#design-principles)
    - [Shape internal structure around the output format](#shape-internal-structure-around-the-output-format)
    - [Use the same vocabulary for the same responsibility](#use-the-same-vocabulary-for-the-same-responsibility)
    - [Separate decisions from representations](#separate-decisions-from-representations)
    - [Separate state with different lifetimes](#separate-state-with-different-lifetimes)
    - [Treat clarity as a prerequisite for performance](#treat-clarity-as-a-prerequisite-for-performance)
  - [Design decisions](#design-decisions)
    - [Separate Table and Stream control flow](#separate-table-and-stream-control-flow)
    - [Chain results by pipeline stage](#chain-results-by-pipeline-stage)
    - [Limit arena to state ownership](#limit-arena-to-state-ownership)
    - [Resolve dynamic footers at render time](#resolve-dynamic-footers-at-render-time)
    - [Choose option ownership based on value semantics](#choose-option-ownership-based-on-value-semantics)
    - [Evaluate only the selected value](#evaluate-only-the-selected-value)
    - [Separate logical columns from display geometry](#separate-logical-columns-from-display-geometry)
    - [Preserve grapheme clusters when wrapping](#preserve-grapheme-clusters-when-wrapping)
    - [Resolve spans at the format-appropriate stage](#resolve-spans-at-the-format-appropriate-stage)
    - [Distinguish displayed values from attribute values](#distinguish-displayed-values-from-attribute-values)
    - [Match error retention to execution state](#match-error-retention-to-execution-state)

## Design principles

### Shape internal structure around the output format

The five output packages consume the same table data, but they cannot express the same features. They share the `Table` and `Stream` execution models and the vocabulary for corresponding pipeline stages. They do not add unnecessary stages merely to make their internal structures identical.

The following omissions are intentional.

| Omitted abstraction or feature                   | Reason                                                                                                     |
| ------------------------------------------------ | ---------------------------------------------------------------------------------------------------------- |
| A pipeline implementation shared by every format | Callbacks or type parameters would obscure format-specific rules.                                          |
| An empty `solver` in every format                | CSV has no geometry to resolve.                                                                            |
| Active rows stored in `option`                   | Mixing state with different lifetimes would blur the boundaries between `Table`, `Stream`, and the pool.   |
| Parsing arbitrary ANSI strings                   | `Attr` already provides a structured ANSI representation.                                                  |
| Permanently caching terminal width               | A permanent cache could not follow terminal resizing.                                                      |
| HTML layout options                              | Width, padding, and borders belong to the browser and CSS.                                                 |
| Visual decoration in CSV                         | CSV cannot preserve such decoration while remaining machine-readable.                                      |

### Use the same vocabulary for the same responsibility

The pipeline stages are named `config`, `compiler`, `solver`, and `painter`. A format does not introduce a different name for the same responsibility.

### Separate decisions from representations

When several formats share a rule, they share only the decision and retain their own output representation. For example, `internal/span` identifies continuation positions for adjacent cells with equal values, but it does not decide whether those positions become borders, empty cells, or HTML attributes.

### Separate state with different lifetimes

Static settings supplied through `Option` are separate from state needed only during an execution. Static `Table` and `Stream` settings do not retain execution results. String views backed by an `arena` are never stored in exported results or settings that outlive the execution. A pooled `arena` may retain only reusable slice and buffer capacity for the next execution.

### Treat clarity as a prerequisite for performance

Capacity estimation applies across the pipeline and does not belong exclusively to one stage. Each stage records accurate sizes or maxima that arise naturally during its work and passes them forward through results or state. Later stages reuse those measurements instead of rescanning values to derive the same information.

Do not add helper types without measured justification or wrappers that merely delegate work. Do not replace a short name with a longer one solely to avoid brevity.

## Design decisions

### Separate Table and Stream control flow

`Table` can inspect every row before producing output, whereas `Stream` cannot revise output already written. Their exported methods retain this difference in control flow. The shared `config`, `compiler`, `solver`, and `painter` stages use only their inputs, the preceding result, and `arena` state.

Shared stages do not receive a flag indicating whether the caller is `Table` or `Stream`, or whether lookahead is available. Such a flag would make a stage's behavior depend on the caller instead of its explicit inputs. The sequence of method calls made by `Table` and `Stream` expresses the availability of lookahead.

### Chain results by pipeline stage

`configResult` begins the result chain, `compilerResult` embeds `configResult`, and `solverResult` embeds `compilerResult`. No dedicated type exists solely to bundle arguments for `config`; `newConfig` and `resumeConfig` construct `configResult` directly. This keeps passive intermediate types out of the pipeline and lets each stage add or resolve only information within its responsibility. Formats without a `solver` pass `compilerResult` directly to `painter`.

The pipeline does not collect every stage's fields into one shared structure. Such a structure would obscure which stage resolved each field and how long it remains valid. Results are created per pass, while their reusable slices and buffers remain owned by the `arena`. This separates information flow from memory ownership.

### Limit arena to state ownership

The `arena` is not a pipeline stage. It holds reusable storage for each stage and continuation state required by the next `Stream` pass. It may supply inputs when constructing stages, but `Table` and `Stream` decide the order of `prepare`, `compile*`, `solve`, and `paint*`. Moving control flow into the `arena` would hide both execution sequences and stage inputs.

Each format defines its own `arena`. A shared type would force CSV to carry an unused `solverState` and every format to carry buffers used only for text line layout. `sync.Pool` is strictly a reuse mechanism for reducing allocations. Correctness never depends on retrieving the same `arena` or on state surviving a return to the pool. Release removes references to external values and retains only storage owned by the arena itself.

### Resolve dynamic footers at render time

`WithFooter` stores a function that returns footer rows rather than storing the rows themselves. `Table.Render` calls it once before compiling the body, while `Stream.Close` calls it once after receiving the final body row. Applications can therefore reflect body-derived aggregates in the footer produced at that moment.

The returned footer is passed to `newConfig` or `resumeConfig` and proceeds through the normal pipeline in the current pass. The function is not called while applying options, and its result is not cached between calls to `Render`. Only footer generation is deferred; column validation and compilation follow the same stages as other input sections.

### Choose option ownership based on value semantics

Options do not clone every referenced value. Values that must be independent of subsequent caller changes are owned; values read only during execution are borrowed. For example, `Columns` clones indexes to create an independently reusable selector, while headers and footer functions remain referenced until needed. Choosing ownership based on value semantics avoids unnecessary allocations.

A reusable `Option` does not mutate captured values on each application. If a value always requires the same escaping or normalization, perform that work once when constructing the option. Values such as captions and cell contents remain unprocessed until the pipeline stage responsible for their output context handles them.

### Evaluate only the selected value

`compiler` calls a configured transformer before default value conversion. A non-empty transformer result becomes the displayed value, so potentially expensive `String()` or `fmt.Sprint` work is not performed. An empty transformer result selects the default representation, and an empty default representation selects the placeholder.

Primitive slices and arrays are appended directly to arena-backed value storage while preserving their `fmt.Sprint` representation. Other values use `fmt.Append`, which provides the same representation without first creating a temporary string when the destination has reusable capacity.

Missingness is retained when the value is resolved rather than inferred later by comparing strings with the placeholder. Spans compare displayed values before markup is added. Format-specific escaping and markup remain outside shared value conversion.

### Separate logical columns from display geometry

The logical column count determined by `config` is distinct from display widths and span counts determined by `solver`. Short rows are extended to the established columns, but a wider later row never expands the table. Otherwise `Table`, which can inspect all rows, and `Stream`, which cannot revise prior output, would derive different logical column counts.

Display width is measured without changing the logical column count. `Table` can measure the complete pass, while `Stream` emits later rows within the conditions established at startup. Keeping validation separate from display adjustment prevents width optimization from changing the input contract.

When a text stream starts rendering body rows with a column whose content has zero display width, the solver reserves one display cell before freezing the geometry. This gives later values a usable wrapping boundary without changing columns after output has begun. A stream closed before its first body row has complete input and does not freeze its geometry.

### Preserve grapheme clusters when wrapping

Text column width is a wrapping boundary rather than a clipping boundary. The painter scans grapheme clusters and never divides one between physical lines. If a cluster is wider than the configured or frozen width, it remains intact and that physical line exceeds the solved column width. Splitting the cluster would corrupt its displayed value; callers that require replacement instead of wrapping can use `WithTruncate`.

### Resolve spans at the format-appropriate stage

`compiler` compares displayed values before markup and records vertical or horizontal continuation positions. Spans never cross header, body, or footer boundaries.

Formats that measure column width mark absorbed values during compilation so they do not influence width. In HTML, `compiler` records continuation positions and `solver` derives the `rowspan` and `colspan` counts from rows in the current pass. Horizontal spanning is then limited to cells with equal `rowspan` values, preserving a rectangular result.

### Distinguish displayed values from attribute values

Displayed values and attribute values are escaped or normalized for their output contexts. Strings supplied to constructors such as `NewDecoration` intentionally define markup and are not escaped. Sending both categories through the same path would sacrifice either safe displayed output or the expressiveness of caller-provided markup.

GFM separates table cells before parsing inline HTML. Markdown attribute values therefore encode vertical bars as character references, preserving both the table structure and the attribute value.

GFM normalizes code-span line endings to spaces, then removes one space from each end when both are present and the content is not entirely spaces. `escapeCode` treats source line endings as the spaces they become and adds one boundary space only when parsing would otherwise remove content. `resolveTicks` independently chooses a fence longer than every backtick run in the value.

`compiler` resolves displayed-value conversion and markup selection, and `painter` combines them in the established order. `painter` does not reinterpret value safety or decoration semantics, avoiding duplicate escaping and divergence from compilation decisions.

### Match error retention to execution state

A column-count error from `Stream.Render` rejects only the current row, so it is not retained and the stream can accept a subsequent valid row. A write error cannot be undone after output has been emitted, so `Stream` retains the first write error and returns it from later calls to both `Render` and `Close`. Errors raised during `Close` are also retained because execution has already ended.

When an HTML stream detects a footer column-count error after opening the table, `Close` omits the invalid footer but still calls `paintFooter` to close the open elements. The footer error occurred first and remains the retained result if closing also produces a write error. Without a footer error, the same closing failure is retained as a write error.
