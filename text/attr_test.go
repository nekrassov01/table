package text

import (
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func TestNewAttr(t *testing.T) {
	type args struct {
		codes []Code
	}
	type want struct {
		val *Attr
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "no codes",
		},
		{
			name: "one code",
			args: args{
				codes: []Code{CodeFgRed},
			},
			want: want{
				val: &Attr{
					Prefix: []byte("\x1b[31m"),
					Suffix: []byte("\x1b[0m"),
				},
			},
		},
		{
			name: "multiple codes",
			args: args{
				codes: []Code{CodeFgGreen, CodeBgWhite, CodeBold},
			},
			want: want{
				val: &Attr{
					Prefix: []byte("\x1b[32;47;1m"),
					Suffix: []byte("\x1b[0m"),
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := NewAttr(test.args.codes...)
			testutil.AssertValue(t, got, test.want.val, "NewAttr")
		})
	}
}

func TestAttr_isZero(t *testing.T) {
	type fields struct {
		Prefix []byte
		Suffix []byte
	}
	type want struct {
		val bool
	}
	tests := []struct {
		name   string
		fields *fields
		want   want
	}{
		{
			name: "nil",
			want: want{
				val: true,
			},
		},
		{
			name:   "empty",
			fields: &fields{},
			want: want{
				val: true,
			},
		},
		{
			name: "suffix only",
			fields: &fields{
				Suffix: []byte("suffix"),
			},
			want: want{
				val: true,
			},
		},
		{
			name: "prefix",
			fields: &fields{
				Prefix: []byte("prefix"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var o *Attr
			if test.fields != nil {
				o = &Attr{
					Prefix: test.fields.Prefix,
					Suffix: test.fields.Suffix,
				}
			}
			got := o.isZero()
			testutil.AssertValue(t, got, test.want.val, "isZero")
		})
	}
}

func TestAttr_size(t *testing.T) {
	type fields struct {
		Prefix []byte
		Suffix []byte
	}
	type want struct {
		val int
	}
	tests := []struct {
		name   string
		fields *fields
		want   want
	}{
		{
			name: "nil",
		},
		{
			name:   "empty",
			fields: &fields{},
		},
		{
			name: "suffix only",
			fields: &fields{
				Suffix: []byte("suffix"),
			},
		},
		{
			name: "prefix only",
			fields: &fields{
				Prefix: []byte("prefix"),
			},
			want: want{
				val: 6,
			},
		},
		{
			name: "prefix and suffix",
			fields: &fields{
				Prefix: []byte("prefix"),
				Suffix: []byte("suffix"),
			},
			want: want{
				val: 12,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var o *Attr
			if test.fields != nil {
				o = &Attr{
					Prefix: test.fields.Prefix,
					Suffix: test.fields.Suffix,
				}
			}
			got := o.size()
			testutil.AssertValue(t, got, test.want.val, "size")
		})
	}
}

func Test_buildSGR(t *testing.T) {
	type args struct {
		codes []Code
	}
	type want struct {
		val []byte
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "no codes",
		},
		{
			name: "one code",
			args: args{
				codes: []Code{CodeUnderline},
			},
			want: want{
				val: []byte("\x1b[4m"),
			},
		},
		{
			name: "multiple codes",
			args: args{
				codes: []Code{CodeFgHiMagenta, CodeBgHiCyan, CodeReverseVideo},
			},
			want: want{
				val: []byte("\x1b[95;106;7m"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := buildSGR(test.args.codes)
			testutil.AssertValue(t, got, test.want.val, "buildSGR")
		})
	}
}
