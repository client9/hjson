package cson

import (
	"testing"
)

func TestCSON(t *testing.T) {
	cases := []struct {
		orig string
		want string
	}{
		{ // 1
			orig: `
foo:
  brother:
    name: "max"
    age: 11
  sister:
    name: "ida"
    age: 9
`,

			want: `{"foo":{"brother":{"name":"max","age":11},"sister":{"name":"ida","age":9}}}`,
		},
		{ // 2
			orig: `
bits: [
   1,0,1
   0,1,0
   1,1,1
]
`,
			want: `{"bits":[1,0,1,0,1,0,1,1,1]}`,
		},
		{ // 3
			orig: `foo: true`,
			want: `{"foo":true}`,
		},
		{ // 4
			orig: `foo: "true"`,
			want: `{"foo":"true"}`,
		},
		{ // 5
			orig: `
# comments!!

# An Array with no commas!
greatDocumentaries: [
   'earthlings.com'
    # love it
   'forksoverknives.com'
   'cowspiracy.com'
]
`,
			want: `{"greatDocumentaries":["earthlings.com","forksoverknives.com","cowspiracy.com"]}`,
		},
		{ // 7
			orig: `
######## banner!
foo: "bar"
`,
			want: `{"foo":"bar"}`,
		},
		{ // 8
			orig: `
###
multiline
###
foo: "bar"
`,
			want: `{"foo":"bar"}`,
		},
	}

	for num, tt := range cases {
		got := ToJSON([]byte(tt.orig))
		if tt.want != string(got) {
			t.Errorf("%d: want %s got %s", num+1, tt.want, got)
		}
	}
}
