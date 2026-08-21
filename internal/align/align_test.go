package align

import (
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func TestSide_String(t *testing.T) {
	type fields struct {
		value Side
	}
	type want struct {
		val string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "default",
			fields: fields{
				value: Default,
			},
			want: want{
				val: "default",
			},
		},
		{
			name: "left",
			fields: fields{
				value: Left,
			},
			want: want{
				val: "left",
			},
		},
		{
			name: "right",
			fields: fields{
				value: Right,
			},
			want: want{
				val: "right",
			},
		},
		{
			name: "center",
			fields: fields{
				value: Center,
			},
			want: want{
				val: "center",
			},
		},
		{
			name: "unknown",
			fields: fields{
				value: Side(255),
			},
			want: want{
				val: "default",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testutil.AssertValue(t, test.fields.value.String(), test.want.val, "String")
		})
	}
}
