package stack

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	cases := []struct {
		stack    []byte
		err      bool
		expected Stack
	}{
		{
			stack: []byte(`goroutine 6 [running]:
runtime/debug.Stack()
	/opt/go/src/runtime/debug/stack.go:24 +0x88
github.com/morelj/httptools/stack.TestParse(0x4000106ea0)
	/home/ubuntu/workspaces/github.com/morelj/httptools/stack/stack_test.go:13 +0x38
testing.tRunner(0x4000106ea0, 0x1c5f50)
	/opt/go/src/testing/testing.go:1259 +0xf8
created by testing.(*T).Run
	/opt/go/src/testing/testing.go:1306 +0x350`),
			expected: Stack{
				Goroutines: []Goroutine{
					{
						Name:  "goroutine 6",
						State: "running",
						Elements: []Element{
							{
								Func:   "runtime/debug.Stack()",
								Source: "/opt/go/src/runtime/debug/stack.go:24 +0x88",
							},
							{
								Func:   "github.com/morelj/httptools/stack.TestParse(0x4000106ea0)",
								Source: "/home/ubuntu/workspaces/github.com/morelj/httptools/stack/stack_test.go:13 +0x38",
							},
							{
								Func:   "testing.tRunner(0x4000106ea0, 0x1c5f50)",
								Source: "/opt/go/src/testing/testing.go:1259 +0xf8",
							},
							{
								Func:   "created by testing.(*T).Run",
								Source: "/opt/go/src/testing/testing.go:1306 +0x350",
							},
						},
					},
				},
			},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			stack, err := Parse(c.stack)
			if c.err {
				assert.Error(err)
			} else {
				require.NoError(err)
				assert.Equal(c.expected, stack)
			}
		})
	}
}
