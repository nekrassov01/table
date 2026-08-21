package text

import (
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func TestHorizontal_maxGlyphLen(t *testing.T) {
	type fields struct {
		Inner Joints
		Outer Joints
		Fill  string
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
			name: "inner joint",
			fields: &fields{
				Inner: Joints{
					UDLR: "inner",
				},
			},
			want: want{
				val: 5,
			},
		},
		{
			name: "outer joint",
			fields: &fields{
				Outer: Joints{
					XXXX: "outside",
				},
			},
			want: want{
				val: 7,
			},
		},
		{
			name: "fill",
			fields: &fields{
				Fill: "horizontal",
			},
			want: want{
				val: 10,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var o *Horizontal
			if test.fields != nil {
				o = &Horizontal{
					Inner: test.fields.Inner,
					Outer: test.fields.Outer,
					Fill:  test.fields.Fill,
				}
			}
			got := o.maxGlyphLen()
			testutil.AssertValue(t, got, test.want.val, "maxGlyphLen")
		})
	}
}

func TestVertical_maxGlyphLen(t *testing.T) {
	type fields struct {
		Outer string
		Inner string
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
			name: "outer",
			fields: &fields{
				Outer: "outside",
				Inner: "in",
			},
			want: want{
				val: 7,
			},
		},
		{
			name: "inner",
			fields: &fields{
				Outer: "out",
				Inner: "inside",
			},
			want: want{
				val: 6,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var o *Vertical
			if test.fields != nil {
				o = &Vertical{
					Outer: test.fields.Outer,
					Inner: test.fields.Inner,
				}
			}
			got := o.maxGlyphLen()
			testutil.AssertValue(t, got, test.want.val, "maxGlyphLen")
		})
	}
}

func TestJoints_maxGlyphLen(t *testing.T) {
	type fields struct {
		UDLR string
		UDLX string
		UDXR string
		UDXX string
		XDLR string
		XDLX string
		XDXR string
		XDXX string
		UXLR string
		UXLX string
		UXXR string
		UXXX string
		XXLR string
		XXLX string
		XXXR string
		XXXX string
	}
	type want struct {
		val int
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "empty",
		},
		{
			name: "up and down row",
			fields: fields{
				UDLX: "1234",
			},
			want: want{
				val: 4,
			},
		},
		{
			name: "down row",
			fields: fields{
				XDXR: "12345",
			},
			want: want{
				val: 5,
			},
		},
		{
			name: "up row",
			fields: fields{
				UXXX: "123456",
			},
			want: want{
				val: 6,
			},
		},
		{
			name: "horizontal row",
			fields: fields{
				XXXX: "1234567",
			},
			want: want{
				val: 7,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &Joints{
				UDLR: test.fields.UDLR,
				UDLX: test.fields.UDLX,
				UDXR: test.fields.UDXR,
				UDXX: test.fields.UDXX,
				XDLR: test.fields.XDLR,
				XDLX: test.fields.XDLX,
				XDXR: test.fields.XDXR,
				XDXX: test.fields.XDXX,
				UXLR: test.fields.UXLR,
				UXLX: test.fields.UXLX,
				UXXR: test.fields.UXXR,
				UXXX: test.fields.UXXX,
				XXLR: test.fields.XXLR,
				XXLX: test.fields.XXLX,
				XXXR: test.fields.XXXR,
				XXXX: test.fields.XXXX,
			}
			got := o.maxGlyphLen()
			testutil.AssertValue(t, got, test.want.val, "maxGlyphLen")
		})
	}
}

func TestJoints_resolve(t *testing.T) {
	type fields struct {
		UDLR string
		UDLX string
		UDXR string
		UDXX string
		XDLR string
		XDLX string
		XDXR string
		XDXX string
		UXLR string
		UXLX string
		UXXR string
		UXXX string
		XXLR string
		XXLX string
		XXXR string
		XXXX string
	}
	type args struct {
		up    bool
		down  bool
		left  bool
		right bool
	}
	type want struct {
		val string
	}
	jointFields := fields{
		UDLR: "UDLR",
		UDLX: "UDLX",
		UDXR: "UDXR",
		UDXX: "UDXX",
		XDLR: "XDLR",
		XDLX: "XDLX",
		XDXR: "XDXR",
		XDXX: "XDXX",
		UXLR: "UXLR",
		UXLX: "UXLX",
		UXXR: "UXXR",
		UXXX: "UXXX",
		XXLR: "XXLR",
		XXLX: "XXLX",
		XXXR: "XXXR",
		XXXX: "XXXX",
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name:   "up down left right",
			fields: jointFields,
			args: args{
				up:    true,
				down:  true,
				left:  true,
				right: true,
			},
			want: want{
				val: "UDLR",
			},
		},
		{
			name:   "up down left",
			fields: jointFields,
			args: args{
				up:   true,
				down: true,
				left: true,
			},
			want: want{
				val: "UDLX",
			},
		},
		{
			name:   "up down right",
			fields: jointFields,
			args: args{
				up:    true,
				down:  true,
				right: true,
			},
			want: want{
				val: "UDXR",
			},
		},
		{
			name:   "up down",
			fields: jointFields,
			args: args{
				up:   true,
				down: true,
			},
			want: want{
				val: "UDXX",
			},
		},
		{
			name:   "down left right",
			fields: jointFields,
			args: args{
				down:  true,
				left:  true,
				right: true,
			},
			want: want{
				val: "XDLR",
			},
		},
		{
			name:   "down left",
			fields: jointFields,
			args: args{
				down: true,
				left: true,
			},
			want: want{
				val: "XDLX",
			},
		},
		{
			name:   "down right",
			fields: jointFields,
			args: args{
				down:  true,
				right: true,
			},
			want: want{
				val: "XDXR",
			},
		},
		{
			name:   "down",
			fields: jointFields,
			args: args{
				down: true,
			},
			want: want{
				val: "XDXX",
			},
		},
		{
			name:   "up left right",
			fields: jointFields,
			args: args{
				up:    true,
				left:  true,
				right: true,
			},
			want: want{
				val: "UXLR",
			},
		},
		{
			name:   "up left",
			fields: jointFields,
			args: args{
				up:   true,
				left: true,
			},
			want: want{
				val: "UXLX",
			},
		},
		{
			name:   "up right",
			fields: jointFields,
			args: args{
				up:    true,
				right: true,
			},
			want: want{
				val: "UXXR",
			},
		},
		{
			name:   "up",
			fields: jointFields,
			args: args{
				up: true,
			},
			want: want{
				val: "UXXX",
			},
		},
		{
			name:   "left right",
			fields: jointFields,
			args: args{
				left:  true,
				right: true,
			},
			want: want{
				val: "XXLR",
			},
		},
		{
			name:   "left",
			fields: jointFields,
			args: args{
				left: true,
			},
			want: want{
				val: "XXLX",
			},
		},
		{
			name:   "right",
			fields: jointFields,
			args: args{
				right: true,
			},
			want: want{
				val: "XXXR",
			},
		},
		{
			name:   "no arms",
			fields: jointFields,
			want: want{
				val: "XXXX",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &Joints{
				UDLR: test.fields.UDLR,
				UDLX: test.fields.UDLX,
				UDXR: test.fields.UDXR,
				UDXX: test.fields.UDXX,
				XDLR: test.fields.XDLR,
				XDLX: test.fields.XDLX,
				XDXR: test.fields.XDXR,
				XDXX: test.fields.XDXX,
				UXLR: test.fields.UXLR,
				UXLX: test.fields.UXLX,
				UXXR: test.fields.UXXR,
				UXXX: test.fields.UXXX,
				XXLR: test.fields.XXLR,
				XXLX: test.fields.XXLX,
				XXXR: test.fields.XXXR,
				XXXX: test.fields.XXXX,
			}
			got := o.resolve(
				test.args.up,
				test.args.down,
				test.args.left,
				test.args.right,
			)
			testutil.AssertValue(t, got, test.want.val, "resolve")
		})
	}
}
