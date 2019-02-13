package gosignedpayapi

import "testing"

func Test_getSignature(t *testing.T) {
	type args struct {
		body []byte
		m    []byte
		pk   []byte
	}

	args1 := args{
		[]byte(`{"name":"case"`),
		[]byte(`merchant`),
		[]byte(`secret`),
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"working case",
			args1,
			"NTM5NmYzMjBjMTU5NDE3MjE1NTY1ZmIxODUwMjVhZjFhZjI3ZDBmNjgyY2NhZjY5ZTBjZGFmYzBmNGQ3NzEwOWM1ZmUxNGQ3ODU0ODRlNDI0Y2ZiNzkzY2I2ZjAyZmFiZTNhMzk2NjNlODVkYzllZDA4ZTVjZTFmN2ZmMWY2ZGQ=",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSignature(tt.args.body, tt.args.m, tt.args.pk); got != tt.want {
				t.Errorf("getSignature() = %v, want %v", got, tt.want)
			}
		})
	}
}
