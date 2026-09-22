package backlog

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
		metrics []columnMetric
		output  solverResult
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "initializes metrics without measuring rows",
			fields: fields{
				input: compilerResult{
					configResult: configResult{
						columns: []columnConfig{{}, {}},
					},
					header: []row{{
						cells: []cell{{width: 2}, {width: 5}},
					}},
					body: []row{{
						cells: []cell{{width: 4}, {width: 2}},
					}},
					footer: []row{{
						cells: []cell{{width: 3}, {width: 1}},
					}},
				},
				state: solverState{
					columnMetrics: make([]columnMetric, 1, 2),
				},
			},
			want: want{
				metrics: []columnMetric{{}, {}},
			},
		},
		{
			name: "grows metrics",
			fields: fields{
				input: compilerResult{
					configResult: configResult{
						columns: []columnConfig{{}, {}},
					},
				},
			},
			want: want{
				metrics: []columnMetric{{}, {}},
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
			expected := test.want
			expected.output = solverResult{
				compilerResult: test.fields.input,
				metrics:        test.want.metrics,
			}
			got := want{
				metrics: state.columnMetrics,
				output:  o.output,
			}
			testutil.AssertValue(t, got, expected, "prepare")
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
			name: "measures header body and footer",
			fields: fields{
				input: compilerResult{
					header: []row{{
						cells: []cell{{width: 2}, {width: 5}},
					}},
					body: []row{{
						cells: []cell{{width: 4}, {width: 2}},
					}},
					footer: []row{{
						cells: []cell{{width: 3}, {width: 1}},
					}},
				},
				state: solverState{
					columnMetrics: []columnMetric{{}, {}},
				},
			},
			want: want{
				metrics: []columnMetric{
					{box: box{width: 4}},
					{box: box{width: 5}},
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
			name: "body rows",
			fields: fields{
				state: solverState{
					columnMetrics: []columnMetric{{}, {}},
				},
			},
			args: args{
				rows: []row{
					{cells: []cell{{width: 2}, {width: 5}}},
					{cells: []cell{{width: 4}, {width: 3}}},
				},
			},
			want: want{
				metrics: []columnMetric{{box: box{width: 4}}, {box: box{width: 5}}},
			},
		},
		{
			name: "single row",
			fields: fields{
				state: solverState{
					columnMetrics: []columnMetric{{}},
				},
			},
			args: args{
				rows: []row{{cells: []cell{{width: 2}}}},
			},
			want: want{
				metrics: []columnMetric{{box: box{width: 2}}},
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
			name: "keeps maximum widths",
			fields: fields{
				state: solverState{
					columnMetrics: []columnMetric{{box: box{width: 5}}, {box: box{width: 1}}},
				},
			},
			args: args{
				row: row{
					cells: []cell{{width: 3}, {width: 4}},
				},
			},
			want: want{
				metrics: []columnMetric{{box: box{width: 5}}, {box: box{width: 4}}},
			},
		},
		{
			name: "empty cell",
			fields: fields{
				state: solverState{
					columnMetrics: []columnMetric{{}},
				},
			},
			args: args{
				row: row{
					cells: []cell{{}},
				},
			},
			want: want{
				metrics: []columnMetric{{}},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			o := &solver{
				state: &state,
			}
			o.measureRow(&test.args.row)
			got := want{
				metrics: state.columnMetrics,
			}
			testutil.AssertValue(t, got, test.want, "measureRow")
		})
	}
}
