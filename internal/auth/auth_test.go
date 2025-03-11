package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name     string
		args     args
		dontWant string
		wantErr  bool
	}{
		{
			name:     "test 1",
			args:     args{password: "password"},
			dontWant: "password",
			wantErr:  false,
		},
		{
			name:     "test 2",
			args:     args{password: "1234pass"},
			dontWant: "1234pass",
			wantErr:  false,
		},
		{
			name:     "test 3",
			args:     args{password: "1234qwerty"},
			dontWant: "1234qwerty",
			wantErr:  false,
		},
		{
			name:     "test 4",
			args:     args{password: ""},
			dontWant: "has to be different from the input",
			wantErr:  true,
		},
		{
			name:     "test 5",
			args:     args{password: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
			dontWant: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HashPassword(tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				err = CheckPasswordHash(tt.args.password, got)
				if err != nil {
					t.Errorf("Generated hash does not match the original password value: %v", err)
				}
			}
			if got == tt.dontWant {
				t.Errorf("HashPassword() = %v, password = %v. Should not be the same", got, tt.dontWant)
			}
		})
	}
}
