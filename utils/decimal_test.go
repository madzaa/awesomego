package utils

import "testing"

func TestFormatBalance(t *testing.T) {
	tests := []struct {
		name    string
		balance uint64
		want    string
	}{
		{
			name:    "format whole number",
			balance: 100,
			want:    "1.00",
		},
		{
			name:    "format with cents",
			balance: 123,
			want:    "1.23",
		},
		{
			name:    "format zero",
			balance: 0,
			want:    "0.00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatBalance(tt.balance); got != tt.want {
				t.Errorf("FormatBalance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseAmount(t *testing.T) {
	tests := []struct {
		name    string
		amount  string
		want    uint64
		wantErr bool
	}{
		{
			name:    "valid amount",
			amount:  "1.23",
			want:    123,
			wantErr: false,
		},
		{
			name:    "invalid format",
			amount:  "1.2",
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid chars",
			amount:  "abc",
			want:    0,
			wantErr: true,
		},
		{
			name:    "empty string",
			amount:  "",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAmount(tt.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAmount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseAmount() = %v, want %v", got, tt.want)
			}
		})
	}
}
