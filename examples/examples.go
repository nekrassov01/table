package examples

import (
	"encoding/json"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/nekrassov01/table"
)

// Data is one example input, shared by the examples and the benchmarks.
// A stacked header spans by repeating a label across its rows, and the two
// formats that allow one header row take the first. Footer is a function so
// examples can calculate its rows for each Table.Render or Stream.Close call.
type Data struct {
	// The header rows, topmost first.
	Header [][]string
	// The data rows.
	Body [][]table.Value
	// Returns the footer rows, topmost first; nil for none.
	Footer func() [][]string
}

// SimpleData is a small instance list of plain strings.
var SimpleData = Data{
	Header: [][]string{
		{
			"INSTANCE ID",
			"INSTANCE NAME",
			"INSTANCE STATE",
		},
	},
	Body: [][]table.Value{
		{
			table.String("i-00000000000000000"),
			table.String("server-1"),
			table.String("running"),
		},
		{
			table.String("i-00000000000000001"),
			table.String("server-2"),
			table.String("stopped"),
		},
		{
			table.String("i-00000000000000002"),
			table.String("server-3"),
			table.String("pending"),
		},
		{
			table.String("i-00000000000000003"),
			table.String("server-4"),
			table.String("terminated"),
		},
		{
			table.String("i-00000000000000004"),
			table.String("server-5"),
			table.String("stopping"),
		},
		{
			table.String("i-00000000000000005"),
			table.String("server-6"),
			table.String("shutting-down"),
		},
	},
}

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
	Body: [][]table.Value{
		{
			table.String("accesslog"),
			table.String(time.Date(2026, 5, 1, 9, 1, 15, 0, time.UTC).String()),
			table.String("healthcheck ok"),
			table.Int(200),
		},
		{
			table.String("accesslog"),
			table.String(time.Date(2026, 5, 1, 9, 1, 16, 0, time.UTC).String()),
			table.String("authentication ok"),
			table.Int(200),
		},
		{
			table.String("accesslog"),
			table.String(time.Date(2026, 5, 1, 9, 1, 19, 0, time.UTC).String()),
			table.String("get resource ok"),
			table.Int(200),
		},
		{
			table.String("application"),
			table.String(time.Date(2026, 5, 1, 0, 19, 21, 0, time.Local).String()),
			table.String("GET /api/v1/users/ HTTP/1.1"),
			table.Int(200),
		},
		{
			table.String("application"),
			table.String(time.Date(2026, 5, 1, 0, 20, 57, 0, time.Local).String()),
			table.String("GET /api/v1/users/alice/ HTTP/1.1"),
			table.Int(200),
		},
		{
			table.String("application"),
			table.String(time.Date(2026, 5, 1, 1, 05, 34, 0, time.Local).String()),
			table.String("POST /api/v1/users/bob/ HTTP/1.1"),
			table.Int(201),
		},
		{
			table.String("application"),
			table.String(time.Date(2026, 5, 1, 1, 07, 56, 0, time.Local).String()),
			table.String("DELETE /api/v1/users/alice/ HTTP/1.1"),
			table.Int(204),
		},
	},
}

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
	Body: [][]table.Value{
		{
			table.String("i-00000000000000000"),
			table.String("sg-10000000000000000"),
			table.String("Ingress"),
			table.String("tcp"),
			table.Int(22),
			table.Int(22),
			table.String("SecurityGroup"),
			table.String("sg-20000000000000000"),
		},
		{
			table.String("i-00000000000000000"),
			table.String("sg-10000000000000000"),
			table.String("Egress"),
			table.String("-1"),
			table.Int(0),
			table.Int(0),
			table.String("Ipv4"),
			table.String("0.0.0.0/0"),
		},
		{
			table.String("i-00000000000000000"),
			table.String("sg-10000000000000001"),
			table.String("Ingress"),
			table.String("tcp"),
			table.Int(443),
			table.Int(443),
			table.String("Ipv4"),
			table.String("0.0.0.0/0"),
		},
		{
			table.String("i-00000000000000000"),
			table.String("sg-10000000000000001"),
			table.String("Egress"),
			table.String("-1"),
			table.Int(0),
			table.Int(0),
			table.String("Ipv4"),
			table.String("0.0.0.0/0"),
		},
		{
			table.String("i-00000000000000001"),
			table.String("sg-10000000000000002"),
			table.String("Ingress"),
			table.String("tcp"),
			table.Int(3389),
			table.Int(3389),
			table.String("Ipv4"),
			table.String("10.1.0.0/16"),
		},
		{
			table.String("i-00000000000000001"),
			table.String("sg-10000000000000002"),
			table.String("Ingress"),
			table.String("tcp"),
			table.Int(0),
			table.Int(65535),
			table.String("PrefixList"),
			table.String("pl-00000000/com.amazonaws.ap-northeast-1.s3"),
		},
		{
			table.String("i-00000000000000001"),
			table.String("sg-10000000000000002"),
			table.String("Egress"),
			table.String("-1"),
			table.Int(0),
			table.Int(0),
			table.String("Ipv4"),
			table.String("0.0.0.0/0"),
		},
		{
			table.String("i-00000000000000002"),
			table.String("sg-10000000000000003"),
			table.String("Ingress"),
			table.String("tcp"),
			table.Int(443),
			table.Int(443),
			table.String("Ipv4"),
			table.String("0.0.0.0/0"),
		},
		{
			table.String("i-00000000000000002"),
			table.String("sg-10000000000000003"),
			table.String("Egress"),
			table.String("-1"),
			table.Int(0),
			table.Int(0),
			table.String("Ipv4"),
			table.String("0.0.0.0/0"),
		},
	},
}

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
	Body: [][]table.Value{
		{
			table.String("John"),
			table.String("Doe"),
			table.Int(30),
			table.String(time.Date(1994, 5, 1, 0, 0, 0, 0, time.UTC).String()),
		},
		{
			table.String("Jane"),
			table.String("Smith"),
			table.Int(25),
			table.String(time.Date(1999, 5, 1, 0, 0, 0, 0, time.UTC).String()),
		},
		{
			table.String("Anonymous"),
			table.String("Anonymous"),
			table.String("Unknown"),
			table.String("Unknown"),
		},
		{
			table.String("Alice"),
			table.String("Johnson"),
			table.Int(28),
			table.String(time.Date(1996, 5, 1, 0, 0, 0, 0, time.UTC).String()),
		},
	},
}

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
		Body: [][]table.Value{
			{
				table.String("蜀"),
				table.String("君主"),
				table.String("劉備 玄徳"),
				table.Int(161),
				table.Int(45),
				table.Int(2800),
				table.Int(350),
				table.Int(280),
				table.Int(5),
				table.Int(62),
				table.Int(68),
				table.Int(55),
				table.Int(58),
				table.Int(50),
			},
			{
				table.String("蜀"),
				table.String("軍神"),
				table.String("関羽 雲長"),
				table.Int(160),
				table.Int(38),
				table.Int(3500),
				table.Int(420),
				table.Int(150),
				table.Int(4),
				table.Int(95),
				table.Int(88),
				table.Int(72),
				table.Int(35),
				table.Int(55),
			},
			{
				table.String("蜀"),
				table.String("猛将"),
				table.String("張飛 翼徳"),
				table.Int(167),
				table.Int(52),
				table.Int(3200),
				table.Int(380),
				table.Int(80),
				table.Int(3),
				table.Int(97),
				table.Int(82),
				table.Int(40),
				table.Int(18),
				table.Int(62),
			},
			{
				table.String("蜀"),
				table.String("軍師"),
				table.String("諸葛亮 孔明"),
				table.Int(181),
				table.Int(30),
				table.Int(1800),
				table.Int(280),
				table.Int(580),
				table.Int(3),
				table.Int(25),
				table.Int(38),
				table.Int(65),
				table.Int(99),
				table.Int(45),
			},
		},
	}
	var cached atomic.Pointer[footerState]
	data.Footer = func() [][]string {
		var totals [valueColumns]int
		for column := range valueColumns {
			for _, row := range data.Body {
				totals[column] += row[labelColumns+column].AsInt()
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
	Body: [][]table.Value{
		{
			table.String("entry 1"),
			table.Int(123),
			table.Float64(3.14),
			table.Any(time.Date(2026, 5, 1, 12, 34, 56, 0, time.UTC)),
			table.Any(time.Hour*2 + time.Minute*30),
			table.Any([]string{"a", "b", "c"}),
			table.Any([3]string{"x", "y", "z"}),
			table.Any([]int{1, 2, 3}),
			table.Any(struct {
				Field1 string
				Field2 int
			}{
				Field1: "value1",
				Field2: 256,
			}),
			table.Any(map[string]string{
				"key1": "value1",
				"key2": "value2",
			}),
			table.Any([]struct {
				Field1 string
				Field2 string
				Field3 string
			}{
				{"Line1", "Line2", "Line3"},
			}),
			table.String("Line1\nLine2\nLine3"),
		},
		{
			table.String("entry 2"),
			table.Int(000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000),
			table.Float64(0.00000000000000000000000000000000000000001),
			table.Any(time.Date(2026, 5, 1, 12, 34, 56, 0, time.Local)),
			table.Any(time.Nanosecond),
			table.Any([]string{"a", "", "c"}),
			table.Any([3]string{" ", "y", ""}),
			table.Any([]int{1, 2, 3, 4, 5}),
			table.Any(struct {
				Field1 string
				Field2 int
				Field3 time.Time
				Field4 time.Duration
			}{
				Field1: "value1",
				Field2: 256,
				Field3: time.Date(2026, 5, 1, 12, 34, 56, 0, time.UTC),
				Field4: time.Hour*2 + time.Minute*30,
			}),
			table.Any(map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
				"key4": "value4",
				"key5": "value5",
			}),
			table.Any([][]string{
				{"Line1", "Line2", "Line3"},
				{"Line4", "Line5", "Line6"},
			}),
			table.Any(func() []byte {
				rows := make([][]any, 2)
				for i, row := range SimpleData.Body[:2] {
					rows[i] = make([]any, len(row))
					for j, cell := range row {
						rows[i][j] = cell.AsAny()
					}
				}
				b, _ := json.MarshalIndent(rows, "", "  ")
				return b
			}()),
		},
	},
}

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
	Body: [][]table.Value{
		{
			table.String("vpc-1"),
			table.String("sub-1"),
			table.String("sg-1"),
			table.String("nacl-1"),
			table.String("i-001"),
		},
		{
			table.String("vpc-2"),
			table.String("sub-2"),
			table.String("sg-2"),
			table.String("nacl-2"),
			table.String("i-002"),
		},
	},
}

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
	Body: [][]table.Value{
		{
			table.String("p-00000000000000000"),
			table.String("product-1"),
			table.String("1,000.00"),
			table.Any([]string{"color:red", "size:large", "weight:1.5kg"}),
		},
		{
			table.String("p-00000000000000001"),
			table.String("product-2"),
			table.String("2,500.00"),
			table.Any([]string{"color:blue", "size:medium", "weight:2.0kg"}),
		},
		{
			table.String("p-00000000000000002"),
			table.String("product-3"),
			table.String("3,750.00"),
			table.Any([]string{"color:green", "size:small", "weight:1.0kg"}),
		},
	},
}

var (
	// SimpleDataLarge is SimpleData with its rows repeated 100 times.
	SimpleDataLarge = largeOf(SimpleData)

	// CompactDataLarge is CompactData with its rows repeated 100 times.
	CompactDataLarge = largeOf(CompactData)

	// RowspanDataLarge is RowspanData with its rows repeated 100 times.
	RowspanDataLarge = largeOf(RowspanData)

	// ColspanDataLarge is ColspanData with its rows repeated 100 times.
	ColspanDataLarge = largeOf(ColspanData)

	// FooterDataLarge is FooterData with its rows repeated 100 times.
	FooterDataLarge = largeOf(FooterData)

	// ComplexDataLarge is ComplexData with its rows repeated 100 times.
	ComplexDataLarge = largeOf(ComplexData)

	// StackedHeaderDataLarge is StackedHeaderData with its rows repeated 100 times.
	StackedHeaderDataLarge = largeOf(StackedHeaderData)

	// CommaIncludedDataLarge is CommaIncludedData with its rows repeated 100 times.
	CommaIncludedDataLarge = largeOf(CommaIncludedData)
)

// largeOf returns d with its data rows repeated 100 times, for comparing
// streaming against batch rendering on a larger input.
func largeOf(d Data) Data {
	body := make([][]table.Value, 0, len(d.Body)*100)
	for range 100 {
		body = append(body, d.Body...)
	}
	return Data{
		Header: d.Header,
		Body:   body,
		Footer: d.Footer,
	}
}
