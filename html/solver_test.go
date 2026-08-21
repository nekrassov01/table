package html

import (
	"testing"

	"github.com/nekrassov01/table/internal/param"
	"github.com/nekrassov01/table/internal/testutil"
)

func Test_solver_solve(t *testing.T) {
	type fields struct {
		input compilerResult
		state solverState
	}
	type want struct {
		headerRowspans []int
		headerColspans []int
		bodyRowspans   []int
		bodyColspans   []int
		footerRowspans []int
		footerColspans []int
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "resolves each table part",
			fields: fields{
				state: func() solverState {
					state := solverState{}
					state.rowspanCounts[0] = 99
					return state
				}(),
				input: compilerResult{
					header: []row{
						{
							cells: []cell{{colspan: 1}, {colspan: 1}},
						},
						{
							cells:    []cell{{colspan: 1}, {colspan: 1}},
							rowspans: 0b01,
						},
					},
					body: []row{
						{
							cells: []cell{
								{colspan: 1},
								{colspan: 1},
							},
							colspans: 0b10,
						},
					},
					footer: []row{
						{
							cells: []cell{{colspan: 1}, {colspan: 1}},
						},
						{
							cells:    []cell{{colspan: 1}, {colspan: 1}},
							rowspans: 0b01,
						},
					},
				},
			},
			want: want{
				headerRowspans: []int{2, 0, 0, 0},
				headerColspans: []int{1, 1, 0, 1},
				bodyRowspans:   []int{0, 0},
				bodyColspans:   []int{2, 0},
				footerRowspans: []int{2, 0, 0, 0},
				footerColspans: []int{1, 1, 0, 1},
			},
		},
		{
			name: "leaves rows without span facts",
			fields: fields{
				input: compilerResult{
					body: []row{
						{
							cells: []cell{{rowspan: 4, colspan: 3}},
						},
					},
				},
			},
			want: want{
				bodyRowspans: []int{4},
				bodyColspans: []int{3},
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
			o.solve()
			got := want{}
			for _, r := range o.output.header {
				for _, compiled := range r.cells {
					got.headerRowspans = append(got.headerRowspans, compiled.rowspan)
					got.headerColspans = append(got.headerColspans, compiled.colspan)
				}
			}
			for _, r := range o.output.body {
				for _, compiled := range r.cells {
					got.bodyRowspans = append(got.bodyRowspans, compiled.rowspan)
					got.bodyColspans = append(got.bodyColspans, compiled.colspan)
				}
			}
			for _, r := range o.output.footer {
				for _, compiled := range r.cells {
					got.footerRowspans = append(got.footerRowspans, compiled.rowspan)
					got.footerColspans = append(got.footerColspans, compiled.colspan)
				}
			}
			testutil.AssertValue(t, got, test.want, "solve")
		})
	}
}

func Test_solver_resolveRows(t *testing.T) {
	type fields struct {
		state solverState
	}
	type args struct {
		rows []row
	}
	type want struct {
		rowspans []int
		colspans []int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "settles rowspans before colspans",
			args: args{
				rows: []row{
					{
						cells: []cell{
							{colspan: 1},
							{colspan: 1},
						},
						colspans: 0b10,
					},
					{
						cells: []cell{
							{colspan: 1},
							{colspan: 1},
						},
						rowspans: 0b01,
					},
				},
			},
			want: want{
				rowspans: []int{2, 0, 0, 0},
				colspans: []int{1, 1, 0, 1},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			o := &solver{
				state: &state,
			}
			o.resolveRows(test.args.rows)
			got := want{}
			for _, r := range test.args.rows {
				for _, compiled := range r.cells {
					got.rowspans = append(got.rowspans, compiled.rowspan)
					got.colspans = append(got.colspans, compiled.colspan)
				}
			}
			testutil.AssertValue(t, got, test.want, "resolveRows")
		})
	}
}

func Test_solver_resolveRowspans(t *testing.T) {
	type fields struct {
		state solverState
	}
	type args struct {
		rows []row
	}
	type want struct {
		rowspans []int
		colspans []int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "counts vertical continuations",
			args: args{
				rows: []row{
					{cells: []cell{{colspan: 1}}},
					{cells: []cell{{colspan: 1}}, rowspans: 0b1},
					{cells: []cell{{colspan: 1}}, rowspans: 0b1},
				},
			},
			want: want{
				rowspans: []int{3, 0, 0},
				colspans: []int{1, 0, 0},
			},
		},
		{
			name: "leaves counts without continuations",
			args: args{
				rows: []row{
					{cells: []cell{{rowspan: 4, colspan: 1}}},
				},
			},
			want: want{
				rowspans: []int{4},
				colspans: []int{1},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			o := &solver{
				state: &state,
			}
			o.resolveRowspans(test.args.rows)
			got := want{}
			for _, r := range test.args.rows {
				for _, compiled := range r.cells {
					got.rowspans = append(got.rowspans, compiled.rowspan)
					got.colspans = append(got.colspans, compiled.colspan)
				}
			}
			testutil.AssertValue(t, got, test.want, "resolveRowspans")
		})
	}
}

func Test_solver_resolveColspans(t *testing.T) {
	type args struct {
		rows []row
	}
	type want struct {
		colspans []int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "absorbs adjacent candidates",
			args: args{
				rows: []row{
					{
						cells: []cell{
							{rowspan: 1, colspan: 1},
							{rowspan: 1, colspan: 1},
							{rowspan: 1, colspan: 1},
						},
						colspans: 0b110,
					},
				},
			},
			want: want{
				colspans: []int{3, 0, 0},
			},
		},
		{
			name: "keeps nonrectangular cells separate",
			args: args{
				rows: []row{
					{
						cells: []cell{
							{rowspan: 2, colspan: 1},
							{rowspan: 1, colspan: 1},
						},
						colspans: 0b10,
					},
				},
			},
			want: want{
				colspans: []int{1, 1},
			},
		},
		{
			name: "keeps cells without candidates",
			args: args{
				rows: []row{
					{
						cells: []cell{
							{rowspan: 1, colspan: 1},
							{rowspan: 1, colspan: 1},
						},
					},
				},
			},
			want: want{
				colspans: []int{1, 1},
			},
		},
		{
			name: "leaves positions beyond the candidate mask",
			args: args{
				rows: func() []row {
					cells := make([]cell, param.SpanLimit+1)
					for index := range cells {
						cells[index] = cell{
							rowspan: 1,
							colspan: 1,
						}
					}
					return []row{
						{
							cells:    cells,
							colspans: 0b10,
						},
					}
				}(),
			},
			want: func() want {
				colspans := make([]int, param.SpanLimit+1)
				for index := range colspans {
					colspans[index] = 1
				}
				colspans[0] = 2
				colspans[1] = 0
				return want{
					colspans: colspans,
				}
			}(),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &solver{}
			o.resolveColspans(test.args.rows)
			got := want{}
			for _, r := range test.args.rows {
				for _, compiled := range r.cells {
					got.colspans = append(got.colspans, compiled.colspan)
				}
			}
			testutil.AssertValue(t, got, test.want, "resolveColspans")
		})
	}
}
