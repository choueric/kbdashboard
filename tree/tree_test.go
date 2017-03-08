package tree

import (
	"bytes"
	"testing"
)

type Val struct {
	name string
}

func Test_PrintTree(t *testing.T) {
	expect := `root
├── level1_1
│   └── level2_1
├── level1_2
│   └── level2_2
│       └── level3_1
├── level1_3
└── level1_4
`

	root := New(&Val{"root"})

	level1_1 := root.AddSubNode(&Val{"level1_1"})
	level1_2 := root.AddSubNode(&Val{"level1_2"})
	root.AddSubNode(&Val{"level1_3"})
	root.AddSubNode(&Val{"level1_4"})

	level1_1.AddSubNode(&Val{"level2_1"})
	level2_2 := level1_2.AddSubNode(&Val{"level2_2"})
	level2_2.AddSubNode(&Val{"level3_1"})

	var result bytes.Buffer
	root.PrintTree(&result, func(n *Node) string { return n.Value.(*Val).name })

	if result.String() != expect {
		t.Error("tree graphic does not match.")
		t.Error(result.String())
	}
}
