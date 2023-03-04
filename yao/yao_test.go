package yao

import (
	"reflect"
	"testing"
)

func TestYaoProcess(t *testing.T) {
	type args struct {
		method string
		args   []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{
				method: "scripts.ai.chatgpt.GetSetting",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := YaoProcess(tt.args.method, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("YaoProcess() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("YaoProcess() = %v, want %v", got, tt.want)
			}
		})
	}
}
