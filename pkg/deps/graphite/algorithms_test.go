package graphite

import (
	"reflect"
	"testing"
)

func generateCyclicGraph() *Graph[string] {
	g, _ := New([]string{"a", "b", "c"})
	g.AddDependency("a", "b")
	g.AddDependency("b", "c")
	g.AddDependency("c", "a")
	return g
}

func Test_dfs(t *testing.T) {
	graph := generateCyclicGraph()
	type args struct {
		node    int
		reverse bool
	}
	tests := []struct {
		name    string
		args    args
		want    []int
		wantErr bool
	}{
		{
			name: "invalid cyclic graph",
			args: args{
				node:    0,
				reverse: false,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			visited := make([]bool, graph.size)
			got, err := dfs(tt.args.node, visited, graph, tt.args.reverse)
			if (err != nil) != tt.wantErr {
				t.Errorf("dfs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dfs() = %v, want %v", got, tt.want)
			}
		})
	}
}
