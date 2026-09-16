package auth

import (
	"latihan-fiber/app/model"
	"testing"
)

func TestValidateRegister(t *testing.T) {
	tests := []struct {
		name     string
		req      model.RegisterRequest
		wantErr  bool
		errField string
	}{
		{
			name: "Valid Request",
			req: model.RegisterRequest{
				Username: "user123",
				Email:    "test@example.com",
				Password: "StrongPassword1",
			},
			wantErr: false,
		},
		{
			name: "Username Too Short",
			req: model.RegisterRequest{
				Username: "us",
				Email:    "test@example.com",
				Password: "StrongPassword1",
			},
			wantErr:  true,
			errField: "username",
		},
		{
			name: "Username Invalid Characters",
			req: model.RegisterRequest{
				Username: "user 123",
				Email:    "test@example.com",
				Password: "StrongPassword1",
			},
			wantErr:  true,
			errField: "username",
		},
		{
			name: "Invalid Email",
			req: model.RegisterRequest{
				Username: "user123",
				Email:    "testexample.com",
				Password: "StrongPassword1",
			},
			wantErr:  true,
			errField: "email",
		},
		{
			name: "Password Too Weak (No Number)",
			req: model.RegisterRequest{
				Username: "user123",
				Email:    "test@example.com",
				Password: "StrongPassword",
			},
			wantErr:  true,
			errField: "password",
		},
		{
			name: "Password Weak Dictionary",
			req: model.RegisterRequest{
				Username: "user123",
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr:  true,
			errField: "password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := ValidateRegister(tt.req)
			hasErr := len(errs) > 0
			if hasErr != tt.wantErr {
				t.Errorf("ValidateRegister() got error = %v, wantErr %v", hasErr, tt.wantErr)
			}
			if tt.wantErr && tt.errField != "" {
				if _, ok := errs[tt.errField]; !ok {
					t.Errorf("ValidateRegister() expected error in field %s, got errors %v", tt.errField, errs)
				}
			}
		})
	}
}

func TestValidateLogin(t *testing.T) {
	tests := []struct {
		name     string
		req      model.LoginRequest
		wantErr  bool
		errField string
	}{
		{
			name: "Valid Login",
			req: model.LoginRequest{
				Username: "user123",
				Password: "anyPassword123",
			},
			wantErr: false,
		},
		{
			name: "Missing Username",
			req: model.LoginRequest{
				Username: "",
				Password: "anyPassword123",
			},
			wantErr:  true,
			errField: "username",
		},
		{
			name: "Missing Password",
			req: model.LoginRequest{
				Username: "user123",
				Password: "",
			},
			wantErr:  true,
			errField: "password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := ValidateLogin(tt.req)
			hasErr := len(errs) > 0
			if hasErr != tt.wantErr {
				t.Errorf("ValidateLogin() got error = %v, wantErr %v", hasErr, tt.wantErr)
			}
			if tt.wantErr && tt.errField != "" {
				if _, ok := errs[tt.errField]; !ok {
					t.Errorf("ValidateLogin() expected error in field %s, got errors %v", tt.errField, errs)
				}
			}
		})
	}
}
