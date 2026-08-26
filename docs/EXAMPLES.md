# Examples

This file is generated from the example data, options, and command mappings. Run `make generate` after changing them; do not edit the catalog directly.

This catalog shows the input, configuration, and output of every runnable example.

## Running examples

Use `target`, `mode`, and `data` to select the output package, API, and scenario. All three variables are optional, but selecting `data` also requires `mode`.

```sh
make example target=text mode=table data=stacked-header
```

Omitting `target` runs every output package, omitting `mode` runs both APIs, and omitting `data` runs every scenario available for the selected package. The catalog below shows the same combinations without requiring the command to be run.

Each scenario contains its shared input declaration. Each output-package section then shows the exact Option declaration and the bytes produced by `Table` and `Stream`. Identical results are shown once.

## Catalog

- [ascii](#ascii)
- [simple](#simple)
- [compact](#compact)
- [rowspan](#rowspan)
- [colspan](#colspan)
- [footer](#footer)
- [transformer](#transformer)
- [complex](#complex)
- [stacked-header](#stacked-header)
- [comma-included](#comma-included)

### ascii

Input:

````go
// SimpleData is a small instance list of plain strings.
var SimpleData = Data{
	Header: [][]string{
		{
			"INSTANCE ID",
			"INSTANCE NAME",
			"INSTANCE STATE",
		},
	},
	Body: [][]any{
		{
			"i-00000000000000000",
			"server-1",
			"running",
		},
		{
			"i-00000000000000001",
			"server-2",
			"stopped",
		},
		{
			"i-00000000000000002",
			"server-3",
			"pending",
		},
		{
			"i-00000000000000003",
			"server-4",
			"terminated",
		},
		{
			"i-00000000000000004",
			"server-5",
			"stopping",
		},
		{
			"i-00000000000000005",
			"server-6",
			"shutting-down",
		},
	},
}
````

#### text

Configuration:

````go
// TextOptionASCII configures the ASCII table example.
var TextOptionASCII = []text.Option{
	text.WithStyle(text.StyleASCII),
	text.WithHeader(SimpleData.Header...),
}
````

`Table` and `Stream` output:

````text
+---------------------+---------------+----------------+
|     INSTANCE ID     | INSTANCE NAME | INSTANCE STATE |
+---------------------+---------------+----------------+
| i-00000000000000000 | server-1      | running        |
+---------------------+---------------+----------------+
| i-00000000000000001 | server-2      | stopped        |
+---------------------+---------------+----------------+
| i-00000000000000002 | server-3      | pending        |
+---------------------+---------------+----------------+
| i-00000000000000003 | server-4      | terminated     |
+---------------------+---------------+----------------+
| i-00000000000000004 | server-5      | stopping       |
+---------------------+---------------+----------------+
| i-00000000000000005 | server-6      | shutting-down  |
+---------------------+---------------+----------------+
````

### simple

Input:

````go
// SimpleData is a small instance list of plain strings.
var SimpleData = Data{
	Header: [][]string{
		{
			"INSTANCE ID",
			"INSTANCE NAME",
			"INSTANCE STATE",
		},
	},
	Body: [][]any{
		{
			"i-00000000000000000",
			"server-1",
			"running",
		},
		{
			"i-00000000000000001",
			"server-2",
			"stopped",
		},
		{
			"i-00000000000000002",
			"server-3",
			"pending",
		},
		{
			"i-00000000000000003",
			"server-4",
			"terminated",
		},
		{
			"i-00000000000000004",
			"server-5",
			"stopping",
		},
		{
			"i-00000000000000005",
			"server-6",
			"shutting-down",
		},
	},
}
````

#### text

Configuration:

````go
// TextOptionSimple configures the simple text table example.
var TextOptionSimple = []text.Option{
	text.WithStyle(text.StyleLight),
	text.WithHeader(SimpleData.Header...),
}
````

`Table` and `Stream` output:

````text
┌─────────────────────┬───────────────┬────────────────┐
│     INSTANCE ID     │ INSTANCE NAME │ INSTANCE STATE │
╞═════════════════════╪═══════════════╪════════════════╡
│ i-00000000000000000 │ server-1      │ running        │
├─────────────────────┼───────────────┼────────────────┤
│ i-00000000000000001 │ server-2      │ stopped        │
├─────────────────────┼───────────────┼────────────────┤
│ i-00000000000000002 │ server-3      │ pending        │
├─────────────────────┼───────────────┼────────────────┤
│ i-00000000000000003 │ server-4      │ terminated     │
├─────────────────────┼───────────────┼────────────────┤
│ i-00000000000000004 │ server-5      │ stopping       │
├─────────────────────┼───────────────┼────────────────┤
│ i-00000000000000005 │ server-6      │ shutting-down  │
└─────────────────────┴───────────────┴────────────────┘
````

#### html

Configuration:

````go
// HTMLOptionSimple configures the simple HTML table example.
var HTMLOptionSimple = []html.Option{
	html.WithHeader(SimpleData.Header...),
}
````

`Table` and `Stream` output:

````html
<table>
  <thead>
    <tr>
      <th>INSTANCE ID</th>
      <th>INSTANCE NAME</th>
      <th>INSTANCE STATE</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>i-00000000000000000</td>
      <td>server-1</td>
      <td>running</td>
    </tr>
    <tr>
      <td>i-00000000000000001</td>
      <td>server-2</td>
      <td>stopped</td>
    </tr>
    <tr>
      <td>i-00000000000000002</td>
      <td>server-3</td>
      <td>pending</td>
    </tr>
    <tr>
      <td>i-00000000000000003</td>
      <td>server-4</td>
      <td>terminated</td>
    </tr>
    <tr>
      <td>i-00000000000000004</td>
      <td>server-5</td>
      <td>stopping</td>
    </tr>
    <tr>
      <td>i-00000000000000005</td>
      <td>server-6</td>
      <td>shutting-down</td>
    </tr>
  </tbody>
</table>
````

#### markdown

Configuration:

````go
// MarkdownOptionSimple configures the simple Markdown table example.
var MarkdownOptionSimple = []markdown.Option{
	markdown.WithHeader(SimpleData.Header[0]),
}
````

`Table` output:

````markdown
| INSTANCE ID         | INSTANCE NAME | INSTANCE STATE |
| ------------------- | ------------- | -------------- |
| i-00000000000000000 | server-1      | running        |
| i-00000000000000001 | server-2      | stopped        |
| i-00000000000000002 | server-3      | pending        |
| i-00000000000000003 | server-4      | terminated     |
| i-00000000000000004 | server-5      | stopping       |
| i-00000000000000005 | server-6      | shutting-down  |
````

`Stream` output:

````markdown
| INSTANCE ID | INSTANCE NAME | INSTANCE STATE |
| --- | --- | --- |
| i-00000000000000000 | server-1 | running |
| i-00000000000000001 | server-2 | stopped |
| i-00000000000000002 | server-3 | pending |
| i-00000000000000003 | server-4 | terminated |
| i-00000000000000004 | server-5 | stopping |
| i-00000000000000005 | server-6 | shutting-down |
````

#### backlog

Configuration:

````go
// BacklogOptionSimple configures the simple Backlog table example.
var BacklogOptionSimple = []backlog.Option{
	backlog.WithHeader(SimpleData.Header...),
}
````

`Table` output:

````text
|~INSTANCE ID         |~INSTANCE NAME  |~INSTANCE STATE  |
| i-00000000000000000 | server-1       | running         |
| i-00000000000000001 | server-2       | stopped         |
| i-00000000000000002 | server-3       | pending         |
| i-00000000000000003 | server-4       | terminated      |
| i-00000000000000004 | server-5       | stopping        |
| i-00000000000000005 | server-6       | shutting-down   |
````

`Stream` output:

````text
|~INSTANCE ID  |~INSTANCE NAME  |~INSTANCE STATE  |
| i-00000000000000000 | server-1 | running |
| i-00000000000000001 | server-2 | stopped |
| i-00000000000000002 | server-3 | pending |
| i-00000000000000003 | server-4 | terminated |
| i-00000000000000004 | server-5 | stopping |
| i-00000000000000005 | server-6 | shutting-down |
````

#### csv

Configuration:

````go
// CSVOptionSimple configures the simple delimiter-separated table example.
var CSVOptionSimple = []csv.Option{
	csv.WithHeader(SimpleData.Header[0]),
}
````

`Table` and `Stream` output:

````text
INSTANCE ID	INSTANCE NAME	INSTANCE STATE
i-00000000000000000	server-1	running
i-00000000000000001	server-2	stopped
i-00000000000000002	server-3	pending
i-00000000000000003	server-4	terminated
i-00000000000000004	server-5	stopping
i-00000000000000005	server-6	shutting-down
````

### compact

Input:

````go
// CompactData is a log-style table with preformatted timestamps.
var CompactData = Data{
	Header: [][]string{
		{
			"LOG TYPE",
			"DATE TIME (PREFORMATTED)",
			"MESSAGE",
			"STATUS CODE",
		},
	},
	Body: [][]any{
		{
			"accesslog",
			time.Date(2026, 5, 1, 9, 1, 15, 0, time.UTC).String(),
			"healthcheck ok",
			200,
		},
		{
			"accesslog",
			time.Date(2026, 5, 1, 9, 1, 16, 0, time.UTC).String(),
			"authentication ok",
			200,
		},
		{
			"accesslog",
			time.Date(2026, 5, 1, 9, 1, 19, 0, time.UTC).String(),
			"get resource ok",
			200,
		},
		{
			"application",
			time.Date(2026, 5, 1, 0, 19, 21, 0, time.Local).String(),
			"GET /api/v1/users/ HTTP/1.1",
			200,
		},
		{
			"application",
			time.Date(2026, 5, 1, 0, 20, 57, 0, time.Local).String(),
			"GET /api/v1/users/alice/ HTTP/1.1",
			200,
		},
		{
			"application",
			time.Date(2026, 5, 1, 1, 05, 34, 0, time.Local).String(),
			"POST /api/v1/users/bob/ HTTP/1.1",
			201,
		},
		{
			"application",
			time.Date(2026, 5, 1, 1, 07, 56, 0, time.Local).String(),
			"DELETE /api/v1/users/alice/ HTTP/1.1",
			204,
		},
	},
}
````

#### text

Configuration:

````go
// TextOptionCompact configures the compact text table example.
var TextOptionCompact = []text.Option{
	text.WithStyle(text.StyleColoredRounded),
	text.WithHeader(CompactData.Header...),
	text.WithRowspan(text.ScopeBody, text.Columns(0)),
	text.WithAlign(text.ScopeBody, text.Columns(3), text.AlignRight),
	text.WithCompact(),
	text.WithWidth(text.Columns(0), 11),
	text.WithWidth(text.Columns(1), 29),
	text.WithWidth(text.Columns(2), 36),
	text.WithWidth(text.Columns(3), 6),
}
````

`Table` and `Stream` output:

````text
╭─────────────┬───────────────────────────────┬──────────────────────────────────────┬────────╮
│  LOG TYPE   │   DATE TIME (PREFORMATTED)    │               MESSAGE                │ STATUS │
│             │                               │                                      │  CODE  │
╞═════════════╪═══════════════════════════════╪══════════════════════════════════════╪════════╡
│ accesslog   │ 2026-05-01 09:01:15 +0000 UTC │ healthcheck ok                       │    200 │
│             │ 2026-05-01 09:01:16 +0000 UTC │ authentication ok                    │    200 │
│             │ 2026-05-01 09:01:19 +0000 UTC │ get resource ok                      │    200 │
├─────────────┼───────────────────────────────┼──────────────────────────────────────┼────────┤
│ application │ 2026-05-01 00:19:21 +0000 UTC │ GET /api/v1/users/ HTTP/1.1          │    200 │
│             │ 2026-05-01 00:20:57 +0000 UTC │ GET /api/v1/users/alice/ HTTP/1.1    │    200 │
│             │ 2026-05-01 01:05:34 +0000 UTC │ POST /api/v1/users/bob/ HTTP/1.1     │    201 │
│             │ 2026-05-01 01:07:56 +0000 UTC │ DELETE /api/v1/users/alice/ HTTP/1.1 │    204 │
╰─────────────┴───────────────────────────────┴──────────────────────────────────────┴────────╯
````

### rowspan

Input:

````go
// RowspanData holds security-group rules whose leading columns repeat
// across consecutive rows.
var RowspanData = Data{
	Header: [][]string{
		{
			"INSTANCE",
			"SECURITY GROUP",
			"DIRECTION",
			"PROTOCOL",
			"FROM PORT",
			"TO PORT",
			"ADDRESS TYPE",
			"CIDR BLOCK",
		},
	},
	Body: [][]any{
		{
			"i-00000000000000000",
			"sg-10000000000000000",
			"Ingress",
			"tcp",
			22,
			22,
			"SecurityGroup",
			"sg-20000000000000000",
		},
		{
			"i-00000000000000000",
			"sg-10000000000000000",
			"Egress",
			"-1",
			0,
			0,
			"Ipv4",
			"0.0.0.0/0",
		},
		{
			"i-00000000000000000",
			"sg-10000000000000001",
			"Ingress",
			"tcp",
			443,
			443,
			"Ipv4",
			"0.0.0.0/0",
		},
		{
			"i-00000000000000000",
			"sg-10000000000000001",
			"Egress",
			"-1",
			0,
			0,
			"Ipv4",
			"0.0.0.0/0",
		},
		{
			"i-00000000000000001",
			"sg-10000000000000002",
			"Ingress",
			"tcp",
			3389,
			3389,
			"Ipv4",
			"10.1.0.0/16",
		},
		{
			"i-00000000000000001",
			"sg-10000000000000002",
			"Ingress",
			"tcp",
			0,
			65535,
			"PrefixList",
			"pl-00000000/com.amazonaws.ap-northeast-1.s3",
		},
		{
			"i-00000000000000001",
			"sg-10000000000000002",
			"Egress",
			"-1",
			0,
			0,
			"Ipv4",
			"0.0.0.0/0",
		},
		{
			"i-00000000000000002",
			"sg-10000000000000003",
			"Ingress",
			"tcp",
			443,
			443,
			"Ipv4",
			"0.0.0.0/0",
		},
		{
			"i-00000000000000002",
			"sg-10000000000000003",
			"Egress",
			"-1",
			0,
			0,
			"Ipv4",
			"0.0.0.0/0",
		},
	},
}
````

#### text

Configuration:

````go
// TextOptionRowspan configures the text table example with row spans.
var TextOptionRowspan = []text.Option{
	text.WithStyle(text.StyleColoredHeavy),
	text.WithHeader(RowspanData.Header...),
	text.WithRowspan(text.ScopeBody, text.Columns(0, 1, 2)),
	text.WithAlign(text.ScopeBody, text.Columns(4, 5), text.AlignRight),
	text.WithAutoFit(),
}
````

`Table` output:

````text
┏━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━━━━┯━━━━━━━━━━┯━━━━━━━━━━━┯━━━━━━━━━┯━━━━━━━━━━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
┃      INSTANCE       │    SECURITY GROUP    │ DIRECTION │ PROTOCOL │ FROM PORT │ TO PORT │ ADDRESS TYPE  │                 CIDR BLOCK                  ┃
┣━━━━━━━━━━━━━━━━━━━━━┿━━━━━━━━━━━━━━━━━━━━━━┿━━━━━━━━━━━┿━━━━━━━━━━┿━━━━━━━━━━━┿━━━━━━━━━┿━━━━━━━━━━━━━━━┿━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫
┃ i-00000000000000000 │ sg-10000000000000000 │ Ingress   │ tcp      │        22 │      22 │ SecurityGroup │ sg-20000000000000000                        ┃
┃                     │                      ├───────────┼──────────┼───────────┼─────────┼───────────────┼─────────────────────────────────────────────┨
┃                     │                      │ Egress    │ -1       │         0 │       0 │ Ipv4          │ 0.0.0.0/0                                   ┃
┃                     ├──────────────────────┼───────────┼──────────┼───────────┼─────────┼───────────────┼─────────────────────────────────────────────┨
┃                     │ sg-10000000000000001 │ Ingress   │ tcp      │       443 │     443 │ Ipv4          │ 0.0.0.0/0                                   ┃
┃                     │                      ├───────────┼──────────┼───────────┼─────────┼───────────────┼─────────────────────────────────────────────┨
┃                     │                      │ Egress    │ -1       │         0 │       0 │ Ipv4          │ 0.0.0.0/0                                   ┃
┠─────────────────────┼──────────────────────┼───────────┼──────────┼───────────┼─────────┼───────────────┼─────────────────────────────────────────────┨
┃ i-00000000000000001 │ sg-10000000000000002 │ Ingress   │ tcp      │      3389 │    3389 │ Ipv4          │ 10.1.0.0/16                                 ┃
┃                     │                      │           ├──────────┼───────────┼─────────┼───────────────┼─────────────────────────────────────────────┨
┃                     │                      │           │ tcp      │         0 │   65535 │ PrefixList    │ pl-00000000/com.amazonaws.ap-northeast-1.s3 ┃
┃                     │                      ├───────────┼──────────┼───────────┼─────────┼───────────────┼─────────────────────────────────────────────┨
┃                     │                      │ Egress    │ -1       │         0 │       0 │ Ipv4          │ 0.0.0.0/0                                   ┃
┠─────────────────────┼──────────────────────┼───────────┼──────────┼───────────┼─────────┼───────────────┼─────────────────────────────────────────────┨
┃ i-00000000000000002 │ sg-10000000000000003 │ Ingress   │ tcp      │       443 │     443 │ Ipv4          │ 0.0.0.0/0                                   ┃
┃                     │                      ├───────────┼──────────┼───────────┼─────────┼───────────────┼─────────────────────────────────────────────┨
┃                     │                      │ Egress    │ -1       │         0 │       0 │ Ipv4          │ 0.0.0.0/0                                   ┃
┗━━━━━━━━━━━━━━━━━━━━━┷━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━━━━━━┷━━━━━━━━━━┷━━━━━━━━━━━┷━━━━━━━━━┷━━━━━━━━━━━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
````

`Stream` output:

````text
┏━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━━━━┯━━━━━━━━━━┯━━━━━━━━━━━┯━━━━━━━━━┯━━━━━━━━━━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━┓
┃      INSTANCE       │    SECURITY GROUP    │ DIRECTION │ PROTOCOL │ FROM PORT │ TO PORT │ ADDRESS TYPE  │      CIDR BLOCK      ┃
┣━━━━━━━━━━━━━━━━━━━━━┿━━━━━━━━━━━━━━━━━━━━━━┿━━━━━━━━━━━┿━━━━━━━━━━┿━━━━━━━━━━━┿━━━━━━━━━┿━━━━━━━━━━━━━━━┿━━━━━━━━━━━━━━━━━━━━━━┫
┃ i-00000000000000000 │ sg-10000000000000000 │ Ingress   │ tcp      │        22 │      22 │ SecurityGroup │ sg-20000000000000000 ┃
┃                     │                      ├───────────┼──────────┼───────────┼─────────┼───────────────┼──────────────────────┨
┃                     │                      │ Egress    │ -1       │         0 │       0 │ Ipv4          │ 0.0.0.0/0            ┃
┃                     ├──────────────────────┼───────────┼──────────┼───────────┼─────────┼───────────────┼──────────────────────┨
┃                     │ sg-10000000000000001 │ Ingress   │ tcp      │       443 │     443 │ Ipv4          │ 0.0.0.0/0            ┃
┃                     │                      ├───────────┼──────────┼───────────┼─────────┼───────────────┼──────────────────────┨
┃                     │                      │ Egress    │ -1       │         0 │       0 │ Ipv4          │ 0.0.0.0/0            ┃
┠─────────────────────┼──────────────────────┼───────────┼──────────┼───────────┼─────────┼───────────────┼──────────────────────┨
┃ i-00000000000000001 │ sg-10000000000000002 │ Ingress   │ tcp      │      3389 │    3389 │ Ipv4          │ 10.1.0.0/16          ┃
┃                     │                      │           ├──────────┼───────────┼─────────┼───────────────┼──────────────────────┨
┃                     │                      │           │ tcp      │         0 │   65535 │ PrefixList    │ pl-00000000/com.amaz ┃
┃                     │                      │           │          │           │         │               │ onaws.ap-northeast-1 ┃
┃                     │                      │           │          │           │         │               │ .s3                  ┃
┃                     │                      ├───────────┼──────────┼───────────┼─────────┼───────────────┼──────────────────────┨
┃                     │                      │ Egress    │ -1       │         0 │       0 │ Ipv4          │ 0.0.0.0/0            ┃
┠─────────────────────┼──────────────────────┼───────────┼──────────┼───────────┼─────────┼───────────────┼──────────────────────┨
┃ i-00000000000000002 │ sg-10000000000000003 │ Ingress   │ tcp      │       443 │     443 │ Ipv4          │ 0.0.0.0/0            ┃
┃                     │                      ├───────────┼──────────┼───────────┼─────────┼───────────────┼──────────────────────┨
┃                     │                      │ Egress    │ -1       │         0 │       0 │ Ipv4          │ 0.0.0.0/0            ┃
┗━━━━━━━━━━━━━━━━━━━━━┷━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━━━━━━┷━━━━━━━━━━┷━━━━━━━━━━━┷━━━━━━━━━┷━━━━━━━━━━━━━━━┷━━━━━━━━━━━━━━━━━━━━━━┛
````

#### html

Configuration:

````go
// HTMLOptionRowspan configures the HTML table example with row spans.
var HTMLOptionRowspan = []html.Option{
	html.WithHeader(RowspanData.Header...),
	html.WithRowspan(html.ScopeBody, html.Columns(0, 1, 2)),
	html.WithAlign(html.ScopeBody, html.Columns(4, 5), html.AlignRight),
}
````

`Table` output:

````html
<table>
  <thead>
    <tr>
      <th>INSTANCE</th>
      <th>SECURITY GROUP</th>
      <th>DIRECTION</th>
      <th>PROTOCOL</th>
      <th>FROM PORT</th>
      <th>TO PORT</th>
      <th>ADDRESS TYPE</th>
      <th>CIDR BLOCK</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td rowspan="4">i-00000000000000000</td>
      <td rowspan="2">sg-10000000000000000</td>
      <td>Ingress</td>
      <td>tcp</td>
      <td style="text-align:right">22</td>
      <td style="text-align:right">22</td>
      <td>SecurityGroup</td>
      <td>sg-20000000000000000</td>
    </tr>
    <tr>
      <td>Egress</td>
      <td>-1</td>
      <td style="text-align:right">0</td>
      <td style="text-align:right">0</td>
      <td>Ipv4</td>
      <td>0.0.0.0/0</td>
    </tr>
    <tr>
      <td rowspan="2">sg-10000000000000001</td>
      <td>Ingress</td>
      <td>tcp</td>
      <td style="text-align:right">443</td>
      <td style="text-align:right">443</td>
      <td>Ipv4</td>
      <td>0.0.0.0/0</td>
    </tr>
    <tr>
      <td>Egress</td>
      <td>-1</td>
      <td style="text-align:right">0</td>
      <td style="text-align:right">0</td>
      <td>Ipv4</td>
      <td>0.0.0.0/0</td>
    </tr>
    <tr>
      <td rowspan="3">i-00000000000000001</td>
      <td rowspan="3">sg-10000000000000002</td>
      <td rowspan="2">Ingress</td>
      <td>tcp</td>
      <td style="text-align:right">3389</td>
      <td style="text-align:right">3389</td>
      <td>Ipv4</td>
      <td>10.1.0.0/16</td>
    </tr>
    <tr>
      <td>tcp</td>
      <td style="text-align:right">0</td>
      <td style="text-align:right">65535</td>
      <td>PrefixList</td>
      <td>pl-00000000/com.amazonaws.ap-northeast-1.s3</td>
    </tr>
    <tr>
      <td>Egress</td>
      <td>-1</td>
      <td style="text-align:right">0</td>
      <td style="text-align:right">0</td>
      <td>Ipv4</td>
      <td>0.0.0.0/0</td>
    </tr>
    <tr>
      <td rowspan="2">i-00000000000000002</td>
      <td rowspan="2">sg-10000000000000003</td>
      <td>Ingress</td>
      <td>tcp</td>
      <td style="text-align:right">443</td>
      <td style="text-align:right">443</td>
      <td>Ipv4</td>
      <td>0.0.0.0/0</td>
    </tr>
    <tr>
      <td>Egress</td>
      <td>-1</td>
      <td style="text-align:right">0</td>
      <td style="text-align:right">0</td>
      <td>Ipv4</td>
      <td>0.0.0.0/0</td>
    </tr>
  </tbody>
</table>
````

`Stream` output:

````html
<table>
  <thead>
    <tr>
      <th>INSTANCE</th>
      <th>SECURITY GROUP</th>
      <th>DIRECTION</th>
      <th>PROTOCOL</th>
      <th>FROM PORT</th>
      <th>TO PORT</th>
      <th>ADDRESS TYPE</th>
      <th>CIDR BLOCK</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>i-00000000000000000</td>
      <td>sg-10000000000000000</td>
      <td>Ingress</td>
      <td>tcp</td>
      <td style="text-align:right">22</td>
      <td style="text-align:right">22</td>
      <td>SecurityGroup</td>
      <td>sg-20000000000000000</td>
    </tr>
    <tr>
      <td></td>
      <td></td>
      <td>Egress</td>
      <td>-1</td>
      <td style="text-align:right">0</td>
      <td style="text-align:right">0</td>
      <td>Ipv4</td>
      <td>0.0.0.0/0</td>
    </tr>
    <tr>
      <td></td>
      <td>sg-10000000000000001</td>
      <td>Ingress</td>
      <td>tcp</td>
      <td style="text-align:right">443</td>
      <td style="text-align:right">443</td>
      <td>Ipv4</td>
      <td>0.0.0.0/0</td>
    </tr>
    <tr>
      <td></td>
      <td></td>
      <td>Egress</td>
      <td>-1</td>
      <td style="text-align:right">0</td>
      <td style="text-align:right">0</td>
      <td>Ipv4</td>
      <td>0.0.0.0/0</td>
    </tr>
    <tr>
      <td>i-00000000000000001</td>
      <td>sg-10000000000000002</td>
      <td>Ingress</td>
      <td>tcp</td>
      <td style="text-align:right">3389</td>
      <td style="text-align:right">3389</td>
      <td>Ipv4</td>
      <td>10.1.0.0/16</td>
    </tr>
    <tr>
      <td></td>
      <td></td>
      <td></td>
      <td>tcp</td>
      <td style="text-align:right">0</td>
      <td style="text-align:right">65535</td>
      <td>PrefixList</td>
      <td>pl-00000000/com.amazonaws.ap-northeast-1.s3</td>
    </tr>
    <tr>
      <td></td>
      <td></td>
      <td>Egress</td>
      <td>-1</td>
      <td style="text-align:right">0</td>
      <td style="text-align:right">0</td>
      <td>Ipv4</td>
      <td>0.0.0.0/0</td>
    </tr>
    <tr>
      <td>i-00000000000000002</td>
      <td>sg-10000000000000003</td>
      <td>Ingress</td>
      <td>tcp</td>
      <td style="text-align:right">443</td>
      <td style="text-align:right">443</td>
      <td>Ipv4</td>
      <td>0.0.0.0/0</td>
    </tr>
    <tr>
      <td></td>
      <td></td>
      <td>Egress</td>
      <td>-1</td>
      <td style="text-align:right">0</td>
      <td style="text-align:right">0</td>
      <td>Ipv4</td>
      <td>0.0.0.0/0</td>
    </tr>
  </tbody>
</table>
````

#### markdown

Configuration:

````go
// MarkdownOptionRowspan configures the Markdown table example with row spans.
var MarkdownOptionRowspan = []markdown.Option{
	markdown.WithHeader(RowspanData.Header[0]),
	markdown.WithRowspan(markdown.Columns(0, 1, 2)),
	markdown.WithAlign(markdown.Columns(4, 5), markdown.AlignRight),
}
````

`Table` output:

````markdown
| INSTANCE            | SECURITY GROUP       | DIRECTION | PROTOCOL | FROM PORT | TO PORT | ADDRESS TYPE  | CIDR BLOCK                                  |
| ------------------- | -------------------- | --------- | -------- | --------: | ------: | ------------- | ------------------------------------------- |
| i-00000000000000000 | sg-10000000000000000 | Ingress   | tcp      |        22 |      22 | SecurityGroup | sg-20000000000000000                        |
|                     |                      | Egress    | -1       |         0 |       0 | Ipv4          | 0.0.0.0/0                                   |
|                     | sg-10000000000000001 | Ingress   | tcp      |       443 |     443 | Ipv4          | 0.0.0.0/0                                   |
|                     |                      | Egress    | -1       |         0 |       0 | Ipv4          | 0.0.0.0/0                                   |
| i-00000000000000001 | sg-10000000000000002 | Ingress   | tcp      |      3389 |    3389 | Ipv4          | 10.1.0.0/16                                 |
|                     |                      |           | tcp      |         0 |   65535 | PrefixList    | pl-00000000/com.amazonaws.ap-northeast-1.s3 |
|                     |                      | Egress    | -1       |         0 |       0 | Ipv4          | 0.0.0.0/0                                   |
| i-00000000000000002 | sg-10000000000000003 | Ingress   | tcp      |       443 |     443 | Ipv4          | 0.0.0.0/0                                   |
|                     |                      | Egress    | -1       |         0 |       0 | Ipv4          | 0.0.0.0/0                                   |
````

`Stream` output:

````markdown
| INSTANCE | SECURITY GROUP | DIRECTION | PROTOCOL | FROM PORT | TO PORT | ADDRESS TYPE | CIDR BLOCK |
| --- | --- | --- | --- | ---: | ---: | --- | --- |
| i-00000000000000000 | sg-10000000000000000 | Ingress | tcp | 22 | 22 | SecurityGroup | sg-20000000000000000 |
|  |  | Egress | -1 | 0 | 0 | Ipv4 | 0.0.0.0/0 |
|  | sg-10000000000000001 | Ingress | tcp | 443 | 443 | Ipv4 | 0.0.0.0/0 |
|  |  | Egress | -1 | 0 | 0 | Ipv4 | 0.0.0.0/0 |
| i-00000000000000001 | sg-10000000000000002 | Ingress | tcp | 3389 | 3389 | Ipv4 | 10.1.0.0/16 |
|  |  |  | tcp | 0 | 65535 | PrefixList | pl-00000000/com.amazonaws.ap-northeast-1.s3 |
|  |  | Egress | -1 | 0 | 0 | Ipv4 | 0.0.0.0/0 |
| i-00000000000000002 | sg-10000000000000003 | Ingress | tcp | 443 | 443 | Ipv4 | 0.0.0.0/0 |
|  |  | Egress | -1 | 0 | 0 | Ipv4 | 0.0.0.0/0 |
````

#### backlog

Configuration:

````go
// BacklogOptionRowspan configures the Backlog table example with row spans.
var BacklogOptionRowspan = []backlog.Option{
	backlog.WithHeader(RowspanData.Header...),
	backlog.WithRowspan(backlog.ScopeBody, backlog.Columns(0, 1, 2)),
}
````

`Table` output:

````text
|~INSTANCE            |~SECURITY GROUP       |~DIRECTION  |~PROTOCOL  |~FROM PORT  |~TO PORT  |~ADDRESS TYPE  |~CIDR BLOCK                                  |
| i-00000000000000000 | sg-10000000000000000 | Ingress    | tcp       | 22         | 22       | SecurityGroup | sg-20000000000000000                        |
|                     |                      | Egress     | -1        | 0          | 0        | Ipv4          | 0.0.0.0/0                                   |
|                     | sg-10000000000000001 | Ingress    | tcp       | 443        | 443      | Ipv4          | 0.0.0.0/0                                   |
|                     |                      | Egress     | -1        | 0          | 0        | Ipv4          | 0.0.0.0/0                                   |
| i-00000000000000001 | sg-10000000000000002 | Ingress    | tcp       | 3389       | 3389     | Ipv4          | 10.1.0.0/16                                 |
|                     |                      |            | tcp       | 0          | 65535    | PrefixList    | pl-00000000/com.amazonaws.ap-northeast-1.s3 |
|                     |                      | Egress     | -1        | 0          | 0        | Ipv4          | 0.0.0.0/0                                   |
| i-00000000000000002 | sg-10000000000000003 | Ingress    | tcp       | 443        | 443      | Ipv4          | 0.0.0.0/0                                   |
|                     |                      | Egress     | -1        | 0          | 0        | Ipv4          | 0.0.0.0/0                                   |
````

`Stream` output:

````text
|~INSTANCE  |~SECURITY GROUP  |~DIRECTION  |~PROTOCOL  |~FROM PORT  |~TO PORT  |~ADDRESS TYPE  |~CIDR BLOCK  |
| i-00000000000000000 | sg-10000000000000000 | Ingress | tcp | 22 | 22 | SecurityGroup | sg-20000000000000000 |
|  |  | Egress | -1 | 0 | 0 | Ipv4 | 0.0.0.0/0 |
|  | sg-10000000000000001 | Ingress | tcp | 443 | 443 | Ipv4 | 0.0.0.0/0 |
|  |  | Egress | -1 | 0 | 0 | Ipv4 | 0.0.0.0/0 |
| i-00000000000000001 | sg-10000000000000002 | Ingress | tcp | 3389 | 3389 | Ipv4 | 10.1.0.0/16 |
|  |  |  | tcp | 0 | 65535 | PrefixList | pl-00000000/com.amazonaws.ap-northeast-1.s3 |
|  |  | Egress | -1 | 0 | 0 | Ipv4 | 0.0.0.0/0 |
| i-00000000000000002 | sg-10000000000000003 | Ingress | tcp | 443 | 443 | Ipv4 | 0.0.0.0/0 |
|  |  | Egress | -1 | 0 | 0 | Ipv4 | 0.0.0.0/0 |
````

### colspan

Input:

````go
// ColspanData contains adjacent repeated values for column-span examples.
var ColspanData = Data{
	Header: [][]string{
		{
			"FIRST NAME",
			"LAST NAME",
			"AGE",
			"BIRTH (PREFORMATTED)",
		},
	},
	Body: [][]any{
		{
			"John",
			"Doe",
			30,
			time.Date(1994, 5, 1, 0, 0, 0, 0, time.UTC).String(),
		},
		{
			"Jane",
			"Smith",
			25,
			time.Date(1999, 5, 1, 0, 0, 0, 0, time.UTC).String(),
		},
		{
			"Anonymous",
			"Anonymous",
			"Unknown",
			"Unknown",
		},
		{
			"Alice",
			"Johnson",
			28,
			time.Date(1996, 5, 1, 0, 0, 0, 0, time.UTC).String(),
		},
	},
}
````

#### text

Configuration:

````go
// TextOptionColspan configures the text table example with column spans.
var TextOptionColspan = []text.Option{
	text.WithStyle(text.StyleColoredLight),
	text.WithHeader(ColspanData.Header...),
	text.WithColspan(text.ScopeBody, text.Columns(0, 1, 2, 3)),
}
````

`Table` and `Stream` output:

````text
┌────────────┬───────────┬─────┬───────────────────────────────┐
│ FIRST NAME │ LAST NAME │ AGE │     BIRTH (PREFORMATTED)      │
╞════════════╪═══════════╪═════╪═══════════════════════════════╡
│ John       │ Doe       │ 30  │ 1994-05-01 00:00:00 +0000 UTC │
├────────────┼───────────┼─────┼───────────────────────────────┤
│ Jane       │ Smith     │ 25  │ 1999-05-01 00:00:00 +0000 UTC │
├────────────┴───────────┼─────┴───────────────────────────────┤
│ Anonymous              │ Unknown                             │
├────────────┬───────────┼─────┬───────────────────────────────┤
│ Alice      │ Johnson   │ 28  │ 1996-05-01 00:00:00 +0000 UTC │
└────────────┴───────────┴─────┴───────────────────────────────┘
````

#### html

Configuration:

````go
// HTMLOptionColspan configures the HTML table example with column spans.
var HTMLOptionColspan = []html.Option{
	html.WithHeader(ColspanData.Header...),
	html.WithColspan(html.ScopeBody, html.Columns(0, 1, 2, 3)),
}
````

`Table` and `Stream` output:

````html
<table>
  <thead>
    <tr>
      <th>FIRST NAME</th>
      <th>LAST NAME</th>
      <th>AGE</th>
      <th>BIRTH (PREFORMATTED)</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>John</td>
      <td>Doe</td>
      <td>30</td>
      <td>1994-05-01 00:00:00 +0000 UTC</td>
    </tr>
    <tr>
      <td>Jane</td>
      <td>Smith</td>
      <td>25</td>
      <td>1999-05-01 00:00:00 +0000 UTC</td>
    </tr>
    <tr>
      <td colspan="2">Anonymous</td>
      <td colspan="2">Unknown</td>
    </tr>
    <tr>
      <td>Alice</td>
      <td>Johnson</td>
      <td>28</td>
      <td>1996-05-01 00:00:00 +0000 UTC</td>
    </tr>
  </tbody>
</table>
````

#### markdown

Configuration:

````go
// MarkdownOptionColspan configures the Markdown table example with column spans.
var MarkdownOptionColspan = []markdown.Option{
	markdown.WithHeader(ColspanData.Header[0]),
	markdown.WithColspan(markdown.Columns(0, 1, 2, 3)),
}
````

`Table` output:

````markdown
| FIRST NAME | LAST NAME | AGE     | BIRTH (PREFORMATTED)          |
| ---------- | --------- | ------- | ----------------------------- |
| John       | Doe       | 30      | 1994-05-01 00:00:00 +0000 UTC |
| Jane       | Smith     | 25      | 1999-05-01 00:00:00 +0000 UTC |
| Anonymous  |           | Unknown |                               |
| Alice      | Johnson   | 28      | 1996-05-01 00:00:00 +0000 UTC |
````

`Stream` output:

````markdown
| FIRST NAME | LAST NAME | AGE | BIRTH (PREFORMATTED) |
| --- | --- | --- | --- |
| John | Doe | 30 | 1994-05-01 00:00:00 +0000 UTC |
| Jane | Smith | 25 | 1999-05-01 00:00:00 +0000 UTC |
| Anonymous |  | Unknown |  |
| Alice | Johnson | 28 | 1996-05-01 00:00:00 +0000 UTC |
````

#### backlog

Configuration:

````go
// BacklogOptionColspan configures the Backlog table example with column spans.
var BacklogOptionColspan = []backlog.Option{
	backlog.WithHeader(ColspanData.Header...),
	backlog.WithColspan(backlog.ScopeBody, backlog.Columns(0, 1, 2, 3)),
}
````

`Table` output:

````text
|~FIRST NAME  |~LAST NAME  |~AGE     |~BIRTH (PREFORMATTED)          |
| John        | Doe        | 30      | 1994-05-01 00:00:00 +0000 UTC |
| Jane        | Smith      | 25      | 1999-05-01 00:00:00 +0000 UTC |
| Anonymous   |            | Unknown |                               |
| Alice       | Johnson    | 28      | 1996-05-01 00:00:00 +0000 UTC |
````

`Stream` output:

````text
|~FIRST NAME  |~LAST NAME  |~AGE  |~BIRTH (PREFORMATTED)  |
| John | Doe | 30 | 1994-05-01 00:00:00 +0000 UTC |
| Jane | Smith | 25 | 1999-05-01 00:00:00 +0000 UTC |
| Anonymous |  | Unknown |  |
| Alice | Johnson | 28 | 1996-05-01 00:00:00 +0000 UTC |
````

### footer

Input:

````go
// FooterData is a wide table with a calculated footer and CJK values.
var FooterData = newFooterData()

func newFooterData() Data {
	const (
		labelColumns = 4
		valueColumns = 10
	)
	type footerState struct {
		bodyRows int
		totals   [valueColumns]int
		rows     [][]string
	}
	data := Data{
		Header: [][]string{
			{
				"TEAM",
				"CLASS",
				"NAME",
				"BIRTH",
				"ATB GAUGE",
				"HIT POINT",
				"SKILL POINT",
				"SPELL POINT",
				"LIFE POINT",
				"STRENGTH",
				"STAMINA",
				"DEXTERITY",
				"MAGIC",
				"SPEED",
			},
		},
		Body: [][]any{
			{
				"蜀",
				"君主",
				"劉備 玄徳",
				161,
				45,
				2800,
				350,
				280,
				5,
				62,
				68,
				55,
				58,
				50,
			},
			{
				"蜀",
				"軍神",
				"関羽 雲長",
				160,
				38,
				3500,
				420,
				150,
				4,
				95,
				88,
				72,
				35,
				55,
			},
			{
				"蜀",
				"猛将",
				"張飛 翼徳",
				167,
				52,
				3200,
				380,
				80,
				3,
				97,
				82,
				40,
				18,
				62,
			},
			{
				"蜀",
				"軍師",
				"諸葛亮 孔明",
				181,
				30,
				1800,
				280,
				580,
				3,
				25,
				38,
				65,
				99,
				45,
			},
		},
	}
	var cached atomic.Pointer[footerState]
	data.Footer = func() [][]string {
		var totals [valueColumns]int
		for column := range valueColumns {
			for _, row := range data.Body {
				totals[column] += row[labelColumns+column].(int)
			}
		}
		bodyRows := len(data.Body)
		if state := cached.Load(); state != nil && state.bodyRows == bodyRows && state.totals == totals {
			return state.rows
		}
		footer := make([]string, labelColumns+valueColumns)
		for column := range labelColumns {
			footer[column] = "平均"
		}
		for column, total := range totals {
			average := float64(total) / float64(bodyRows)
			footer[labelColumns+column] = strconv.FormatFloat(average, 'f', -1, 64)
		}
		state := &footerState{
			bodyRows: bodyRows,
			totals:   totals,
			rows:     [][]string{footer},
		}
		cached.Store(state)
		return state.rows
	}
	return data
}
````

#### text

Configuration:

````go
// TextOptionFooter configures the text table example with a footer.
var TextOptionFooter = []text.Option{
	text.WithStyle(text.StyleColoredLight),
	text.WithHeader(FooterData.Header...),
	text.WithFooter(FooterData.Footer),
	text.WithRowspan(text.ScopeBody, text.Columns(0)),
	text.WithColspan(text.ScopeFooter, text.Columns(0, 1, 2, 3)),
	text.WithAlign(text.ScopeBody, text.Columns(3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13), text.AlignRight),
	text.WithAlign(text.ScopeFooter, text.Columns(4, 5, 6, 7, 8, 9, 10, 11, 12, 13), text.AlignRight),
	text.WithAlign(text.ScopeFooter, text.Columns(0), text.AlignCenter),
	text.WithAutoFit(),
}
````

`Table` output:

````text
┌──────┬───────┬─────────────┬───────┬───────────┬───────────┬─────────────┬─────────────┬────────────┬──────────┬─────────┬───────────┬───────┬───────┐
│ TEAM │ CLASS │    NAME     │ BIRTH │ ATB GAUGE │ HIT POINT │ SKILL POINT │ SPELL POINT │ LIFE POINT │ STRENGTH │ STAMINA │ DEXTERITY │ MAGIC │ SPEED │
╞══════╪═══════╪═════════════╪═══════╪═══════════╪═══════════╪═════════════╪═════════════╪════════════╪══════════╪═════════╪═══════════╪═══════╪═══════╡
│ 蜀   │ 君主  │ 劉備 玄徳   │   161 │        45 │      2800 │         350 │         280 │          5 │       62 │      68 │        55 │    58 │    50 │
│      ├───────┼─────────────┼───────┼───────────┼───────────┼─────────────┼─────────────┼────────────┼──────────┼─────────┼───────────┼───────┼───────┤
│      │ 軍神  │ 関羽 雲長   │   160 │        38 │      3500 │         420 │         150 │          4 │       95 │      88 │        72 │    35 │    55 │
│      ├───────┼─────────────┼───────┼───────────┼───────────┼─────────────┼─────────────┼────────────┼──────────┼─────────┼───────────┼───────┼───────┤
│      │ 猛将  │ 張飛 翼徳   │   167 │        52 │      3200 │         380 │          80 │          3 │       97 │      82 │        40 │    18 │    62 │
│      ├───────┼─────────────┼───────┼───────────┼───────────┼─────────────┼─────────────┼────────────┼──────────┼─────────┼───────────┼───────┼───────┤
│      │ 軍師  │ 諸葛亮 孔明 │   181 │        30 │      1800 │         280 │         580 │          3 │       25 │      38 │        65 │    99 │    45 │
╞══════╧═══════╧═════════════╧═══════╪═══════════╪═══════════╪═════════════╪═════════════╪════════════╪══════════╪═════════╪═══════════╪═══════╪═══════╡
│                平均                │     41.25 │      2825 │       357.5 │       272.5 │       3.75 │    69.75 │      69 │        58 │  52.5 │    53 │
└────────────────────────────────────┴───────────┴───────────┴─────────────┴─────────────┴────────────┴──────────┴─────────┴───────────┴───────┴───────┘
````

`Stream` output:

````text
┌──────┬───────┬───────────┬───────┬───────────┬───────────┬─────────────┬─────────────┬────────────┬──────────┬─────────┬───────────┬───────┬───────┐
│ TEAM │ CLASS │   NAME    │ BIRTH │ ATB GAUGE │ HIT POINT │ SKILL POINT │ SPELL POINT │ LIFE POINT │ STRENGTH │ STAMINA │ DEXTERITY │ MAGIC │ SPEED │
╞══════╪═══════╪═══════════╪═══════╪═══════════╪═══════════╪═════════════╪═════════════╪════════════╪══════════╪═════════╪═══════════╪═══════╪═══════╡
│ 蜀   │ 君主  │ 劉備 玄徳 │   161 │        45 │      2800 │         350 │         280 │          5 │       62 │      68 │        55 │    58 │    50 │
│      ├───────┼───────────┼───────┼───────────┼───────────┼─────────────┼─────────────┼────────────┼──────────┼─────────┼───────────┼───────┼───────┤
│      │ 軍神  │ 関羽 雲長 │   160 │        38 │      3500 │         420 │         150 │          4 │       95 │      88 │        72 │    35 │    55 │
│      ├───────┼───────────┼───────┼───────────┼───────────┼─────────────┼─────────────┼────────────┼──────────┼─────────┼───────────┼───────┼───────┤
│      │ 猛将  │ 張飛 翼徳 │   167 │        52 │      3200 │         380 │          80 │          3 │       97 │      82 │        40 │    18 │    62 │
│      ├───────┼───────────┼───────┼───────────┼───────────┼─────────────┼─────────────┼────────────┼──────────┼─────────┼───────────┼───────┼───────┤
│      │ 軍師  │ 諸葛亮 孔 │   181 │        30 │      1800 │         280 │         580 │          3 │       25 │      38 │        65 │    99 │    45 │
│      │       │ 明        │       │           │           │             │             │            │          │         │           │       │       │
╞══════╧═══════╧═══════════╧═══════╪═══════════╪═══════════╪═════════════╪═════════════╪════════════╪══════════╪═════════╪═══════════╪═══════╪═══════╡
│               平均               │     41.25 │      2825 │       357.5 │       272.5 │       3.75 │    69.75 │      69 │        58 │  52.5 │    53 │
└──────────────────────────────────┴───────────┴───────────┴─────────────┴─────────────┴────────────┴──────────┴─────────┴───────────┴───────┴───────┘
````

#### html

Configuration:

````go
// HTMLOptionFooter configures the HTML table example with a footer.
var HTMLOptionFooter = []html.Option{
	html.WithHeader(FooterData.Header...),
	html.WithFooter(FooterData.Footer),
	html.WithRowspan(html.ScopeBody, html.Columns(0)),
	html.WithColspan(html.ScopeFooter, html.Columns(0, 1, 2, 3)),
	html.WithAlign(html.ScopeBody, html.Columns(3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13), html.AlignRight),
	html.WithAlign(html.ScopeFooter, html.Columns(4, 5, 6, 7, 8, 9, 10, 11, 12, 13), html.AlignRight),
	html.WithAlign(html.ScopeFooter, html.Columns(0), html.AlignCenter),
}
````

`Table` output:

````html
<table>
  <thead>
    <tr>
      <th>TEAM</th>
      <th>CLASS</th>
      <th>NAME</th>
      <th>BIRTH</th>
      <th>ATB GAUGE</th>
      <th>HIT POINT</th>
      <th>SKILL POINT</th>
      <th>SPELL POINT</th>
      <th>LIFE POINT</th>
      <th>STRENGTH</th>
      <th>STAMINA</th>
      <th>DEXTERITY</th>
      <th>MAGIC</th>
      <th>SPEED</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td rowspan="4">蜀</td>
      <td>君主</td>
      <td>劉備 玄徳</td>
      <td style="text-align:right">161</td>
      <td style="text-align:right">45</td>
      <td style="text-align:right">2800</td>
      <td style="text-align:right">350</td>
      <td style="text-align:right">280</td>
      <td style="text-align:right">5</td>
      <td style="text-align:right">62</td>
      <td style="text-align:right">68</td>
      <td style="text-align:right">55</td>
      <td style="text-align:right">58</td>
      <td style="text-align:right">50</td>
    </tr>
    <tr>
      <td>軍神</td>
      <td>関羽 雲長</td>
      <td style="text-align:right">160</td>
      <td style="text-align:right">38</td>
      <td style="text-align:right">3500</td>
      <td style="text-align:right">420</td>
      <td style="text-align:right">150</td>
      <td style="text-align:right">4</td>
      <td style="text-align:right">95</td>
      <td style="text-align:right">88</td>
      <td style="text-align:right">72</td>
      <td style="text-align:right">35</td>
      <td style="text-align:right">55</td>
    </tr>
    <tr>
      <td>猛将</td>
      <td>張飛 翼徳</td>
      <td style="text-align:right">167</td>
      <td style="text-align:right">52</td>
      <td style="text-align:right">3200</td>
      <td style="text-align:right">380</td>
      <td style="text-align:right">80</td>
      <td style="text-align:right">3</td>
      <td style="text-align:right">97</td>
      <td style="text-align:right">82</td>
      <td style="text-align:right">40</td>
      <td style="text-align:right">18</td>
      <td style="text-align:right">62</td>
    </tr>
    <tr>
      <td>軍師</td>
      <td>諸葛亮 孔明</td>
      <td style="text-align:right">181</td>
      <td style="text-align:right">30</td>
      <td style="text-align:right">1800</td>
      <td style="text-align:right">280</td>
      <td style="text-align:right">580</td>
      <td style="text-align:right">3</td>
      <td style="text-align:right">25</td>
      <td style="text-align:right">38</td>
      <td style="text-align:right">65</td>
      <td style="text-align:right">99</td>
      <td style="text-align:right">45</td>
    </tr>
  </tbody>
  <tfoot>
    <tr>
      <td style="text-align:center" colspan="4">平均</td>
      <td style="text-align:right">41.25</td>
      <td style="text-align:right">2825</td>
      <td style="text-align:right">357.5</td>
      <td style="text-align:right">272.5</td>
      <td style="text-align:right">3.75</td>
      <td style="text-align:right">69.75</td>
      <td style="text-align:right">69</td>
      <td style="text-align:right">58</td>
      <td style="text-align:right">52.5</td>
      <td style="text-align:right">53</td>
    </tr>
  </tfoot>
</table>
````

`Stream` output:

````html
<table>
  <thead>
    <tr>
      <th>TEAM</th>
      <th>CLASS</th>
      <th>NAME</th>
      <th>BIRTH</th>
      <th>ATB GAUGE</th>
      <th>HIT POINT</th>
      <th>SKILL POINT</th>
      <th>SPELL POINT</th>
      <th>LIFE POINT</th>
      <th>STRENGTH</th>
      <th>STAMINA</th>
      <th>DEXTERITY</th>
      <th>MAGIC</th>
      <th>SPEED</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>蜀</td>
      <td>君主</td>
      <td>劉備 玄徳</td>
      <td style="text-align:right">161</td>
      <td style="text-align:right">45</td>
      <td style="text-align:right">2800</td>
      <td style="text-align:right">350</td>
      <td style="text-align:right">280</td>
      <td style="text-align:right">5</td>
      <td style="text-align:right">62</td>
      <td style="text-align:right">68</td>
      <td style="text-align:right">55</td>
      <td style="text-align:right">58</td>
      <td style="text-align:right">50</td>
    </tr>
    <tr>
      <td></td>
      <td>軍神</td>
      <td>関羽 雲長</td>
      <td style="text-align:right">160</td>
      <td style="text-align:right">38</td>
      <td style="text-align:right">3500</td>
      <td style="text-align:right">420</td>
      <td style="text-align:right">150</td>
      <td style="text-align:right">4</td>
      <td style="text-align:right">95</td>
      <td style="text-align:right">88</td>
      <td style="text-align:right">72</td>
      <td style="text-align:right">35</td>
      <td style="text-align:right">55</td>
    </tr>
    <tr>
      <td></td>
      <td>猛将</td>
      <td>張飛 翼徳</td>
      <td style="text-align:right">167</td>
      <td style="text-align:right">52</td>
      <td style="text-align:right">3200</td>
      <td style="text-align:right">380</td>
      <td style="text-align:right">80</td>
      <td style="text-align:right">3</td>
      <td style="text-align:right">97</td>
      <td style="text-align:right">82</td>
      <td style="text-align:right">40</td>
      <td style="text-align:right">18</td>
      <td style="text-align:right">62</td>
    </tr>
    <tr>
      <td></td>
      <td>軍師</td>
      <td>諸葛亮 孔明</td>
      <td style="text-align:right">181</td>
      <td style="text-align:right">30</td>
      <td style="text-align:right">1800</td>
      <td style="text-align:right">280</td>
      <td style="text-align:right">580</td>
      <td style="text-align:right">3</td>
      <td style="text-align:right">25</td>
      <td style="text-align:right">38</td>
      <td style="text-align:right">65</td>
      <td style="text-align:right">99</td>
      <td style="text-align:right">45</td>
    </tr>
  </tbody>
  <tfoot>
    <tr>
      <td style="text-align:center" colspan="4">平均</td>
      <td style="text-align:right">41.25</td>
      <td style="text-align:right">2825</td>
      <td style="text-align:right">357.5</td>
      <td style="text-align:right">272.5</td>
      <td style="text-align:right">3.75</td>
      <td style="text-align:right">69.75</td>
      <td style="text-align:right">69</td>
      <td style="text-align:right">58</td>
      <td style="text-align:right">52.5</td>
      <td style="text-align:right">53</td>
    </tr>
  </tfoot>
</table>
````

#### backlog

Configuration:

````go
// BacklogOptionFooter configures the Backlog table example with a footer.
var BacklogOptionFooter = []backlog.Option{
	backlog.WithHeader(FooterData.Header...),
	backlog.WithFooter(FooterData.Footer),
	backlog.WithRowspan(backlog.ScopeBody, backlog.Columns(0)),
	backlog.WithColspan(backlog.ScopeFooter, backlog.Columns(0, 1, 2, 3)),
}
````

`Table` output:

````text
|~TEAM  |~CLASS  |~NAME        |~BIRTH  |~ATB GAUGE  |~HIT POINT  |~SKILL POINT  |~SPELL POINT  |~LIFE POINT  |~STRENGTH  |~STAMINA  |~DEXTERITY  |~MAGIC  |~SPEED  |
| 蜀    | 君主   | 劉備 玄徳   | 161    | 45         | 2800       | 350          | 280          | 5           | 62        | 68       | 55         | 58     | 50     |
|       | 軍神   | 関羽 雲長   | 160    | 38         | 3500       | 420          | 150          | 4           | 95        | 88       | 72         | 35     | 55     |
|       | 猛将   | 張飛 翼徳   | 167    | 52         | 3200       | 380          | 80           | 3           | 97        | 82       | 40         | 18     | 62     |
|       | 軍師   | 諸葛亮 孔明 | 181    | 30         | 1800       | 280          | 580          | 3           | 25        | 38       | 65         | 99     | 45     |
|~平均  |~       |~            |~       |~41.25      |~2825       |~357.5        |~272.5        |~3.75        |~69.75     |~69       |~58         |~52.5   |~53     |
````

`Stream` output:

````text
|~TEAM  |~CLASS  |~NAME  |~BIRTH  |~ATB GAUGE  |~HIT POINT  |~SKILL POINT  |~SPELL POINT  |~LIFE POINT  |~STRENGTH  |~STAMINA  |~DEXTERITY  |~MAGIC  |~SPEED  |
| 蜀 | 君主 | 劉備 玄徳 | 161 | 45 | 2800 | 350 | 280 | 5 | 62 | 68 | 55 | 58 | 50 |
|  | 軍神 | 関羽 雲長 | 160 | 38 | 3500 | 420 | 150 | 4 | 95 | 88 | 72 | 35 | 55 |
|  | 猛将 | 張飛 翼徳 | 167 | 52 | 3200 | 380 | 80 | 3 | 97 | 82 | 40 | 18 | 62 |
|  | 軍師 | 諸葛亮 孔明 | 181 | 30 | 1800 | 280 | 580 | 3 | 25 | 38 | 65 | 99 | 45 |
|~平均  |~  |~  |~  |~41.25  |~2825  |~357.5  |~272.5  |~3.75  |~69.75  |~69  |~58  |~52.5  |~53  |
````

#### csv

Configuration:

````go
// CSVOptionFooter configures the delimiter-separated table example with a footer.
var CSVOptionFooter = []csv.Option{
	csv.WithHeader(FooterData.Header[0]),
	csv.WithFooter(FooterData.Footer),
}
````

`Table` and `Stream` output:

````text
TEAM	CLASS	NAME	BIRTH	ATB GAUGE	HIT POINT	SKILL POINT	SPELL POINT	LIFE POINT	STRENGTH	STAMINA	DEXTERITY	MAGIC	SPEED
蜀	君主	劉備 玄徳	161	45	2800	350	280	5	62	68	55	58	50
蜀	軍神	関羽 雲長	160	38	3500	420	150	4	95	88	72	35	55
蜀	猛将	張飛 翼徳	167	52	3200	380	80	3	97	82	40	18	62
蜀	軍師	諸葛亮 孔明	181	30	1800	280	580	3	25	38	65	99	45
平均	平均	平均	平均	41.25	2825	357.5	272.5	3.75	69.75	69	58	52.5	53
````

### transformer

Input:

````go
// FooterData is a wide table with a calculated footer and CJK values.
var FooterData = newFooterData()

func newFooterData() Data {
	const (
		labelColumns = 4
		valueColumns = 10
	)
	type footerState struct {
		bodyRows int
		totals   [valueColumns]int
		rows     [][]string
	}
	data := Data{
		Header: [][]string{
			{
				"TEAM",
				"CLASS",
				"NAME",
				"BIRTH",
				"ATB GAUGE",
				"HIT POINT",
				"SKILL POINT",
				"SPELL POINT",
				"LIFE POINT",
				"STRENGTH",
				"STAMINA",
				"DEXTERITY",
				"MAGIC",
				"SPEED",
			},
		},
		Body: [][]any{
			{
				"蜀",
				"君主",
				"劉備 玄徳",
				161,
				45,
				2800,
				350,
				280,
				5,
				62,
				68,
				55,
				58,
				50,
			},
			{
				"蜀",
				"軍神",
				"関羽 雲長",
				160,
				38,
				3500,
				420,
				150,
				4,
				95,
				88,
				72,
				35,
				55,
			},
			{
				"蜀",
				"猛将",
				"張飛 翼徳",
				167,
				52,
				3200,
				380,
				80,
				3,
				97,
				82,
				40,
				18,
				62,
			},
			{
				"蜀",
				"軍師",
				"諸葛亮 孔明",
				181,
				30,
				1800,
				280,
				580,
				3,
				25,
				38,
				65,
				99,
				45,
			},
		},
	}
	var cached atomic.Pointer[footerState]
	data.Footer = func() [][]string {
		var totals [valueColumns]int
		for column := range valueColumns {
			for _, row := range data.Body {
				totals[column] += row[labelColumns+column].(int)
			}
		}
		bodyRows := len(data.Body)
		if state := cached.Load(); state != nil && state.bodyRows == bodyRows && state.totals == totals {
			return state.rows
		}
		footer := make([]string, labelColumns+valueColumns)
		for column := range labelColumns {
			footer[column] = "平均"
		}
		for column, total := range totals {
			average := float64(total) / float64(bodyRows)
			footer[labelColumns+column] = strconv.FormatFloat(average, 'f', -1, 64)
		}
		state := &footerState{
			bodyRows: bodyRows,
			totals:   totals,
			rows:     [][]string{footer},
		}
		cached.Store(state)
		return state.rows
	}
	return data
}
````

#### text

Configuration:

````go
// TextOptionTransformer configures the transformed text table example.
var TextOptionTransformer = []text.Option{
	text.WithStyle(text.StyleColoredLight),
	text.WithHeader(FooterData.Header...),
	text.WithFooter(FooterData.Footer),
	text.WithRowspan(text.ScopeBody, text.Columns(0)),
	text.WithColspan(text.ScopeFooter, text.Columns(0, 1, 2, 3)),
	text.WithAlign(text.ScopeBody, text.Columns(3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13), text.AlignRight),
	text.WithAlign(text.ScopeFooter, text.Columns(4, 5, 6, 7, 8, 9, 10, 11, 12, 13), text.AlignRight),
	text.WithAlign(text.ScopeFooter, text.Columns(0), text.AlignCenter),
	text.WithAutoFit(),
	text.WithTransformer(text.Columns(5), func(v any) (string, *text.Attr) {
		n, ok := v.(int)
		if !ok {
			return "", nil
		}
		if n >= 3000 {
			return fmt.Sprintf("*%d", n), textFgRedBold
		}
		return "", nil
	}),
	text.WithTransformer(text.Columns(9), func(v any) (string, *text.Attr) {
		if n, ok := v.(int); ok && n >= 90 {
			return fmt.Sprintf("*%d", n), textFgYellowBold
		}
		return "", nil
	}),
	text.WithTransformer(text.Columns(13), func(v any) (string, *text.Attr) {
		if n, ok := v.(int); ok && n >= 60 {
			return fmt.Sprintf("*%d", n), textFgGreenBold
		}
		return "", nil
	}),
}
````

`Table` output:

````text
┌──────┬───────┬─────────────┬───────┬───────────┬───────────┬─────────────┬─────────────┬────────────┬──────────┬─────────┬───────────┬───────┬───────┐
│ TEAM │ CLASS │    NAME     │ BIRTH │ ATB GAUGE │ HIT POINT │ SKILL POINT │ SPELL POINT │ LIFE POINT │ STRENGTH │ STAMINA │ DEXTERITY │ MAGIC │ SPEED │
╞══════╪═══════╪═════════════╪═══════╪═══════════╪═══════════╪═════════════╪═════════════╪════════════╪══════════╪═════════╪═══════════╪═══════╪═══════╡
│ 蜀   │ 君主  │ 劉備 玄徳   │   161 │        45 │      2800 │         350 │         280 │          5 │       62 │      68 │        55 │    58 │    50 │
│      ├───────┼─────────────┼───────┼───────────┼───────────┼─────────────┼─────────────┼────────────┼──────────┼─────────┼───────────┼───────┼───────┤
│      │ 軍神  │ 関羽 雲長   │   160 │        38 │     *3500 │         420 │         150 │          4 │      *95 │      88 │        72 │    35 │    55 │
│      ├───────┼─────────────┼───────┼───────────┼───────────┼─────────────┼─────────────┼────────────┼──────────┼─────────┼───────────┼───────┼───────┤
│      │ 猛将  │ 張飛 翼徳   │   167 │        52 │     *3200 │         380 │          80 │          3 │      *97 │      82 │        40 │    18 │   *62 │
│      ├───────┼─────────────┼───────┼───────────┼───────────┼─────────────┼─────────────┼────────────┼──────────┼─────────┼───────────┼───────┼───────┤
│      │ 軍師  │ 諸葛亮 孔明 │   181 │        30 │      1800 │         280 │         580 │          3 │       25 │      38 │        65 │    99 │    45 │
╞══════╧═══════╧═════════════╧═══════╪═══════════╪═══════════╪═════════════╪═════════════╪════════════╪══════════╪═════════╪═══════════╪═══════╪═══════╡
│                平均                │     41.25 │      2825 │       357.5 │       272.5 │       3.75 │    69.75 │      69 │        58 │  52.5 │    53 │
└────────────────────────────────────┴───────────┴───────────┴─────────────┴─────────────┴────────────┴──────────┴─────────┴───────────┴───────┴───────┘
````

`Stream` output:

````text
┌──────┬───────┬───────────┬───────┬───────────┬───────────┬─────────────┬─────────────┬────────────┬──────────┬─────────┬───────────┬───────┬───────┐
│ TEAM │ CLASS │   NAME    │ BIRTH │ ATB GAUGE │ HIT POINT │ SKILL POINT │ SPELL POINT │ LIFE POINT │ STRENGTH │ STAMINA │ DEXTERITY │ MAGIC │ SPEED │
╞══════╪═══════╪═══════════╪═══════╪═══════════╪═══════════╪═════════════╪═════════════╪════════════╪══════════╪═════════╪═══════════╪═══════╪═══════╡
│ 蜀   │ 君主  │ 劉備 玄徳 │   161 │        45 │      2800 │         350 │         280 │          5 │       62 │      68 │        55 │    58 │    50 │
│      ├───────┼───────────┼───────┼───────────┼───────────┼─────────────┼─────────────┼────────────┼──────────┼─────────┼───────────┼───────┼───────┤
│      │ 軍神  │ 関羽 雲長 │   160 │        38 │     *3500 │         420 │         150 │          4 │      *95 │      88 │        72 │    35 │    55 │
│      ├───────┼───────────┼───────┼───────────┼───────────┼─────────────┼─────────────┼────────────┼──────────┼─────────┼───────────┼───────┼───────┤
│      │ 猛将  │ 張飛 翼徳 │   167 │        52 │     *3200 │         380 │          80 │          3 │      *97 │      82 │        40 │    18 │   *62 │
│      ├───────┼───────────┼───────┼───────────┼───────────┼─────────────┼─────────────┼────────────┼──────────┼─────────┼───────────┼───────┼───────┤
│      │ 軍師  │ 諸葛亮 孔 │   181 │        30 │      1800 │         280 │         580 │          3 │       25 │      38 │        65 │    99 │    45 │
│      │       │ 明        │       │           │           │             │             │            │          │         │           │       │       │
╞══════╧═══════╧═══════════╧═══════╪═══════════╪═══════════╪═════════════╪═════════════╪════════════╪══════════╪═════════╪═══════════╪═══════╪═══════╡
│               平均               │     41.25 │      2825 │       357.5 │       272.5 │       3.75 │    69.75 │      69 │        58 │  52.5 │    53 │
└──────────────────────────────────┴───────────┴───────────┴─────────────┴─────────────┴────────────┴──────────┴─────────┴───────────┴───────┴───────┘
````

#### html

Configuration:

````go
// HTMLOptionTransformer configures the transformed HTML table example.
var HTMLOptionTransformer = []html.Option{
	html.WithHeader(FooterData.Header...),
	html.WithFooter(FooterData.Footer),
	html.WithRowspan(html.ScopeBody, html.Columns(0)),
	html.WithColspan(html.ScopeFooter, html.Columns(0, 1, 2, 3)),
	html.WithAlign(html.ScopeBody, html.Columns(3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13), html.AlignRight),
	html.WithAlign(html.ScopeFooter, html.Columns(4, 5, 6, 7, 8, 9, 10, 11, 12, 13), html.AlignRight),
	html.WithAlign(html.ScopeFooter, html.Columns(0), html.AlignCenter),
	html.WithTransformer(html.Columns(5), func(v any) (string, *html.Color, *html.Decoration) {
		n, ok := v.(int)
		if !ok {
			return "", nil, nil
		}
		if n >= 3000 {
			return fmt.Sprintf("*%d", n), html.ColorFgRed, html.DecorationBold
		}
		return "", nil, nil
	}),
	html.WithTransformer(html.Columns(9), func(v any) (string, *html.Color, *html.Decoration) {
		if n, ok := v.(int); ok && n >= 90 {
			return fmt.Sprintf("*%d", n), html.ColorFgYellow, html.DecorationBold
		}
		return "", nil, nil
	}),
	html.WithTransformer(html.Columns(13), func(v any) (string, *html.Color, *html.Decoration) {
		if n, ok := v.(int); ok && n >= 60 {
			return fmt.Sprintf("*%d", n), html.ColorFgGreen, html.DecorationBold
		}
		return "", nil, nil
	}),
}
````

`Table` output:

````html
<table>
  <thead>
    <tr>
      <th>TEAM</th>
      <th>CLASS</th>
      <th>NAME</th>
      <th>BIRTH</th>
      <th>ATB GAUGE</th>
      <th>HIT POINT</th>
      <th>SKILL POINT</th>
      <th>SPELL POINT</th>
      <th>LIFE POINT</th>
      <th>STRENGTH</th>
      <th>STAMINA</th>
      <th>DEXTERITY</th>
      <th>MAGIC</th>
      <th>SPEED</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td rowspan="4">蜀</td>
      <td>君主</td>
      <td>劉備 玄徳</td>
      <td style="text-align:right">161</td>
      <td style="text-align:right">45</td>
      <td style="text-align:right">2800</td>
      <td style="text-align:right">350</td>
      <td style="text-align:right">280</td>
      <td style="text-align:right">5</td>
      <td style="text-align:right">62</td>
      <td style="text-align:right">68</td>
      <td style="text-align:right">55</td>
      <td style="text-align:right">58</td>
      <td style="text-align:right">50</td>
    </tr>
    <tr>
      <td>軍神</td>
      <td>関羽 雲長</td>
      <td style="text-align:right">160</td>
      <td style="text-align:right">38</td>
      <td style="text-align:right"><strong><span style="color:red">*3500</span></strong></td>
      <td style="text-align:right">420</td>
      <td style="text-align:right">150</td>
      <td style="text-align:right">4</td>
      <td style="text-align:right"><strong><span style="color:yellow">*95</span></strong></td>
      <td style="text-align:right">88</td>
      <td style="text-align:right">72</td>
      <td style="text-align:right">35</td>
      <td style="text-align:right">55</td>
    </tr>
    <tr>
      <td>猛将</td>
      <td>張飛 翼徳</td>
      <td style="text-align:right">167</td>
      <td style="text-align:right">52</td>
      <td style="text-align:right"><strong><span style="color:red">*3200</span></strong></td>
      <td style="text-align:right">380</td>
      <td style="text-align:right">80</td>
      <td style="text-align:right">3</td>
      <td style="text-align:right"><strong><span style="color:yellow">*97</span></strong></td>
      <td style="text-align:right">82</td>
      <td style="text-align:right">40</td>
      <td style="text-align:right">18</td>
      <td style="text-align:right"><strong><span style="color:green">*62</span></strong></td>
    </tr>
    <tr>
      <td>軍師</td>
      <td>諸葛亮 孔明</td>
      <td style="text-align:right">181</td>
      <td style="text-align:right">30</td>
      <td style="text-align:right">1800</td>
      <td style="text-align:right">280</td>
      <td style="text-align:right">580</td>
      <td style="text-align:right">3</td>
      <td style="text-align:right">25</td>
      <td style="text-align:right">38</td>
      <td style="text-align:right">65</td>
      <td style="text-align:right">99</td>
      <td style="text-align:right">45</td>
    </tr>
  </tbody>
  <tfoot>
    <tr>
      <td style="text-align:center" colspan="4">平均</td>
      <td style="text-align:right">41.25</td>
      <td style="text-align:right">2825</td>
      <td style="text-align:right">357.5</td>
      <td style="text-align:right">272.5</td>
      <td style="text-align:right">3.75</td>
      <td style="text-align:right">69.75</td>
      <td style="text-align:right">69</td>
      <td style="text-align:right">58</td>
      <td style="text-align:right">52.5</td>
      <td style="text-align:right">53</td>
    </tr>
  </tfoot>
</table>
````

`Stream` output:

````html
<table>
  <thead>
    <tr>
      <th>TEAM</th>
      <th>CLASS</th>
      <th>NAME</th>
      <th>BIRTH</th>
      <th>ATB GAUGE</th>
      <th>HIT POINT</th>
      <th>SKILL POINT</th>
      <th>SPELL POINT</th>
      <th>LIFE POINT</th>
      <th>STRENGTH</th>
      <th>STAMINA</th>
      <th>DEXTERITY</th>
      <th>MAGIC</th>
      <th>SPEED</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>蜀</td>
      <td>君主</td>
      <td>劉備 玄徳</td>
      <td style="text-align:right">161</td>
      <td style="text-align:right">45</td>
      <td style="text-align:right">2800</td>
      <td style="text-align:right">350</td>
      <td style="text-align:right">280</td>
      <td style="text-align:right">5</td>
      <td style="text-align:right">62</td>
      <td style="text-align:right">68</td>
      <td style="text-align:right">55</td>
      <td style="text-align:right">58</td>
      <td style="text-align:right">50</td>
    </tr>
    <tr>
      <td></td>
      <td>軍神</td>
      <td>関羽 雲長</td>
      <td style="text-align:right">160</td>
      <td style="text-align:right">38</td>
      <td style="text-align:right"><strong><span style="color:red">*3500</span></strong></td>
      <td style="text-align:right">420</td>
      <td style="text-align:right">150</td>
      <td style="text-align:right">4</td>
      <td style="text-align:right"><strong><span style="color:yellow">*95</span></strong></td>
      <td style="text-align:right">88</td>
      <td style="text-align:right">72</td>
      <td style="text-align:right">35</td>
      <td style="text-align:right">55</td>
    </tr>
    <tr>
      <td></td>
      <td>猛将</td>
      <td>張飛 翼徳</td>
      <td style="text-align:right">167</td>
      <td style="text-align:right">52</td>
      <td style="text-align:right"><strong><span style="color:red">*3200</span></strong></td>
      <td style="text-align:right">380</td>
      <td style="text-align:right">80</td>
      <td style="text-align:right">3</td>
      <td style="text-align:right"><strong><span style="color:yellow">*97</span></strong></td>
      <td style="text-align:right">82</td>
      <td style="text-align:right">40</td>
      <td style="text-align:right">18</td>
      <td style="text-align:right"><strong><span style="color:green">*62</span></strong></td>
    </tr>
    <tr>
      <td></td>
      <td>軍師</td>
      <td>諸葛亮 孔明</td>
      <td style="text-align:right">181</td>
      <td style="text-align:right">30</td>
      <td style="text-align:right">1800</td>
      <td style="text-align:right">280</td>
      <td style="text-align:right">580</td>
      <td style="text-align:right">3</td>
      <td style="text-align:right">25</td>
      <td style="text-align:right">38</td>
      <td style="text-align:right">65</td>
      <td style="text-align:right">99</td>
      <td style="text-align:right">45</td>
    </tr>
  </tbody>
  <tfoot>
    <tr>
      <td style="text-align:center" colspan="4">平均</td>
      <td style="text-align:right">41.25</td>
      <td style="text-align:right">2825</td>
      <td style="text-align:right">357.5</td>
      <td style="text-align:right">272.5</td>
      <td style="text-align:right">3.75</td>
      <td style="text-align:right">69.75</td>
      <td style="text-align:right">69</td>
      <td style="text-align:right">58</td>
      <td style="text-align:right">52.5</td>
      <td style="text-align:right">53</td>
    </tr>
  </tfoot>
</table>
````

#### markdown

Configuration:

````go
// MarkdownOptionTransformer configures the transformed Markdown table example.
var MarkdownOptionTransformer = []markdown.Option{
	markdown.WithHeader(FooterData.Header[0]),
	markdown.WithRowspan(markdown.Columns(0)),
	markdown.WithColspan(markdown.Columns(0, 1, 2, 3)),
	markdown.WithAlign(markdown.Columns(3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13), markdown.AlignRight),
	markdown.WithTransformer(markdown.Columns(5), func(v any) (string, *markdown.Color, *markdown.Decoration) {
		n, ok := v.(int)
		if !ok {
			return "", nil, nil
		}
		if n >= 3000 {
			return fmt.Sprintf("*%d", n), markdown.ColorFgRed, markdown.DecorationBold
		}
		return "", nil, nil
	}),
	markdown.WithTransformer(markdown.Columns(9), func(v any) (string, *markdown.Color, *markdown.Decoration) {
		if n, ok := v.(int); ok && n >= 90 {
			return fmt.Sprintf("*%d", n), markdown.ColorFgYellow, markdown.DecorationBold
		}
		return "", nil, nil
	}),
	markdown.WithTransformer(markdown.Columns(13), func(v any) (string, *markdown.Color, *markdown.Decoration) {
		if n, ok := v.(int); ok && n >= 60 {
			return fmt.Sprintf("*%d", n), markdown.ColorFgGreen, markdown.DecorationBold
		}
		return "", nil, nil
	}),
}
````

`Table` output:

````markdown
| TEAM | CLASS | NAME        | BIRTH | ATB GAUGE |                                 HIT POINT | SKILL POINT | SPELL POINT | LIFE POINT |                                   STRENGTH | STAMINA | DEXTERITY | MAGIC |                                     SPEED |
| ---- | ----- | ----------- | ----: | --------: | ----------------------------------------: | ----------: | ----------: | ---------: | -----------------------------------------: | ------: | --------: | ----: | ----------------------------------------: |
| 蜀   | 君主  | 劉備 玄徳   |   161 |        45 |                                      2800 |         350 |         280 |          5 |                                         62 |      68 |        55 |    58 |                                        50 |
|      | 軍神  | 関羽 雲長   |   160 |        38 | **<span style="color:red">\*3500</span>** |         420 |         150 |          4 | **<span style="color:yellow">\*95</span>** |      88 |        72 |    35 |                                        55 |
|      | 猛将  | 張飛 翼徳   |   167 |        52 | **<span style="color:red">\*3200</span>** |         380 |          80 |          3 | **<span style="color:yellow">\*97</span>** |      82 |        40 |    18 | **<span style="color:green">\*62</span>** |
|      | 軍師  | 諸葛亮 孔明 |   181 |        30 |                                      1800 |         280 |         580 |          3 |                                         25 |      38 |        65 |    99 |                                        45 |
````

`Stream` output:

````markdown
| TEAM | CLASS | NAME | BIRTH | ATB GAUGE | HIT POINT | SKILL POINT | SPELL POINT | LIFE POINT | STRENGTH | STAMINA | DEXTERITY | MAGIC | SPEED |
| --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| 蜀 | 君主 | 劉備 玄徳 | 161 | 45 | 2800 | 350 | 280 | 5 | 62 | 68 | 55 | 58 | 50 |
|  | 軍神 | 関羽 雲長 | 160 | 38 | **<span style="color:red">\*3500</span>** | 420 | 150 | 4 | **<span style="color:yellow">\*95</span>** | 88 | 72 | 35 | 55 |
|  | 猛将 | 張飛 翼徳 | 167 | 52 | **<span style="color:red">\*3200</span>** | 380 | 80 | 3 | **<span style="color:yellow">\*97</span>** | 82 | 40 | 18 | **<span style="color:green">\*62</span>** |
|  | 軍師 | 諸葛亮 孔明 | 181 | 30 | 1800 | 280 | 580 | 3 | 25 | 38 | 65 | 99 | 45 |
````

#### backlog

Configuration:

````go
// BacklogOptionTransformer configures the transformed Backlog table example.
var BacklogOptionTransformer = []backlog.Option{
	backlog.WithHeader(FooterData.Header...),
	backlog.WithFooter(FooterData.Footer),
	backlog.WithRowspan(backlog.ScopeBody, backlog.Columns(0)),
	backlog.WithColspan(backlog.ScopeFooter, backlog.Columns(0, 1, 2, 3)),
	backlog.WithTransformer(backlog.Columns(5), func(v any) (string, *backlog.Color, *backlog.Decoration) {
		n, ok := v.(int)
		if !ok {
			return "", nil, nil
		}
		if n >= 3000 {
			return fmt.Sprintf("*%d", n), backlog.ColorFgRed, backlog.DecorationBold
		}
		return "", nil, nil
	}),
	backlog.WithTransformer(backlog.Columns(9), func(v any) (string, *backlog.Color, *backlog.Decoration) {
		if n, ok := v.(int); ok && n >= 90 {
			return fmt.Sprintf("*%d", n), backlog.ColorFgYellow, backlog.DecorationBold
		}
		return "", nil, nil
	}),
	backlog.WithTransformer(backlog.Columns(13), func(v any) (string, *backlog.Color, *backlog.Decoration) {
		if n, ok := v.(int); ok && n >= 60 {
			return fmt.Sprintf("*%d", n), backlog.ColorFgGreen, backlog.DecorationBold
		}
		return "", nil, nil
	}),
}
````

`Table` output:

````text
|~TEAM  |~CLASS  |~NAME        |~BIRTH  |~ATB GAUGE  |~HIT POINT              |~SKILL POINT  |~SPELL POINT  |~LIFE POINT  |~STRENGTH                |~STAMINA  |~DEXTERITY  |~MAGIC  |~SPEED                  |
| 蜀    | 君主   | 劉備 玄徳   | 161    | 45         | 2800                   | 350          | 280          | 5           | 62                      | 68       | 55         | 58     | 50                     |
|       | 軍神   | 関羽 雲長   | 160    | 38         | &color(red){''*3500''} | 420          | 150          | 4           | &color(yellow){''*95''} | 88       | 72         | 35     | 55                     |
|       | 猛将   | 張飛 翼徳   | 167    | 52         | &color(red){''*3200''} | 380          | 80           | 3           | &color(yellow){''*97''} | 82       | 40         | 18     | &color(green){''*62''} |
|       | 軍師   | 諸葛亮 孔明 | 181    | 30         | 1800                   | 280          | 580          | 3           | 25                      | 38       | 65         | 99     | 45                     |
|~平均  |~       |~            |~       |~41.25      |~2825                   |~357.5        |~272.5        |~3.75        |~69.75                   |~69       |~58         |~52.5   |~53                     |
````

`Stream` output:

````text
|~TEAM  |~CLASS  |~NAME  |~BIRTH  |~ATB GAUGE  |~HIT POINT  |~SKILL POINT  |~SPELL POINT  |~LIFE POINT  |~STRENGTH  |~STAMINA  |~DEXTERITY  |~MAGIC  |~SPEED  |
| 蜀 | 君主 | 劉備 玄徳 | 161 | 45 | 2800 | 350 | 280 | 5 | 62 | 68 | 55 | 58 | 50 |
|  | 軍神 | 関羽 雲長 | 160 | 38 | &color(red){''*3500''} | 420 | 150 | 4 | &color(yellow){''*95''} | 88 | 72 | 35 | 55 |
|  | 猛将 | 張飛 翼徳 | 167 | 52 | &color(red){''*3200''} | 380 | 80 | 3 | &color(yellow){''*97''} | 82 | 40 | 18 | &color(green){''*62''} |
|  | 軍師 | 諸葛亮 孔明 | 181 | 30 | 1800 | 280 | 580 | 3 | 25 | 38 | 65 | 99 | 45 |
|~平均  |~  |~  |~  |~41.25  |~2825  |~357.5  |~272.5  |~3.75  |~69.75  |~69  |~58  |~52.5  |~53  |
````

#### csv

Configuration:

````go
// CSVOptionTransformer configures the transformed delimiter-separated table example.
var CSVOptionTransformer = []csv.Option{
	csv.WithHeader(FooterData.Header[0]),
	csv.WithFooter(FooterData.Footer),
	csv.WithTransformer(csv.Columns(5), func(v any) string {
		n, ok := v.(int)
		if !ok {
			return ""
		}
		if n >= 3000 {
			return fmt.Sprintf("*%d", n)
		}
		return ""
	}),
	csv.WithTransformer(csv.Columns(9), func(v any) string {
		if n, ok := v.(int); ok && n >= 90 {
			return fmt.Sprintf("*%d", n)
		}
		return ""
	}),
	csv.WithTransformer(csv.Columns(13), func(v any) string {
		if n, ok := v.(int); ok && n >= 60 {
			return fmt.Sprintf("*%d", n)
		}
		return ""
	}),
}
````

`Table` and `Stream` output:

````text
TEAM	CLASS	NAME	BIRTH	ATB GAUGE	HIT POINT	SKILL POINT	SPELL POINT	LIFE POINT	STRENGTH	STAMINA	DEXTERITY	MAGIC	SPEED
蜀	君主	劉備 玄徳	161	45	2800	350	280	5	62	68	55	58	50
蜀	軍神	関羽 雲長	160	38	*3500	420	150	4	*95	88	72	35	55
蜀	猛将	張飛 翼徳	167	52	*3200	380	80	3	*97	82	40	18	*62
蜀	軍師	諸葛亮 孔明	181	30	1800	280	580	3	25	38	65	99	45
平均	平均	平均	平均	41.25	2825	357.5	272.5	3.75	69.75	69	58	52.5	53
````

### complex

Input:

````go
// ComplexData exercises arbitrary Go values, from scalar slices to
// nested composites.
var ComplexData = Data{
	Header: [][]string{
		{
			"STRING",
			"NUMBER",
			"FLOAT",
			"TIME.TIME - STRING()",
			"TIME.DURATION - STRING()",
			"STRING SLICE",
			"STRING ARRAY",
			"INT SLICE",
			"STRUCT",
			"MAP",
			"NESTED SLICE",
			"WRAPPED CONTENT",
		},
	},
	Body: [][]any{
		{
			"entry 1",
			123,
			3.14,
			time.Date(2026, 5, 1, 12, 34, 56, 0, time.UTC),
			time.Hour*2 + time.Minute*30,
			[]string{"a", "b", "c"},
			[3]string{"x", "y", "z"},
			[]int{1, 2, 3},
			struct {
				Field1 string
				Field2 int
			}{
				Field1: "value1",
				Field2: 256,
			},
			map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			[]struct {
				Field1 string
				Field2 string
				Field3 string
			}{
				{"Line1", "Line2", "Line3"},
			},
			"Line1\nLine2\nLine3",
		},
		{
			"entry 2",
			000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000,
			0.00000000000000000000000000000000000000001,
			time.Date(2026, 5, 1, 12, 34, 56, 0, time.Local),
			time.Nanosecond,
			[]string{"a", "", "c"},
			[3]string{" ", "y", ""},
			[]int{1, 2, 3, 4, 5},
			struct {
				Field1 string
				Field2 int
				Field3 time.Time
				Field4 time.Duration
			}{
				Field1: "value1",
				Field2: 256,
				Field3: time.Date(2026, 5, 1, 12, 34, 56, 0, time.UTC),
				Field4: time.Hour*2 + time.Minute*30,
			},
			map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
				"key4": "value4",
				"key5": "value5",
			},
			[][]string{
				{"Line1", "Line2", "Line3"},
				{"Line4", "Line5", "Line6"},
			},
			func() []byte {
				b, _ := json.MarshalIndent(SimpleData.Body[:2], "", "  ")
				return b
			}(),
		},
	},
}
````

#### text

Configuration:

````go
// TextOptionComplex configures the complex text table example.
var TextOptionComplex = []text.Option{
	text.WithStyle(text.StyleColoredDouble),
	text.WithHeader(ComplexData.Header...),
	text.WithAutoFit(),
	text.WithAttr(text.ScopeBody, text.Columns(8, 9, 10), text.ColorFgBlack),
	text.WithTransformer(text.Columns(5), func(v any) (string, *text.Attr) {
		values, ok := v.([]string)
		if !ok {
			return "", nil
		}
		tokens := make([]string, len(values))
		for i, s := range values {
			tokens[i] = fmt.Sprintf("%d/%d: %s", i+1, len(values), s)
		}
		return strings.Join(tokens, "\n"), textBgGreenUnderline
	}),
	text.WithTransformer(text.Columns(6), func(v any) (string, *text.Attr) {
		values, ok := v.([3]string)
		if !ok {
			return "", nil
		}
		tokens := make([]string, len(values))
		for i, s := range values {
			tokens[i] = fmt.Sprintf("%d/%d: %s", i+1, len(values), s)
		}
		return strings.Join(tokens, "\n"), textBgMagentaItalic
	}),
	text.WithTransformer(text.Columns(7), func(v any) (string, *text.Attr) {
		values, ok := v.([]int)
		if !ok {
			return "", nil
		}
		sum := 0
		for _, value := range values {
			sum += value
		}
		return fmt.Sprintf("sum=%d", sum), textFgCyanBold
	}),
	text.WithCaption("⚡️ Rendered by <github.com/nekrassov01/table/text>", text.CaptionDefault),
}
````

`Table` output:

````text
╔═════════╤════════╤═══════╤═══════════════════════════════╤══════════════════════════╤══════════════╤══════════════╤═══════════╤════════════════════════════════════════════════════╤══════════════════════════════════════════════════════════════════╤═══════════════════════════════════════════╤════════════════════════════╗
║ STRING  │ NUMBER │ FLOAT │     TIME.TIME - STRING()      │ TIME.DURATION - STRING() │ STRING SLICE │ STRING ARRAY │ INT SLICE │                       STRUCT                       │                               MAP                                │               NESTED SLICE                │      WRAPPED CONTENT       ║
╠═════════╪════════╪═══════╪═══════════════════════════════╪══════════════════════════╪══════════════╪══════════════╪═══════════╪════════════════════════════════════════════════════╪══════════════════════════════════════════════════════════════════╪═══════════════════════════════════════════╪════════════════════════════╣
║ entry 1 │ 123    │ 3.14  │ 2026-05-01 12:34:56 +0000 UTC │ 2h30m0s                  │ 1/3: a       │ 1/3: x       │ sum=6     │ {value1 256}                                       │ map[key1:value1 key2:value2]                                     │ [{Line1 Line2 Line3}]                     │ Line1                      ║
║         │        │       │                               │                          │ 2/3: b       │ 2/3: y       │           │                                                    │                                                                  │                                           │ Line2                      ║
║         │        │       │                               │                          │ 3/3: c       │ 3/3: z       │           │                                                    │                                                                  │                                           │ Line3                      ║
╟─────────┼────────┼───────┼───────────────────────────────┼──────────────────────────┼──────────────┼──────────────┼───────────┼────────────────────────────────────────────────────┼──────────────────────────────────────────────────────────────────┼───────────────────────────────────────────┼────────────────────────────╢
║ entry 2 │ 0      │ 1e-41 │ 2026-05-01 12:34:56 +0000 UTC │ 1ns                      │ 1/3: a       │ 1/3:         │ sum=15    │ {value1 256 2026-05-01 12:34:56 +0000 UTC 2h30m0s} │ map[key1:value1 key2:value2 key3:value3 key4:value4 key5:value5] │ [[Line1 Line2 Line3] [Line4 Line5 Line6]] │ [                          ║
║         │        │       │                               │                          │ 2/3:         │ 2/3: y       │           │                                                    │                                                                  │                                           │   [                        ║
║         │        │       │                               │                          │ 3/3: c       │ 3/3:         │           │                                                    │                                                                  │                                           │     "i-00000000000000000", ║
║         │        │       │                               │                          │              │              │           │                                                    │                                                                  │                                           │     "server-1",            ║
║         │        │       │                               │                          │              │              │           │                                                    │                                                                  │                                           │     "running"              ║
║         │        │       │                               │                          │              │              │           │                                                    │                                                                  │                                           │   ],                       ║
║         │        │       │                               │                          │              │              │           │                                                    │                                                                  │                                           │   [                        ║
║         │        │       │                               │                          │              │              │           │                                                    │                                                                  │                                           │     "i-00000000000000001", ║
║         │        │       │                               │                          │              │              │           │                                                    │                                                                  │                                           │     "server-2",            ║
║         │        │       │                               │                          │              │              │           │                                                    │                                                                  │                                           │     "stopped"              ║
║         │        │       │                               │                          │              │              │           │                                                    │                                                                  │                                           │   ]                        ║
║         │        │       │                               │                          │              │              │           │                                                    │                                                                  │                                           │ ]                          ║
╚═════════╧════════╧═══════╧═══════════════════════════════╧══════════════════════════╧══════════════╧══════════════╧═══════════╧════════════════════════════════════════════════════╧══════════════════════════════════════════════════════════════════╧═══════════════════════════════════════════╧════════════════════════════╝
⚡️ Rendered by <github.com/nekrassov01/table/text>
````

`Stream` output:

````text
╔═════════╤════════╤═══════╤═══════════════════════════════╤══════════════════════════╤══════════════╤══════════════╤═══════════╤══════════════╤══════════════════════════════╤═══════════════════════╤═════════════════╗
║ STRING  │ NUMBER │ FLOAT │     TIME.TIME - STRING()      │ TIME.DURATION - STRING() │ STRING SLICE │ STRING ARRAY │ INT SLICE │    STRUCT    │             MAP              │     NESTED SLICE      │ WRAPPED CONTENT ║
╠═════════╪════════╪═══════╪═══════════════════════════════╪══════════════════════════╪══════════════╪══════════════╪═══════════╪══════════════╪══════════════════════════════╪═══════════════════════╪═════════════════╣
║ entry 1 │ 123    │ 3.14  │ 2026-05-01 12:34:56 +0000 UTC │ 2h30m0s                  │ 1/3: a       │ 1/3: x       │ sum=6     │ {value1 256} │ map[key1:value1 key2:value2] │ [{Line1 Line2 Line3}] │ Line1           ║
║         │        │       │                               │                          │ 2/3: b       │ 2/3: y       │           │              │                              │                       │ Line2           ║
║         │        │       │                               │                          │ 3/3: c       │ 3/3: z       │           │              │                              │                       │ Line3           ║
╟─────────┼────────┼───────┼───────────────────────────────┼──────────────────────────┼──────────────┼──────────────┼───────────┼──────────────┼──────────────────────────────┼───────────────────────┼─────────────────╢
║ entry 2 │ 0      │ 1e-41 │ 2026-05-01 12:34:56 +0000 UTC │ 1ns                      │ 1/3: a       │ 1/3:         │ sum=15    │ {value1 256  │ map[key1:value1 key2:value2  │ [[Line1 Line2 Line3]  │ [               ║
║         │        │       │                               │                          │ 2/3:         │ 2/3: y       │           │ 2026-05-01 1 │ key3:value3 key4:value4 key5 │ [Line4 Line5 Line6]]  │   [             ║
║         │        │       │                               │                          │ 3/3: c       │ 3/3:         │           │ 2:34:56 +000 │ :value5]                     │                       │     "i-00000000 ║
║         │        │       │                               │                          │              │              │           │ 0 UTC 2h30m0 │                              │                       │ 000000000",     ║
║         │        │       │                               │                          │              │              │           │ s}           │                              │                       │     "server-1", ║
║         │        │       │                               │                          │              │              │           │              │                              │                       │     "running"   ║
║         │        │       │                               │                          │              │              │           │              │                              │                       │   ],            ║
║         │        │       │                               │                          │              │              │           │              │                              │                       │   [             ║
║         │        │       │                               │                          │              │              │           │              │                              │                       │     "i-00000000 ║
║         │        │       │                               │                          │              │              │           │              │                              │                       │ 000000001",     ║
║         │        │       │                               │                          │              │              │           │              │                              │                       │     "server-2", ║
║         │        │       │                               │                          │              │              │           │              │                              │                       │     "stopped"   ║
║         │        │       │                               │                          │              │              │           │              │                              │                       │   ]             ║
║         │        │       │                               │                          │              │              │           │              │                              │                       │ ]               ║
╚═════════╧════════╧═══════╧═══════════════════════════════╧══════════════════════════╧══════════════╧══════════════╧═══════════╧══════════════╧══════════════════════════════╧═══════════════════════╧═════════════════╝
⚡️ Rendered by <github.com/nekrassov01/table/text>
````

#### html

Configuration:

````go
// HTMLOptionComplex configures the complex HTML table example.
var HTMLOptionComplex = []html.Option{
	html.WithHeader(ComplexData.Header...),
	html.WithColor(html.ScopeBody, html.Columns(8, 9, 10), html.ColorFgBlack),
	html.WithDecoration(html.ScopeBody, html.Columns(11), html.DecorationPreformatted),
	html.WithTransformer(html.Columns(5), func(v any) (string, *html.Color, *html.Decoration) {
		values, ok := v.([]string)
		if !ok {
			return "", nil, nil
		}
		tokens := make([]string, len(values))
		for i, s := range values {
			tokens[i] = fmt.Sprintf("%d/%d: %s", i+1, len(values), s)
		}
		return strings.Join(tokens, "\n"), html.ColorBgGreen, html.DecorationUnderline
	}),
	html.WithTransformer(html.Columns(6), func(v any) (string, *html.Color, *html.Decoration) {
		values, ok := v.([3]string)
		if !ok {
			return "", nil, nil
		}
		tokens := make([]string, len(values))
		for i, s := range values {
			tokens[i] = fmt.Sprintf("%d/%d: %s", i+1, len(values), s)
		}
		return strings.Join(tokens, "\n"), html.ColorBgMagenta, html.DecorationItalic
	}),
	html.WithTransformer(html.Columns(7), func(v any) (string, *html.Color, *html.Decoration) {
		values, ok := v.([]int)
		if !ok {
			return "", nil, nil
		}
		sum := 0
		for _, value := range values {
			sum += value
		}
		return fmt.Sprintf("sum=%d", sum), html.ColorFgCyan, html.DecorationBold
	}),
	html.WithCaption("⚡️ Rendered by <github.com/nekrassov01/table/html>", html.CaptionDefault),
}
````

`Table` and `Stream` output:

````html
<table>
  <caption>⚡️ Rendered by &lt;github.com/nekrassov01/table/html&gt;</caption>
  <thead>
    <tr>
      <th>STRING</th>
      <th>NUMBER</th>
      <th>FLOAT</th>
      <th>TIME.TIME - STRING()</th>
      <th>TIME.DURATION - STRING()</th>
      <th>STRING SLICE</th>
      <th>STRING ARRAY</th>
      <th>INT SLICE</th>
      <th>STRUCT</th>
      <th>MAP</th>
      <th>NESTED SLICE</th>
      <th>WRAPPED CONTENT</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>entry 1</td>
      <td>123</td>
      <td>3.14</td>
      <td>2026-05-01 12:34:56 +0000 UTC</td>
      <td>2h30m0s</td>
      <td><u><span style="background-color:green">1/3: a<br>2/3: b<br>3/3: c</span></u></td>
      <td><em><span style="background-color:magenta">1/3: x<br>2/3: y<br>3/3: z</span></em></td>
      <td><strong><span style="color:cyan">sum=6</span></strong></td>
      <td><span style="color:black">{value1 256}</span></td>
      <td><span style="color:black">map[key1:value1 key2:value2]</span></td>
      <td><span style="color:black">[{Line1 Line2 Line3}]</span></td>
      <td><pre>Line1<br>Line2<br>Line3</pre></td>
    </tr>
    <tr>
      <td>entry 2</td>
      <td>0</td>
      <td>1e-41</td>
      <td>2026-05-01 12:34:56 +0000 UTC</td>
      <td>1ns</td>
      <td><u><span style="background-color:green">1/3: a<br>2/3: <br>3/3: c</span></u></td>
      <td><em><span style="background-color:magenta">1/3:  <br>2/3: y<br>3/3: </span></em></td>
      <td><strong><span style="color:cyan">sum=15</span></strong></td>
      <td><span style="color:black">{value1 256 2026-05-01 12:34:56 +0000 UTC 2h30m0s}</span></td>
      <td><span style="color:black">map[key1:value1 key2:value2 key3:value3 key4:value4 key5:value5]</span></td>
      <td><span style="color:black">[[Line1 Line2 Line3] [Line4 Line5 Line6]]</span></td>
      <td><pre>[<br>  [<br>    &quot;i-00000000000000000&quot;,<br>    &quot;server-1&quot;,<br>    &quot;running&quot;<br>  ],<br>  [<br>    &quot;i-00000000000000001&quot;,<br>    &quot;server-2&quot;,<br>    &quot;stopped&quot;<br>  ]<br>]</pre></td>
    </tr>
  </tbody>
</table>
````

#### markdown

Configuration:

````go
// MarkdownOptionComplex configures the complex Markdown table example.
var MarkdownOptionComplex = []markdown.Option{
	markdown.WithHeader(ComplexData.Header[0]),
	markdown.WithColor(markdown.ScopeBody, markdown.Columns(8, 9, 10), markdown.ColorFgBlack),
	markdown.WithDecoration(markdown.ScopeBody, markdown.Columns(11), markdown.DecorationUnderline),
	markdown.WithTransformer(markdown.Columns(5), func(v any) (string, *markdown.Color, *markdown.Decoration) {
		values, ok := v.([]string)
		if !ok {
			return "", nil, nil
		}
		tokens := make([]string, len(values))
		for i, s := range values {
			tokens[i] = fmt.Sprintf("%d/%d: %s", i+1, len(values), s)
		}
		return strings.Join(tokens, "\n"), markdown.ColorBgGreen, markdown.DecorationBold
	}),
	markdown.WithTransformer(markdown.Columns(6), func(v any) (string, *markdown.Color, *markdown.Decoration) {
		values, ok := v.([3]string)
		if !ok {
			return "", nil, nil
		}
		tokens := make([]string, len(values))
		for i, s := range values {
			tokens[i] = fmt.Sprintf("%d/%d: %s", i+1, len(values), s)
		}
		return strings.Join(tokens, "\n"), markdown.ColorBgMagenta, markdown.DecorationItalic
	}),
	markdown.WithTransformer(markdown.Columns(7), func(v any) (string, *markdown.Color, *markdown.Decoration) {
		values, ok := v.([]int)
		if !ok {
			return "", nil, nil
		}
		sum := 0
		for _, value := range values {
			sum += value
		}
		return fmt.Sprintf("sum=%d", sum), markdown.ColorFgCyan, markdown.DecorationBold
	}),
}
````

`Table` output:

````markdown
| STRING  | NUMBER | FLOAT | TIME.TIME - STRING()          | TIME.DURATION - STRING() | STRING SLICE                                                               | STRING ARRAY                                                               | INT SLICE                                  | STRUCT                                                                              | MAP                                                                                                 | NESTED SLICE                                                                     | WRAPPED CONTENT                                                                                                                                                                      |
| ------- | ------ | ----- | ----------------------------- | ------------------------ | -------------------------------------------------------------------------- | -------------------------------------------------------------------------- | ------------------------------------------ | ----------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| entry 1 | 123    | 3.14  | 2026-05-01 12:34:56 +0000 UTC | 2h30m0s                  | **<span style="background-color:green">1/3: a<br>2/3: b<br>3/3: c</span>** | *<span style="background-color:magenta">1/3: x<br>2/3: y<br>3/3: z</span>* | **<span style="color:cyan">sum=6</span>**  | <span style="color:black">{value1 256}</span>                                       | <span style="color:black">map\[key1:value1 key2:value2\]</span>                                     | <span style="color:black">\[{Line1 Line2 Line3}\]</span>                         | <u>Line1<br>Line2<br>Line3</u>                                                                                                                                                       |
| entry 2 | 0      | 1e-41 | 2026-05-01 12:34:56 +0000 UTC | 1ns                      | **<span style="background-color:green">1/3: a<br>2/3: <br>3/3: c</span>**  | *<span style="background-color:magenta">1/3:  <br>2/3: y<br>3/3: </span>*  | **<span style="color:cyan">sum=15</span>** | <span style="color:black">{value1 256 2026-05-01 12:34:56 +0000 UTC 2h30m0s}</span> | <span style="color:black">map\[key1:value1 key2:value2 key3:value3 key4:value4 key5:value5\]</span> | <span style="color:black">\[\[Line1 Line2 Line3\] \[Line4 Line5 Line6\]\]</span> | <u>\[<br>  \[<br>    "i-00000000000000000",<br>    "server-1",<br>    "running"<br>  \],<br>  \[<br>    "i-00000000000000001",<br>    "server-2",<br>    "stopped"<br>  \]<br>\]</u> |
````

`Stream` output:

````markdown
| STRING | NUMBER | FLOAT | TIME.TIME - STRING() | TIME.DURATION - STRING() | STRING SLICE | STRING ARRAY | INT SLICE | STRUCT | MAP | NESTED SLICE | WRAPPED CONTENT |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| entry 1 | 123 | 3.14 | 2026-05-01 12:34:56 +0000 UTC | 2h30m0s | **<span style="background-color:green">1/3: a<br>2/3: b<br>3/3: c</span>** | *<span style="background-color:magenta">1/3: x<br>2/3: y<br>3/3: z</span>* | **<span style="color:cyan">sum=6</span>** | <span style="color:black">{value1 256}</span> | <span style="color:black">map\[key1:value1 key2:value2\]</span> | <span style="color:black">\[{Line1 Line2 Line3}\]</span> | <u>Line1<br>Line2<br>Line3</u> |
| entry 2 | 0 | 1e-41 | 2026-05-01 12:34:56 +0000 UTC | 1ns | **<span style="background-color:green">1/3: a<br>2/3: <br>3/3: c</span>** | *<span style="background-color:magenta">1/3:  <br>2/3: y<br>3/3: </span>* | **<span style="color:cyan">sum=15</span>** | <span style="color:black">{value1 256 2026-05-01 12:34:56 +0000 UTC 2h30m0s}</span> | <span style="color:black">map\[key1:value1 key2:value2 key3:value3 key4:value4 key5:value5\]</span> | <span style="color:black">\[\[Line1 Line2 Line3\] \[Line4 Line5 Line6\]\]</span> | <u>\[<br>  \[<br>    "i-00000000000000000",<br>    "server-1",<br>    "running"<br>  \],<br>  \[<br>    "i-00000000000000001",<br>    "server-2",<br>    "stopped"<br>  \]<br>\]</u> |
````

#### backlog

Configuration:

````go
// BacklogOptionComplex configures the complex Backlog table example.
var BacklogOptionComplex = []backlog.Option{
	backlog.WithHeader(ComplexData.Header...),
	backlog.WithColor(backlog.ScopeBody, backlog.Columns(8, 9, 10), backlog.ColorFgBlack),
	backlog.WithDecoration(backlog.ScopeBody, backlog.Columns(11), backlog.DecorationBold),
	backlog.WithTransformer(backlog.Columns(5), func(v any) (string, *backlog.Color, *backlog.Decoration) {
		values, ok := v.([]string)
		if !ok {
			return "", nil, nil
		}
		tokens := make([]string, len(values))
		for i, s := range values {
			tokens[i] = fmt.Sprintf("%d/%d: %s", i+1, len(values), s)
		}
		return strings.Join(tokens, "\n"), backlog.ColorBgGreen, backlog.DecorationBold
	}),
	backlog.WithTransformer(backlog.Columns(6), func(v any) (string, *backlog.Color, *backlog.Decoration) {
		values, ok := v.([3]string)
		if !ok {
			return "", nil, nil
		}
		tokens := make([]string, len(values))
		for i, s := range values {
			tokens[i] = fmt.Sprintf("%d/%d: %s", i+1, len(values), s)
		}
		return strings.Join(tokens, "\n"), backlog.ColorBgYellow, backlog.DecorationItalic
	}),
	backlog.WithTransformer(backlog.Columns(7), func(v any) (string, *backlog.Color, *backlog.Decoration) {
		values, ok := v.([]int)
		if !ok {
			return "", nil, nil
		}
		sum := 0
		for _, value := range values {
			sum += value
		}
		return fmt.Sprintf("sum=%d", sum), backlog.ColorFgBlue, backlog.DecorationStrikethrough
	}),
}
````

`Table` output:

````text
|~STRING  |~NUMBER  |~FLOAT  |~TIME.TIME - STRING()          |~TIME.DURATION - STRING()  |~STRING SLICE                                      |~STRING ARRAY                                         |~INT SLICE                |~STRUCT                                                            |~MAP                                                                             |~NESTED SLICE                                                     |~WRAPPED CONTENT                                                                                                                                                             |
| entry 1 | 123     | 3.14   | 2026-05-01 12:34:56 +0000 UTC | 2h30m0s                   | &color("", green){''1/3: a&br;2/3: b&br;3/3: c''} | &color("", yellow){'''1/3: x&br;2/3: y&br;3/3: z'''} | &color(blue){%%sum=6%%}  | &color(black){{value1 256}}                                       | &color(black){map[key1:value1 key2:value2]}                                     | &color(black){[{Line1 Line2 Line3}]}                             | ''Line1&br;Line2&br;Line3''                                                                                                                                                 |
| entry 2 | 0       | 1e-41  | 2026-05-01 12:34:56 +0000 UTC | 1ns                       | &color("", green){''1/3: a&br;2/3: &br;3/3: c''}  | &color("", yellow){'''1/3:  &br;2/3: y&br;3/3: '''}  | &color(blue){%%sum=15%%} | &color(black){{value1 256 2026-05-01 12:34:56 +0000 UTC 2h30m0s}} | &color(black){map[key1:value1 key2:value2 key3:value3 key4:value4 key5:value5]} | &color(black){\\[\\[Line1 Line2 Line3] [Line4 Line5 Line6\\]\\]} | ''[&br;  [&br;    "i-00000000000000000",&br;    "server-1",&br;    "running"&br;  ],&br;  [&br;    "i-00000000000000001",&br;    "server-2",&br;    "stopped"&br;  ]&br;]'' |
````

`Stream` output:

````text
|~STRING  |~NUMBER  |~FLOAT  |~TIME.TIME - STRING()  |~TIME.DURATION - STRING()  |~STRING SLICE  |~STRING ARRAY  |~INT SLICE  |~STRUCT  |~MAP  |~NESTED SLICE  |~WRAPPED CONTENT  |
| entry 1 | 123 | 3.14 | 2026-05-01 12:34:56 +0000 UTC | 2h30m0s | &color("", green){''1/3: a&br;2/3: b&br;3/3: c''} | &color("", yellow){'''1/3: x&br;2/3: y&br;3/3: z'''} | &color(blue){%%sum=6%%} | &color(black){{value1 256}} | &color(black){map[key1:value1 key2:value2]} | &color(black){[{Line1 Line2 Line3}]} | ''Line1&br;Line2&br;Line3'' |
| entry 2 | 0 | 1e-41 | 2026-05-01 12:34:56 +0000 UTC | 1ns | &color("", green){''1/3: a&br;2/3: &br;3/3: c''} | &color("", yellow){'''1/3:  &br;2/3: y&br;3/3: '''} | &color(blue){%%sum=15%%} | &color(black){{value1 256 2026-05-01 12:34:56 +0000 UTC 2h30m0s}} | &color(black){map[key1:value1 key2:value2 key3:value3 key4:value4 key5:value5]} | &color(black){\\[\\[Line1 Line2 Line3] [Line4 Line5 Line6\\]\\]} | ''[&br;  [&br;    "i-00000000000000000",&br;    "server-1",&br;    "running"&br;  ],&br;  [&br;    "i-00000000000000001",&br;    "server-2",&br;    "stopped"&br;  ]&br;]'' |
````

#### csv

Configuration:

````go
// CSVOptionComplex configures the complex delimiter-separated table example.
var CSVOptionComplex = []csv.Option{
	csv.WithHeader(ComplexData.Header[0]),
	csv.WithTransformer(csv.Columns(5), func(v any) string {
		values, ok := v.([]string)
		if !ok {
			return ""
		}
		tokens := make([]string, len(values))
		for i, s := range values {
			tokens[i] = fmt.Sprintf("%d/%d: %s", i+1, len(values), s)
		}
		return strings.Join(tokens, "\n")
	}),
	csv.WithTransformer(csv.Columns(6), func(v any) string {
		values, ok := v.([3]string)
		if !ok {
			return ""
		}
		tokens := make([]string, len(values))
		for i, s := range values {
			tokens[i] = fmt.Sprintf("%d/%d: %s", i+1, len(values), s)
		}
		return strings.Join(tokens, "\n")
	}),
	csv.WithTransformer(csv.Columns(7), func(v any) string {
		values, ok := v.([]int)
		if !ok {
			return ""
		}
		sum := 0
		for _, value := range values {
			sum += value
		}
		return fmt.Sprintf("sum=%d", sum)
	}),
}
````

`Table` and `Stream` output:

````text
STRING	NUMBER	FLOAT	TIME.TIME - STRING()	TIME.DURATION - STRING()	STRING SLICE	STRING ARRAY	INT SLICE	STRUCT	MAP	NESTED SLICE	WRAPPED CONTENT
entry 1	123	3.14	2026-05-01 12:34:56 +0000 UTC	2h30m0s	"1/3: a
2/3: b
3/3: c"	"1/3: x
2/3: y
3/3: z"	sum=6	{value1 256}	map[key1:value1 key2:value2]	[{Line1 Line2 Line3}]	"Line1
Line2
Line3"
entry 2	0	1e-41	2026-05-01 12:34:56 +0000 UTC	1ns	"1/3: a
2/3: 
3/3: c"	"1/3:  
2/3: y
3/3: "	sum=15	{value1 256 2026-05-01 12:34:56 +0000 UTC 2h30m0s}	map[key1:value1 key2:value2 key3:value3 key4:value4 key5:value5]	[[Line1 Line2 Line3] [Line4 Line5 Line6]]	"[
  [
    ""i-00000000000000000"",
    ""server-1"",
    ""running""
  ],
  [
    ""i-00000000000000001"",
    ""server-2"",
    ""stopped""
  ]
]"
````

### stacked-header

Input:

````go
// StackedHeaderData provides a multi-row header for formats that support one.
// Repeated adjacent labels provide the values used for column spanning, so the
// rows need no separate grouping metadata.
var StackedHeaderData = Data{
	Header: [][]string{
		{
			"AWS RESOURCE",
			"AWS RESOURCE",
			"AWS RESOURCE",
			"AWS RESOURCE",
			"ID",
		},
		{
			"NETWORK",
			"NETWORK",
			"SECURITY",
			"SECURITY",
			"ID",
		},
		{
			"VPC",
			"SUBNET",
			"SG",
			"NACL",
			"ID",
		},
	},
	Body: [][]any{
		{
			"vpc-1",
			"sub-1",
			"sg-1",
			"nacl-1",
			"i-001",
		},
		{
			"vpc-2",
			"sub-2",
			"sg-2",
			"nacl-2",
			"i-002",
		},
	},
}
````

#### text

Configuration:

````go
// TextOptionStackedHeader configures the text table example with a stacked header.
var TextOptionStackedHeader = []text.Option{
	text.WithStyle(text.StyleColoredLight),
	text.WithHeader(StackedHeaderData.Header...),
	text.WithRowspan(text.ScopeHeader, text.Columns(4)),
	text.WithColspan(text.ScopeHeader, text.Columns(0, 1, 2, 3)),
}
````

`Table` and `Stream` output:

````text
┌────────────────────────────────┬───────┐
│          AWS RESOURCE          │       │
├────────────────┬───────────────┤       │
│    NETWORK     │   SECURITY    │       │
├───────┬────────┼──────┬────────┤       │
│  VPC  │ SUBNET │  SG  │  NACL  │  ID   │
╞═══════╪════════╪══════╪════════╪═══════╡
│ vpc-1 │ sub-1  │ sg-1 │ nacl-1 │ i-001 │
├───────┼────────┼──────┼────────┼───────┤
│ vpc-2 │ sub-2  │ sg-2 │ nacl-2 │ i-002 │
└───────┴────────┴──────┴────────┴───────┘
````

#### html

Configuration:

````go
// HTMLOptionStackedHeader configures the HTML table example with a stacked header.
var HTMLOptionStackedHeader = []html.Option{
	html.WithHeader(StackedHeaderData.Header...),
	html.WithRowspan(html.ScopeHeader, html.Columns(4)),
	html.WithColspan(html.ScopeHeader, html.Columns(0, 1, 2, 3)),
}
````

`Table` and `Stream` output:

````html
<table>
  <thead>
    <tr>
      <th colspan="4">AWS RESOURCE</th>
      <th rowspan="3">ID</th>
    </tr>
    <tr>
      <th colspan="2">NETWORK</th>
      <th colspan="2">SECURITY</th>
    </tr>
    <tr>
      <th>VPC</th>
      <th>SUBNET</th>
      <th>SG</th>
      <th>NACL</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>vpc-1</td>
      <td>sub-1</td>
      <td>sg-1</td>
      <td>nacl-1</td>
      <td>i-001</td>
    </tr>
    <tr>
      <td>vpc-2</td>
      <td>sub-2</td>
      <td>sg-2</td>
      <td>nacl-2</td>
      <td>i-002</td>
    </tr>
  </tbody>
</table>
````

#### backlog

Configuration:

````go
// BacklogOptionStackedHeader configures the Backlog table example with a stacked header.
var BacklogOptionStackedHeader = []backlog.Option{
	backlog.WithHeader(StackedHeaderData.Header...),
	backlog.WithRowspan(backlog.ScopeHeader, backlog.Columns(4)),
	backlog.WithColspan(backlog.ScopeHeader, backlog.Columns(0, 1, 2, 3)),
}
````

`Table` output:

````text
|~AWS RESOURCE  |~        |~          |~       |~      |
|~NETWORK       |~        |~SECURITY  |~       |~      |
|~VPC           |~SUBNET  |~SG        |~NACL   |~ID    |
| vpc-1         | sub-1   | sg-1      | nacl-1 | i-001 |
| vpc-2         | sub-2   | sg-2      | nacl-2 | i-002 |
````

`Stream` output:

````text
|~AWS RESOURCE  |~  |~  |~  |~  |
|~NETWORK  |~  |~SECURITY  |~  |~  |
|~VPC  |~SUBNET  |~SG  |~NACL  |~ID  |
| vpc-1 | sub-1 | sg-1 | nacl-1 | i-001 |
| vpc-2 | sub-2 | sg-2 | nacl-2 | i-002 |
````

### comma-included

Input:

````go
// CommaIncludedData contains comma-bearing values that exercise CSV
// quoting.
var CommaIncludedData = Data{
	Header: [][]string{
		{
			"PRODUCT ID",
			"PRODUCT NAME",
			"PRICE",
			"PROPERTIES",
		},
	},
	Body: [][]any{
		{
			"p-00000000000000000",
			"product-1",
			"1,000.00",
			[]string{"color:red", "size:large", "weight:1.5kg"},
		},
		{
			"p-00000000000000001",
			"product-2",
			"2,500.00",
			[]string{"color:blue", "size:medium", "weight:2.0kg"},
		},
		{
			"p-00000000000000002",
			"product-3",
			"3,750.00",
			[]string{"color:green", "size:small", "weight:1.0kg"},
		},
	},
}
````

#### csv

Configuration:

````go
// CSVOptionCommaIncluded configures the comma-separated table example.
var CSVOptionCommaIncluded = []csv.Option{
	csv.WithHeader(CommaIncludedData.Header[0]),
	csv.WithDelimiter(','),
}
````

`Table` and `Stream` output:

````text
PRODUCT ID,PRODUCT NAME,PRICE,PROPERTIES
p-00000000000000000,product-1,"1,000.00",[color:red size:large weight:1.5kg]
p-00000000000000001,product-2,"2,500.00",[color:blue size:medium weight:2.0kg]
p-00000000000000002,product-3,"3,750.00",[color:green size:small weight:1.0kg]
````
