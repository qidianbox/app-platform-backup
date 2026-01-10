package validator

import (
	"testing"
)

func TestValidateAppCreate(t *testing.T) {
	tests := []struct {
		name    string
		req     *AppCreateRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request with name",
			req: &AppCreateRequest{
				Name:        "测试应用",
				PackageName: "com.example.app",
				Description: "这是一个测试应用",
			},
			wantErr: false,
		},
		{
			name: "valid request with app_name",
			req: &AppCreateRequest{
				AppName:     "测试应用2",
				PackageName: "com.example.app2",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			req: &AppCreateRequest{
				Name:        "",
				PackageName: "com.example.app",
			},
			wantErr: true,
			errMsg:  "应用名称不能为空",
		},
		{
			name: "name too short",
			req: &AppCreateRequest{
				Name: "A",
			},
			wantErr: true,
			errMsg:  "应用名称至少需要2个字符",
		},
		{
			name: "name too long",
			req: &AppCreateRequest{
				Name: "这是一个非常非常非常非常非常非常非常非常非常非常非常非常非常非常非常非常非常非常非常非常非常非常非常非常非常非常非常非常非常长的应用名称",
			},
			wantErr: true,
			errMsg:  "应用名称不能超过50个字符",
		},
		{
			name: "invalid package name",
			req: &AppCreateRequest{
				Name:        "测试应用",
				PackageName: "invalid",
			},
			wantErr: true,
			errMsg:  "包名格式不正确，应为类似 com.example.app 的格式",
		},
		{
			name: "valid package name",
			req: &AppCreateRequest{
				Name:        "测试应用",
				PackageName: "com.example.myapp",
			},
			wantErr: false,
		},
		{
			name: "description too long",
			req: &AppCreateRequest{
				Name:        "测试应用",
				Description: string(make([]byte, 501)),
			},
			wantErr: true,
			errMsg:  "描述不能超过500个字符",
		},
		{
			name: "too many modules",
			req: &AppCreateRequest{
				Name:    "测试应用",
				Modules: make([]string, 21),
			},
			wantErr: true,
			errMsg:  "启用的模块数量不能超过20个",
		},
		{
			name: "name with special characters",
			req: &AppCreateRequest{
				Name: "测试<script>应用",
			},
			wantErr: true,
			errMsg:  "应用名称不能包含特殊字符",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAppCreate(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAppCreate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && tt.errMsg != "" {
				if err.Error() != tt.errMsg {
					t.Errorf("ValidateAppCreate() error message = %v, want %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestValidateAppUpdate(t *testing.T) {
	tests := []struct {
		name    string
		req     *AppUpdateRequest
		wantErr bool
	}{
		{
			name: "valid update",
			req: &AppUpdateRequest{
				Name:        "新名称",
				Description: "新描述",
			},
			wantErr: false,
		},
		{
			name: "empty update (valid)",
			req:     &AppUpdateRequest{},
			wantErr: false,
		},
		{
			name: "invalid status",
			req: &AppUpdateRequest{
				Status: intPtr(2),
			},
			wantErr: true,
		},
		{
			name: "valid status 0",
			req: &AppUpdateRequest{
				Status: intPtr(0),
			},
			wantErr: false,
		},
		{
			name: "valid status 1",
			req: &AppUpdateRequest{
				Status: intPtr(1),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAppUpdate(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAppUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		want    uint
		wantErr bool
	}{
		{
			name:    "valid id",
			id:      "123",
			want:    123,
			wantErr: false,
		},
		{
			name:    "empty id",
			id:      "",
			want:    0,
			wantErr: true,
		},
		{
			name:    "zero id",
			id:      "0",
			want:    0,
			wantErr: true,
		},
		{
			name:    "negative id",
			id:      "-1",
			want:    0,
			wantErr: true,
		},
		{
			name:    "non-numeric id",
			id:      "abc",
			want:    0,
			wantErr: true,
		},
		{
			name:    "id with leading zero",
			id:      "01",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidatePagination(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		size     int
		wantPage int
		wantSize int
	}{
		{
			name:     "valid pagination",
			page:     1,
			size:     10,
			wantPage: 1,
			wantSize: 10,
		},
		{
			name:     "negative page",
			page:     -1,
			size:     10,
			wantPage: 1,
			wantSize: 10,
		},
		{
			name:     "zero page",
			page:     0,
			size:     10,
			wantPage: 1,
			wantSize: 10,
		},
		{
			name:     "negative size",
			page:     1,
			size:     -1,
			wantPage: 1,
			wantSize: 10,
		},
		{
			name:     "zero size",
			page:     1,
			size:     0,
			wantPage: 1,
			wantSize: 10,
		},
		{
			name:     "size too large",
			page:     1,
			size:     200,
			wantPage: 1,
			wantSize: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPage, gotSize := ValidatePagination(tt.page, tt.size)
			if gotPage != tt.wantPage {
				t.Errorf("ValidatePagination() page = %v, want %v", gotPage, tt.wantPage)
			}
			if gotSize != tt.wantSize {
				t.Errorf("ValidatePagination() size = %v, want %v", gotSize, tt.wantSize)
			}
		})
	}
}

func TestSanitizeString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "normal string",
			input: "Hello World",
			want:  "Hello World",
		},
		{
			name:  "string with HTML tags",
			input: "<script>alert('xss')</script>",
			want:  "alert(&#39;xss&#39;)",
		},
		{
			name:  "string with special chars",
			input: "Test & \"quote\"",
			want:  "Test &amp; &quot;quote&quot;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeString(tt.input)
			if got != tt.want {
				t.Errorf("SanitizeString() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function
func intPtr(i int) *int {
	return &i
}
