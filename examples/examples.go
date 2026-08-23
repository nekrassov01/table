package examples

import (
	"encoding/json"
	"strconv"
	"sync/atomic"
	"time"
)

// Data is one example input, shared by the examples and the benchmarks.
// A stacked header spans by repeating a label across its rows, and the two
// formats that allow one header row take the first. Footer is a function so
// examples can calculate its rows for each Table.Render or Stream.Close call.
type Data struct {
	Header [][]string        // The header rows, topmost first.
	Body   [][]any           // The data rows.
	Footer func() [][]string // Returns the footer rows, topmost first; nil for none.
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
	body := make([][]any, 0, len(d.Body)*100)
	for range 100 {
		body = append(body, d.Body...)
	}
	return Data{
		Header: d.Header,
		Body:   body,
		Footer: d.Footer,
	}
}
