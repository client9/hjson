package hjson

import (
	"testing"
)

func TestHJSON(t *testing.T) {
	cases := []struct {
		orig string
		want string
	}{
		{ // 1
			orig: `{"foo": "bar", }`,
			want: `{"foo":"bar"}`,
		},
		{ // 2
			orig: `  "foo": "bar",  `,
			want: `{"foo":"bar"}`,
		},
		{ // 3
			orig: `  "foo": "bar"  `,
			want: `{"foo":"bar"}`,
		},
		{ // 4
			orig: `{foo: bar}`,
			want: `{"foo":"bar"}`,
		},
		{ // 5
			orig: `{foo: bar bar}`,
			want: `{"foo":"bar bar"}`,
		},
		{ // 6
			orig: `{"foo" : "bar", "ding": "bat"}`,
			want: `{"foo":"bar","ding":"bat"}`,
		},
		{ // 7
			orig: `{ foo : bar, ding: bat}`,
			want: `{"foo":"bar","ding":"bat"}`,
		},
		{ // 8
			orig: `{foo:  false  }`,
			want: `{"foo":false}`,
		},
		{ // 9
			orig: `{ foo:  [1,2,3,4]  }`,
			want: `{"foo":[1,2,3,4]}`,
		},
		{ // 10
			orig: `{ foo:  [ 1 , 2 , 3 , 4 ]  }`,
			want: `{"foo":[1,2,3,4]}`,
		},
		{ // 11
			orig: "{ foo:  [  \n 1 \n 2  \n 3  \n  4 \n]\n}",
			want: `{"foo":[1,2,3,4]}`,
		},
		{ // 12
			orig: `{ foo:  [ "1", "2", "3", "4",  ]  }`,
			want: `{"foo":["1","2","3","4"]}`,
		},
		{ // 13
			orig: `{ 日本語:  [ "1", "2", "3", "4",  ]  }`,
			want: `{"日本語":["1","2","3","4"]}`,
		},
		{ // 14
			orig: `
# junk
foo: "bar",
`,
			want: `{"foo":"bar"}`,
		},
		{ // 15
			orig: `
foo: '''
bar
''',
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
