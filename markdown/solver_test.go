package markdown

import (
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func Test_solver_prepare(t *testing.T) {
	type fields struct {
		input compilerResult
		state solverState
	}
	type want struct {
		metrics    []columnMetric
		metricsCap int
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "initializes alignment without measuring rows",
			fields: fields{
				input: compilerResult{
					configResult: configResult{
						option: &option{
							indexOffset: 1,
						},
						columns: []columnConfig{
							{},
							{
								align: AlignCenter,
							},
						},
					},
					header: row{
						cells: []cell{
							{
								value: "#",
								width: 1,
								size:  1,
							},
							{
								value: "A",
								width: 1,
								size:  1,
							},
						},
					},
					body: []row{
						{
							cells: []cell{
								{
									value: "1",
									width: 1,
									size:  1,
								},
								{
									value: "日本",
									width: 4,
									size:  6,
								},
							},
						},
					},
				},
				state: solverState{
					columnMetrics: make([]columnMetric, 1, 4),
				},
			},
			want: want{
				metrics: []columnMetric{
					{
						box: box{
							align: AlignRight,
						},
						separator: separator{
							width: 4,
							trail: true,
						},
					},
					{
						box: box{
							align: AlignCenter,
						},
						separator: separator{
							width: 5,
							lead:  true,
							trail: true,
						},
					},
				},
				metricsCap: 4,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			o := &solver{
				input: test.fields.input,
				state: &state,
				output: solverResult{
					compilerResult: test.fields.input,
				},
			}
			o.prepare()
			got := want{
				metrics:    o.output.metrics,
				metricsCap: cap(state.columnMetrics),
			}
			testutil.AssertValue(t, got, test.want, "prepare")
		})
	}
}

func Test_solver_measureRows(t *testing.T) {
	type fields struct {
		state solverState
	}
	type args struct {
		rows []row
	}
	type want struct {
		metrics []columnMetric
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "takes maxima across rows",
			fields: fields{
				state: solverState{
					columnMetrics: []columnMetric{{}},
				},
			},
			args: args{
				rows: []row{
					{
						cells: []cell{
							{
								value: "a",
								width: 1,
							},
						},
					},
					{
						cells: []cell{
							{
								value: "wide",
								width: 4,
							},
						},
					},
				},
			},
			want: want{
				metrics: []columnMetric{
					{
						box: box{width: 4},
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			o := &solver{
				state: &state,
			}
			o.measureRows(test.args.rows)
			got := want{
				metrics: state.columnMetrics,
			}
			testutil.AssertValue(t, got, test.want, "measureRows")
		})
	}
}

func Test_solver_measureRow(t *testing.T) {
	type fields struct {
		state solverState
	}
	type args struct {
		row row
	}
	type want struct {
		metrics []columnMetric
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "measures maximum display width",
			fields: fields{
				state: solverState{
					columnMetrics: []columnMetric{{}},
				},
			},
			args: args{
				row: row{
					cells: []cell{
						{
							value: "日",
							width: 37,
						},
					},
				},
			},
			want: want{
				metrics: []columnMetric{
					{
						box: box{width: 37},
					},
				},
			},
		},
		{
			name: "keeps existing maxima",
			fields: fields{
				state: solverState{
					columnMetrics: []columnMetric{
						{
							box: box{width: 5},
						},
					},
				},
			},
			args: args{
				row: row{
					cells: []cell{
						{
							value: "x",
							width: 1,
							size:  1,
						},
					},
				},
			},
			want: want{
				metrics: []columnMetric{
					{
						box: box{width: 5},
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			o := &solver{
				state: &state,
			}
			r := test.args.row
			o.measureRow(&r)
			got := want{
				metrics: state.columnMetrics,
			}
			testutil.AssertValue(t, got, test.want, "measureRow")
		})
	}
}

func Test_solver_solve(t *testing.T) {
	type fields struct {
		input compilerResult
		state solverState
	}
	type want struct {
		metrics []columnMetric
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "raises widths to content and delimiter minima",
			fields: fields{
				input: compilerResult{
					header: row{
						cells: []cell{
							{width: 2},
							{width: 6},
							{width: 1},
						},
					},
					body: []row{
						{
							cells: []cell{
								{width: 1},
								{width: 7},
								{width: 1},
							},
						},
					},
				},
				state: solverState{
					columnMetrics: []columnMetric{
						{
							separator: separator{
								width: 3,
							},
						},
						{
							box: box{
								align: AlignCenter,
							},
							separator: separator{
								width: 5,
								lead:  true,
								trail: true,
							},
						},
						{
							box: box{
								align: AlignRight,
							},
							separator: separator{
								width: 4,
								trail: true,
							},
						},
					},
				},
			},
			want: want{
				metrics: []columnMetric{
					{
						box: box{
							width: 3,
						},
						separator: separator{
							width: 3,
						},
					},
					{
						box: box{
							width: 7,
							align: AlignCenter,
						},
						separator: separator{
							width: 5,
							lead:  true,
							trail: true,
						},
					},
					{
						box: box{
							width: 4,
							align: AlignRight,
						},
						separator: separator{
							width: 4,
							trail: true,
						},
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			o := &solver{
				input: test.fields.input,
				state: &state,
			}
			o.solve()
			got := want{
				metrics: state.columnMetrics,
			}
			testutil.AssertValue(t, got, test.want, "solve")
		})
	}
}

func TestResolveSeparator(t *testing.T) {
	type args struct {
		side AlignSide
	}
	type want struct {
		separator separator
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "default",
			want: want{
				separator: separator{
					width: 3,
				},
			},
		},
		{
			name: "left",
			args: args{
				side: AlignLeft,
			},
			want: want{
				separator: separator{
					width: 4,
					lead:  true,
				},
			},
		},
		{
			name: "right",
			args: args{
				side: AlignRight,
			},
			want: want{
				separator: separator{
					width: 4,
					trail: true,
				},
			},
		},
		{
			name: "center",
			args: args{
				side: AlignCenter,
			},
			want: want{
				separator: separator{
					width: 5,
					lead:  true,
					trail: true,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := want{
				separator: resolveSeparator(test.args.side),
			}
			testutil.AssertValue(t, got, test.want, "resolveSeparator")
		})
	}
}
