package repeat

import (
	"strings"
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func TestAppendSpaces(t *testing.T) {
	type args struct {
		dst   []byte
		count int
	}
	type want struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "appends spaces",
			args: args{
				dst:   []byte("prefix"),
				count: 3,
			},
			want: want{
				value: "prefix   ",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := want{
				value: string(AppendSpaces(test.args.dst, test.args.count)),
			}
			testutil.AssertValue(t, got, test.want, "AppendSpaces")
		})
	}
}

func TestAppendDashes(t *testing.T) {
	type args struct {
		dst   []byte
		count int
	}
	type want struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "appends dashes",
			args: args{
				dst:   []byte("prefix"),
				count: 257,
			},
			want: want{
				value: "prefix" + strings.Repeat("-", 257),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := want{
				value: string(AppendDashes(test.args.dst, test.args.count)),
			}
			testutil.AssertValue(t, got, test.want, "AppendDashes")
		})
	}
}

func Test_appendBlock(t *testing.T) {
	type args struct {
		dst   []byte
		block []byte
		count int
	}
	type want struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "keeps destination for non-positive count",
			args: args{
				dst:   []byte("prefix"),
				block: []byte("--"),
				count: -1,
			},
			want: want{
				value: "prefix",
			},
		},
		{
			name: "appends multiple blocks",
			args: args{
				dst:   []byte("prefix"),
				block: []byte("---"),
				count: 7,
			},
			want: want{
				value: "prefix-------",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := want{
				value: string(appendBlock(test.args.dst, test.args.block, test.args.count)),
			}
			testutil.AssertValue(t, got, test.want, "appendBlock")
		})
	}
}
