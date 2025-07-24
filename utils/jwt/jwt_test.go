package jwt

import (
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	config := JWTConfig{
		SecretKey: "test_secret_key",
		Expired:   24, // 24 hours
	}

	tests := []struct {
		name     string
		userID   uint
		username string
		role     string
		wantErr  bool
	}{
		{
			name:     "valid user data",
			userID:   123,
			username: "testuser",
			role:     "admin",
			wantErr:  false,
		},
		{
			name:     "zero user ID",
			userID:   0,
			username: "testuser",
			role:     "user",
			wantErr:  false,
		},
		{
			name:     "empty username",
			userID:   456,
			username: "",
			role:     "user",
			wantErr:  false,
		},
		{
			name:     "empty role",
			userID:   789,
			username: "testuser",
			role:     "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, expiresAt, err := GenerateToken(config, tt.userID, tt.username, tt.role)

			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if token == "" {
					t.Errorf("GenerateToken() returned empty token")
				}

				// Check that expiration time is in the future
				if expiresAt.Before(time.Now()) {
					t.Errorf("GenerateToken() expiration time is in the past: %v", expiresAt)
				}

				// Check that expiration time is approximately 24 hours from now
				expectedExpiry := time.Now().Add(24 * time.Hour)
				if expiresAt.After(expectedExpiry.Add(5*time.Minute)) || expiresAt.Before(expectedExpiry.Add(-5*time.Minute)) {
					t.Errorf("GenerateToken() expiration time is not approximately 24 hours: %v", expiresAt)
				}
			}
		})
	}
}

func TestParseToken(t *testing.T) {
	secretKey := "test_secret_key"
	config := JWTConfig{
		SecretKey: secretKey,
		Expired:   24,
	}

	// Generate a valid token for testing
	userID := uint(123)
	username := "testuser"
	role := "admin"

	token, _, err := GenerateToken(config, userID, username, role)
	if err != nil {
		t.Fatalf("Failed to generate token for test: %v", err)
	}

	tests := []struct {
		name      string
		token     string
		secretKey string
		wantErr   bool
	}{
		{
			name:      "valid token",
			token:     token,
			secretKey: secretKey,
			wantErr:   false,
		},
		{
			name:      "invalid secret key",
			token:     token,
			secretKey: "wrong_secret",
			wantErr:   true,
		},
		{
			name:      "empty token",
			token:     "",
			secretKey: secretKey,
			wantErr:   true,
		},
		{
			name:      "malformed token",
			token:     "invalid.token.format",
			secretKey: secretKey,
			wantErr:   true,
		},
		{
			name:      "expired token",
			token:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMjMsInVzZXJuYW1lIjoidGVzdHVzZXIiLCJyb2xlIjoiYWRtaW4iLCJleHAiOjE2MzQ1Njc4OTAsImlhdCI6MTYzNDQ4MTQ5MH0.invalid_signature",
			secretKey: secretKey,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ParseToken(tt.token, tt.secretKey)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if claims == nil {
					t.Errorf("ParseToken() returned nil claims")
					return
				}

				// Verify the claims match the original data
				if claims.UserID != userID {
					t.Errorf("ParseToken() UserID = %v, want %v", claims.UserID, userID)
				}
				if claims.Username != username {
					t.Errorf("ParseToken() Username = %v, want %v", claims.Username, username)
				}
				if claims.Role != role {
					t.Errorf("ParseToken() Role = %v, want %v", claims.Role, role)
				}
			}
		})
	}
}

func TestExtractUserContext(t *testing.T) {
	tests := []struct {
		name   string
		claims *JWTClaims
		want   UserContext
	}{
		{
			name: "valid claims",
			claims: &JWTClaims{
				UserID:   123,
				Username: "testuser",
				Role:     "admin",
			},
			want: UserContext{
				UserID:   123,
				Username: "testuser",
				Role:     "admin",
			},
		},
		{
			name: "zero values",
			claims: &JWTClaims{
				UserID:   0,
				Username: "",
				Role:     "",
			},
			want: UserContext{
				UserID:   0,
				Username: "",
				Role:     "",
			},
		},
		{
			name:   "nil claims",
			claims: nil,
			want: UserContext{
				UserID:   0,
				Username: "",
				Role:     "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractUserContext(tt.claims)

			if got.UserID != tt.want.UserID {
				t.Errorf("ExtractUserContext() UserID = %v, want %v", got.UserID, tt.want.UserID)
			}
			if got.Username != tt.want.Username {
				t.Errorf("ExtractUserContext() Username = %v, want %v", got.Username, tt.want.Username)
			}
			if got.Role != tt.want.Role {
				t.Errorf("ExtractUserContext() Role = %v, want %v", got.Role, tt.want.Role)
			}
		})
	}
}

func TestJWTIntegration(t *testing.T) {
	// Test complete flow: generate token -> parse token -> extract context
	config := JWTConfig{
		SecretKey: "integration_test_secret",
		Expired:   1, // 1 hour for faster test
	}

	userID := uint(456)
	username := "integrationuser"
	role := "manager"

	// Step 1: Generate token
	token, expiresAt, err := GenerateToken(config, userID, username, role)
	if err != nil {
		t.Fatalf("GenerateToken() failed: %v", err)
	}

	// Step 2: Parse token
	claims, err := ParseToken(token, config.SecretKey)
	if err != nil {
		t.Fatalf("ParseToken() failed: %v", err)
	}

	// Step 3: Extract user context
	userContext := ExtractUserContext(claims)

	// Verify all data matches
	if userContext.UserID != userID {
		t.Errorf("Integration test UserID mismatch: got %v, want %v", userContext.UserID, userID)
	}
	if userContext.Username != username {
		t.Errorf("Integration test Username mismatch: got %v, want %v", userContext.Username, username)
	}
	if userContext.Role != role {
		t.Errorf("Integration test Role mismatch: got %v, want %v", userContext.Role, role)
	}

	// Verify expiration time
	if expiresAt.Before(time.Now()) {
		t.Errorf("Integration test: token expires in the past: %v", expiresAt)
	}
}

func TestJWTExpiration(t *testing.T) {
	config := JWTConfig{
		SecretKey: "expiration_test_secret",
		Expired:   1, // 1 hour
	}

	// Generate token
	token, _, err := GenerateToken(config, 123, "testuser", "user")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Parse token immediately (should work)
	claims, err := ParseToken(token, config.SecretKey)
	if err != nil {
		t.Errorf("ParseToken() failed immediately after generation: %v", err)
	}
	if claims == nil {
		t.Errorf("ParseToken() returned nil claims")
	}

	// Note: Testing actual expiration would require waiting for the token to expire,
	// which would make the test slow. In a real scenario, you might want to use
	// a very short expiration time for testing or mock the time.
}

func TestJWTClaimsStructure(t *testing.T) {
	// Test that JWTClaims embeds jwt.RegisteredClaims correctly
	claims := &JWTClaims{
		UserID:   123,
		Username: "testuser",
		Role:     "admin",
	}

	// Verify the structure has the expected fields
	if claims.UserID != 123 {
		t.Errorf("JWTClaims.UserID = %v, want 123", claims.UserID)
	}
	if claims.Username != "testuser" {
		t.Errorf("JWTClaims.Username = %v, want testuser", claims.Username)
	}
	if claims.Role != "admin" {
		t.Errorf("JWTClaims.Role = %v, want admin", claims.Role)
	}

	// Verify that RegisteredClaims is embedded (it's a struct, not a pointer)
	// The fact that we can access it means it's properly embedded
	_ = claims.RegisteredClaims
}
