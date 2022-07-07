package tm

import "testing"

func TestConfig_validate(t *testing.T) {
	type fields struct {
		Force      bool
		FilePath   string
		MaxRetries uint8
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "file path is empty",
			fields: fields{
				FilePath:   "",
				MaxRetries: 0,
			},
			wantErr: true,
		},
		{
			name: "max retries = 0",
			fields: fields{
				FilePath:   "some",
				MaxRetries: 0,
			},
			wantErr: true,
		},
		{
			name: "max retries = 1",
			fields: fields{
				FilePath:   "some",
				MaxRetries: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Force:      tt.fields.Force,
				FilePath:   tt.fields.FilePath,
				MaxRetries: tt.fields.MaxRetries,
			}
			if err := c.validate(); (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
