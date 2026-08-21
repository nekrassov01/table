package text

import (
	"testing"

	"github.com/nekrassov01/table/internal/testutil"
)

func Test_solver_prepare(t *testing.T) {
	type fields struct {
		input      compilerResult
		state      solverState
		widthLimit int
	}
	type want struct {
		metrics          []columnMetric
		spanRequirements []spanRequirement
		outputMetrics    []columnMetric
		ownsMetrics      bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "allocates column metrics",
			fields: fields{
				input: func() compilerResult {
					return compilerResult{
						configResult: configResult{
							option: &option{
								placeholder: "-",
								style: Style{
									Border: BorderStyle{
										Vertical: &Vertical{
											Inner: "|",
										},
									},
								},
							},
							columns: []column{
								{
									lPad: 1,
									rPad: 2,
								},
								{
									limit: 4,
									rPad:  1,
								},
							},
						},
						header: []row{
							{
								cells: []cell{
									{
										value: "界",
									},
									{
										value: "abc",
									},
								},
							},
						},
					}
				}(),
			},
			want: want{
				metrics: []columnMetric{
					{
						box: box{
							lPad: 1,
							rPad: 2,
						},
					},
					{
						box: box{
							rPad: 1,
						},
						limit: 4,
					},
				},
				outputMetrics: []columnMetric{
					{
						box: box{
							lPad: 1,
							rPad: 2,
						},
					},
					{
						box: box{
							rPad: 1,
						},
						limit: 4,
					},
				},
				ownsMetrics: true,
			},
		},
		{
			name: "reuses metric storage and clears requirements",
			fields: fields{
				input: compilerResult{
					configResult: configResult{
						option: &option{},
						columns: []column{
							{
								lPad: 2,
								rPad: 3,
							},
						},
					},
				},
				state: solverState{
					columnMetrics: make([]columnMetric, 2, 3),
					spanRequirements: []spanRequirement{
						{
							start: 0,
							end:   2,
							width: 8,
						},
					},
				},
			},
			want: want{
				metrics: []columnMetric{
					{
						box: box{
							lPad: 2,
							rPad: 3,
						},
					},
				},
				spanRequirements: []spanRequirement{},
				outputMetrics: []columnMetric{
					{
						box: box{
							lPad: 2,
							rPad: 3,
						},
					},
				},
				ownsMetrics: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			o := &solver{
				input:      test.fields.input,
				state:      &state,
				widthLimit: test.fields.widthLimit,
				output: solverResult{
					compilerResult: test.fields.input,
				},
			}
			o.prepare()
			got := want{
				metrics:          state.columnMetrics,
				spanRequirements: state.spanRequirements,
				outputMetrics:    o.output.metrics,
				ownsMetrics:      &o.output.metrics[0] == &state.columnMetrics[0],
			}
			testutil.AssertValue(t, got, test.want, "prepare")
		})
	}
}

func Test_solver_solve(t *testing.T) {
	type fields struct {
		input compilerResult
		state solverState
	}
	type want struct {
		metrics          []columnMetric
		spanRequirements []spanRequirement
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "resolves column geometry",
			fields: fields{
				input: func() compilerResult {
					return compilerResult{
						configResult: configResult{
							option: &option{
								placeholder: "-",
								style: Style{
									Border: BorderStyle{
										Vertical: &Vertical{
											Inner: "|",
										},
									},
								},
							},
						},
						header: []row{
							{
								cells: []cell{
									{value: "界"},
									{value: "abc"},
								},
							},
						},
					}
				}(),
				state: solverState{
					columnMetrics: []columnMetric{
						{
							box: box{
								lPad: 1,
								rPad: 2,
							},
						},
						{
							box: box{
								rPad: 1,
							},
							limit: 4,
						},
					},
				},
			},
			want: want{
				metrics: []columnMetric{
					{
						box: box{
							width: 2,
							lPad:  1,
							rPad:  2,
						},
						overhead: 1,
					},
					{
						box: box{
							offset: 6,
							width:  4,
							rPad:   1,
						},
						limit: 4,
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
				metrics:          state.columnMetrics,
				spanRequirements: state.spanRequirements,
			}
			testutil.AssertValue(t, got, test.want, "solve")
		})
	}
}

func Test_solver_freeze(t *testing.T) {
	type fields struct {
		input      compilerResult
		state      solverState
		widthLimit int
		output     solverResult
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
			name: "fixes only unconstrained widths",
			fields: fields{
				state: solverState{
					columnMetrics: []columnMetric{
						{
							box: box{
								width: 3,
							},
						},
						{
							box: box{
								width: 5,
							},
							limit: 2,
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
						limit: 3,
					},
					{
						box: box{
							width: 5,
						},
						limit: 2,
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			o := &solver{
				input:      test.fields.input,
				state:      &state,
				widthLimit: test.fields.widthLimit,
				output:     test.fields.output,
			}
			o.freeze()
			got := want{
				metrics: state.columnMetrics,
			}
			testutil.AssertValue(t, got, test.want, "freeze")
		})
	}
}

func Test_solver_measureRows(t *testing.T) {
	type fields struct {
		input      compilerResult
		state      solverState
		widthLimit int
		output     solverResult
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
			name: "accumulates every row",
			fields: fields{
				state: solverState{
					columnMetrics: make([]columnMetric, 1),
				},
			},
			args: args{
				rows: []row{
					{
						cells: []cell{
							{
								value: "a",
							},
						},
					},
					{
						cells: []cell{
							{
								value: "界",
							},
						},
					},
				},
			},
			want: want{
				metrics: []columnMetric{
					{
						box: box{
							width: 2,
						},
						overhead: 1,
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			o := &solver{
				input:      test.fields.input,
				state:      &state,
				widthLimit: test.fields.widthLimit,
				output:     test.fields.output,
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
		input      compilerResult
		state      solverState
		widthLimit int
		output     solverResult
	}
	type args struct {
		row row
	}
	type want struct {
		metrics          []columnMetric
		spanRequirements []spanRequirement
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "measures visible cells and skips spans",
			fields: fields{
				state: solverState{
					columnMetrics: make([]columnMetric, 3),
				},
			},
			args: args{
				row: row{
					cells: []cell{
						{
							value: "界",
						},
						{
							value: "rowspan",
						},
						{
							value: "colspan",
						},
					},
					rowspans: 0b010,
					colspans: 0b100,
				},
			},
			want: want{
				metrics: []columnMetric{
					{
						box: box{
							width: 2,
						},
						overhead: 1,
					},
					{},
					{},
				},
			},
		},
		{
			name: "measures the widest physical line",
			fields: fields{
				state: solverState{
					columnMetrics: make([]columnMetric, 1),
				},
			},
			args: args{
				row: row{
					cells: []cell{
						{
							value: "ab\n界界",
						},
					},
				},
			},
			want: want{
				metrics: []columnMetric{
					{
						box: box{
							width: 4,
						},
						overhead: 2,
					},
				},
			},
		},
		{
			name: "updates and orders span requirements",
			fields: fields{
				state: solverState{
					columnMetrics: make([]columnMetric, 6),
					spanRequirements: []spanRequirement{
						{
							start: 0,
							end:   2,
							width: 4,
						},
						{
							start: 3,
							end:   4,
							width: 1,
						},
					},
				},
			},
			args: args{
				row: row{
					cells: []cell{
						{
							value: "aaaaaa",
						},
						{},
						{
							value: "bbb",
						},
						{},
						{
							value: "cc",
						},
						{},
					},
					colspans: 0b101010,
				},
			},
			want: want{
				metrics: []columnMetric{
					{},
					{},
					{},
					{},
					{},
					{},
				},
				spanRequirements: []spanRequirement{
					{
						start: 0,
						end:   2,
						width: 6,
					},
					{
						start: 2,
						end:   4,
						width: 3,
					},
					{
						start: 3,
						end:   4,
						width: 1,
					},
					{
						start: 4,
						end:   6,
						width: 2,
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			o := &solver{
				input:      test.fields.input,
				state:      &state,
				widthLimit: test.fields.widthLimit,
				output:     test.fields.output,
			}
			r := test.args.row
			o.measureRow(&r)
			got := want{
				metrics:          state.columnMetrics,
				spanRequirements: state.spanRequirements,
			}
			testutil.AssertValue(t, got, test.want, "measureRow")
		})
	}
}

func Test_solver_resolveWidths(t *testing.T) {
	type fields struct {
		input      compilerResult
		state      solverState
		widthLimit int
		output     solverResult
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
			name: "resolves index placeholder and fixed widths",
			fields: fields{
				input: compilerResult{
					configResult: configResult{
						option: &option{
							placeholder: "界",
							indexOffset: 1,
							indexWidth:  3,
						},
					},
				},
				state: solverState{
					columnMetrics: []columnMetric{
						{
							box: box{
								width: 1,
							},
						},
						{
							box: box{
								width: 1,
							},
						},
						{
							box: box{
								width: 5,
							},
							limit:    4,
							overhead: 2,
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
					},
					{
						box: box{
							width: 2,
						},
						overhead: 1,
					},
					{
						box: box{
							width: 4,
						},
						limit:    4,
						overhead: 2,
					},
				},
			},
		},
		{
			name: "satisfied span keeps widths",
			fields: fields{
				input: compilerResult{
					configResult: configResult{
						option: &option{
							style: Style{
								Border: BorderStyle{
									Vertical: &Vertical{
										Inner: "|",
									},
								},
							},
						},
					},
				},
				state: solverState{
					columnMetrics: []columnMetric{
						{
							box: box{
								width: 3,
								lPad:  1,
								rPad:  1,
							},
						},
						{
							box: box{
								width: 3,
								lPad:  1,
								rPad:  1,
							},
						},
					},
					spanRequirements: []spanRequirement{
						{
							start: 0,
							end:   2,
							width: 8,
						},
					},
				},
			},
			want: want{
				metrics: []columnMetric{
					{
						box: box{
							width: 3,
							lPad:  1,
							rPad:  1,
						},
					},
					{
						box: box{
							width: 3,
							lPad:  1,
							rPad:  1,
						},
					},
				},
			},
		},
		{
			name: "fixed span cannot absorb deficit",
			fields: fields{
				input: compilerResult{
					configResult: configResult{
						option: &option{},
					},
				},
				state: solverState{
					columnMetrics: []columnMetric{
						{
							box: box{
								width: 1,
							},
							limit: 2,
						},
						{
							box: box{
								width: 1,
							},
							limit: 2,
						},
					},
					spanRequirements: []spanRequirement{
						{
							start: 0,
							end:   2,
							width: 10,
						},
					},
				},
			},
			want: want{
				metrics: []columnMetric{
					{
						box: box{
							width: 2,
						},
						limit: 2,
					},
					{
						box: box{
							width: 2,
						},
						limit: 2,
					},
				},
			},
		},
		{
			name: "distributes span deficit across flexible columns",
			fields: fields{
				input: compilerResult{
					configResult: configResult{
						option: &option{
							style: Style{
								Border: BorderStyle{
									Vertical: &Vertical{
										Inner: "|",
									},
								},
							},
						},
					},
				},
				state: solverState{
					columnMetrics: []columnMetric{
						{
							box: box{
								width: 1,
								lPad:  1,
								rPad:  1,
							},
						},
						{
							box: box{
								width: 2,
								lPad:  1,
								rPad:  1,
							},
							limit: 3,
						},
						{
							box: box{
								width: 1,
								lPad:  1,
								rPad:  1,
							},
						},
					},
					spanRequirements: []spanRequirement{
						{
							start: 0,
							end:   3,
							width: 16,
						},
					},
				},
			},
			want: want{
				metrics: []columnMetric{
					{
						box: box{
							width: 4,
							lPad:  1,
							rPad:  1,
						},
					},
					{
						box: box{
							width: 3,
							lPad:  1,
							rPad:  1,
						},
						limit: 3,
					},
					{
						box: box{
							width: 3,
							lPad:  1,
							rPad:  1,
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
				input:      test.fields.input,
				state:      &state,
				widthLimit: test.fields.widthLimit,
				output:     test.fields.output,
			}
			o.resolveWidths()
			got := want{
				metrics: state.columnMetrics,
			}
			testutil.AssertValue(t, got, test.want, "resolveWidths")
		})
	}
}

func Test_solver_fitColumns(t *testing.T) {
	type fields struct {
		input      compilerResult
		state      solverState
		widthLimit int
		output     solverResult
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
			name: "disabled automatic fitting",
			fields: fields{
				input: compilerResult{
					configResult: configResult{
						option: &option{},
					},
				},
				state: solverState{
					columnMetrics: []columnMetric{
						{
							box: box{
								width: 5,
							},
						},
					},
				},
				widthLimit: 1,
			},
			want: want{
				metrics: []columnMetric{
					{
						box: box{
							width: 5,
						},
					},
				},
			},
		},
		{
			name: "nonpositive width limit",
			fields: fields{
				input: compilerResult{
					configResult: configResult{
						option: &option{
							autoFit: true,
						},
					},
				},
				state: solverState{
					columnMetrics: []columnMetric{
						{
							box: box{
								width: 5,
							},
						},
					},
				},
			},
			want: want{
				metrics: []columnMetric{
					{
						box: box{
							width: 5,
						},
					},
				},
			},
		},
		{
			name: "empty columns",
			fields: fields{
				input: compilerResult{
					configResult: configResult{
						option: &option{
							autoFit: true,
						},
					},
				},
				widthLimit: 10,
			},
			want: want{},
		},
		{
			name: "natural width fits budget",
			fields: fields{
				input: compilerResult{
					configResult: configResult{
						option: &option{
							autoFit: true,
						},
					},
				},
				state: solverState{
					columnMetrics: []columnMetric{
						{
							box: box{
								width: 2,
							},
						},
						{
							box: box{
								width: 3,
							},
						},
					},
				},
				widthLimit: 10,
			},
			want: want{
				metrics: []columnMetric{
					{
						box: box{
							width: 2,
						},
					},
					{
						box: box{
							width: 3,
						},
					},
				},
			},
		},
		{
			name: "wide vertical glyphs consume budget",
			fields: fields{
				input: compilerResult{
					configResult: configResult{
						option: &option{
							autoFit: true,
							style: Style{
								Border: BorderStyle{
									Vertical: &Vertical{
										Outer: "界",
										Inner: "界",
									},
								},
							},
						},
					},
				},
				state: solverState{
					columnMetrics: []columnMetric{
						{
							box: box{
								width: 5,
							},
						},
						{
							box: box{
								width: 5,
							},
						},
					},
				},
				widthLimit: 10,
			},
			want: want{
				metrics: []columnMetric{
					{
						box: box{
							width: 2,
						},
						limit: 2,
					},
					{
						box: box{
							width: 2,
						},
						limit: 2,
					},
				},
			},
		},
		{
			name: "wide outer glyphs leave uneven budget",
			fields: fields{
				input: compilerResult{
					configResult: configResult{
						option: &option{
							autoFit: true,
							style: Style{
								Border: BorderStyle{
									Vertical: &Vertical{
										Outer: "界",
										Inner: "|",
									},
								},
							},
						},
					},
				},
				state: solverState{
					columnMetrics: []columnMetric{
						{
							box: box{
								width: 5,
							},
						},
						{
							box: box{
								width: 5,
							},
						},
					},
				},
				widthLimit: 10,
			},
			want: want{
				metrics: []columnMetric{
					{
						box: box{
							width: 3,
						},
						limit: 3,
					},
					{
						box: box{
							width: 2,
						},
						limit: 2,
					},
				},
			},
		},
		{
			name: "index and naturally narrow column stay fixed",
			fields: fields{
				input: compilerResult{
					configResult: configResult{
						option: &option{
							autoFit:     true,
							indexOffset: 1,
						},
					},
				},
				state: solverState{
					columnMetrics: []columnMetric{
						{
							box: box{
								width: 2,
							},
						},
						{
							box: box{
								width: 1,
							},
						},
						{
							box: box{
								width: 8,
							},
						},
					},
				},
				widthLimit: 8,
			},
			want: want{
				metrics: []columnMetric{
					{
						box: box{
							width: 2,
						},
						limit: 2,
					},
					{
						box: box{
							width: 1,
						},
						limit: 1,
					},
					{
						box: box{
							width: 5,
						},
						limit: 5,
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			o := &solver{
				input:      test.fields.input,
				state:      &state,
				widthLimit: test.fields.widthLimit,
				output:     test.fields.output,
			}
			o.fitColumns()
			got := want{
				metrics: state.columnMetrics,
			}
			testutil.AssertValue(t, got, test.want, "fitColumns")
		})
	}
}

func Test_solver_offsetColumns(t *testing.T) {
	type fields struct {
		input      compilerResult
		state      solverState
		widthLimit int
		output     solverResult
	}
	type want struct {
		offsets []int
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "adjacent boxes without vertical border",
			fields: fields{
				input: compilerResult{
					configResult: configResult{
						option: &option{},
					},
				},
				state: solverState{
					columnMetrics: []columnMetric{
						{
							box: box{
								width: 2,
								lPad:  1,
								rPad:  1,
							},
						},
						{
							box: box{
								width: 3,
								rPad:  2,
							},
						},
					},
				},
			},
			want: want{
				offsets: []int{0, 4},
			},
		},
		{
			name: "wide inner border separates boxes",
			fields: fields{
				input: compilerResult{
					configResult: configResult{
						option: &option{
							style: Style{
								Border: BorderStyle{
									Vertical: &Vertical{
										Inner: "界",
									},
								},
							},
						},
					},
				},
				state: solverState{
					columnMetrics: []columnMetric{
						{
							box: box{
								width: 2,
								lPad:  1,
								rPad:  1,
							},
						},
						{
							box: box{
								width: 3,
								rPad:  2,
							},
						},
					},
				},
			},
			want: want{
				offsets: []int{0, 6},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := test.fields.state
			o := &solver{
				input:      test.fields.input,
				state:      &state,
				widthLimit: test.fields.widthLimit,
				output:     test.fields.output,
			}
			o.offsetColumns()
			got := want{
				offsets: make([]int, len(state.columnMetrics)),
			}
			for i := range state.columnMetrics {
				got.offsets[i] = state.columnMetrics[i].box.offset
			}
			testutil.AssertValue(t, got, test.want, "offsetColumns")
		})
	}
}

func Test_box_totalWidth(t *testing.T) {
	type fields struct {
		offset int
		width  int
		lPad   int
		rPad   int
	}
	type want struct {
		width int
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "includes both paddings",
			fields: fields{
				offset: 9,
				width:  4,
				lPad:   2,
				rPad:   3,
			},
			want: want{
				width: 9,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &box{
				offset: test.fields.offset,
				width:  test.fields.width,
				lPad:   test.fields.lPad,
				rPad:   test.fields.rPad,
			}
			got := want{
				width: o.totalWidth(),
			}
			testutil.AssertValue(t, got, test.want, "totalWidth")
		})
	}
}
