package emoji

import "testing"

func TestRemove(t *testing.T) {
	type args struct {
		content string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "normal",
			args: args{
				content: "abc123",
			},
			want: "abc123",
		},
		{
			name: "chinese",
			args: args{
				content: "abc123大大泡泡糖🎏😊",
			},
			want: "abc123大大泡泡糖",
		},
		{
			name: "symbol",
			args: args{
				content: `abc123大大泡泡糖,./;'[]-=\~!@#$%^&*(){}|:"<>?，。、；’【】、·{}|：“》？《~！@#￥%……&*（）——+`,
			},
			want: `abc123大大泡泡糖,./;'[]-=\~!@#$%^&*(){}|:"<>?，。、；’【】、·{}|：“》？《~！@#￥%……&*（）——+`,
		},
		{
			name: "emoji",
			args: args{
				content: `abc123大大泡泡糖🎏😊`,
			},
			want: "abc123大大泡泡糖",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Remove(tt.args.content); got != tt.want {
				t.Errorf("Remove() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkRemove(b *testing.B) {
	content := string(make([]rune, 10000))
	for i := 0; i < b.N; i++ {
		result := Remove(content)
		if result != content {
			b.Fatal("NOT SAME")
		}
	}
}
