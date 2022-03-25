package graphite

import (
	"reflect"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func setupSimpleGraphForTest() *Graph[string] {
	// a>b>c
	simpleGraph, _ := New([]string{"a", "b", "c"})
	simpleGraph.AddDependency("a", "b")
	simpleGraph.AddDependency("b", "c")
	return simpleGraph
}

func setupNoDepsGraphForTest() *Graph[string] {
	// a b c
	noDepsGraph, _ := New([]string{"a", "b", "c"})
	return noDepsGraph
}

func setupComplexGraphForTest() *Graph[string] {
	// a>b>c<d e>f
	complexGraph, _ := New([]string{"a", "b", "c", "d", "e", "f"})
	complexGraph.AddDependency("a", "b")
	complexGraph.AddDependency("b", "c")
	complexGraph.AddDependency("d", "c")
	complexGraph.AddDependency("e", "f")
	return complexGraph
}

func TestGraph_GetNodeDependentBranches(t *testing.T) {
	simpleGraph := setupSimpleGraphForTest()
	noDepsGraph := setupNoDepsGraphForTest()
	complexGraph := setupComplexGraphForTest()
	tests := []struct {
		name    string
		graph   *Graph[string]
		node    string
		want    [][]string
		wantErr bool
	}{
		{
			name:    "a->b->(c)",
			graph:   simpleGraph,
			node:    "c",
			want:    [][]string{{"a", "b"}},
			wantErr: false,
		},
		{
			name:    "a->(b)->c",
			graph:   simpleGraph,
			node:    "b",
			want:    [][]string{{"a"}},
			wantErr: false,
		},
		{
			name:    "(a)->b->c",
			graph:   simpleGraph,
			node:    "a",
			want:    [][]string{},
			wantErr: false,
		},
		{
			name:    "a|b|(c)",
			graph:   noDepsGraph,
			node:    "c",
			want:    [][]string{},
			wantErr: false,
		},
		{
			name:  "a->b->(c)<-d e->f",
			graph: complexGraph,
			node:  "c",
			want:  [][]string{{"a", "b"}, {"d"}},
		},
		{
			name:  "a->b->c<-d e->(f)",
			graph: complexGraph,
			node:  "f",
			want:  [][]string{{"e"}},
		},
		{
			name:    "a->b->c<-d e->f (invalid)",
			graph:   complexGraph,
			node:    "z",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graph := tt.graph
			got, err := graph.GetDepsFor(tt.node)
			if (err != nil) != tt.wantErr {
				t.Errorf("Graph.GetNodeDependentBranches() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Graph.GetNodeDependentBranches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGraph_Connect(t *testing.T) {
	type args struct {
		a string
		b string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "a->b",
			args: args{
				a: "a",
				b: "b",
			},
			wantErr: false,
		},
		{
			name: "z->b",
			args: args{
				a: "z",
				b: "b",
			},
			wantErr: true,
		},
		{
			name: "z->y",
			args: args{
				a: "z",
				b: "y",
			},
			wantErr: true,
		},
		{
			name: "a->y",
			args: args{
				a: "a",
				b: "y",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graph := setupNoDepsGraphForTest()
			if err := graph.AddDependency(tt.args.a, tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("Graph.Connect() error = %v, wantErr %v", err, tt.wantErr)
			}
			indexA := lo.IndexOf(graph.reference, tt.args.a)
			indexB := lo.IndexOf(graph.reference, tt.args.b)
			if tt.wantErr == false {
				assert.NotEqual(t, indexA, -1)
				assert.NotEqual(t, indexB, -1)
				assert.True(t, graph.reversed[indexB][indexA])
				assert.True(t, graph.connections[indexA][indexB])
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		nodes   []string
		wantErr bool
	}{
		{
			name:    "invalid repeated nodes",
			nodes:   []string{"a", "a"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.nodes)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
