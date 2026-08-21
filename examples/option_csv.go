package examples

import (
	"fmt"
	"strings"

	"github.com/nekrassov01/table/csv"
)

// CSVOptionSimple configures the simple delimiter-separated table example.
var CSVOptionSimple = []csv.Option{
	csv.WithHeader(SimpleData.Header[0]),
}

// CSVOptionFooter configures the delimiter-separated table example with a footer.
var CSVOptionFooter = []csv.Option{
	csv.WithHeader(FooterData.Header[0]),
	csv.WithFooter(FooterData.Footer),
}

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

// CSVOptionCommaIncluded configures the comma-separated table example.
var CSVOptionCommaIncluded = []csv.Option{
	csv.WithHeader(CommaIncludedData.Header[0]),
	csv.WithDelimiter(','),
}
