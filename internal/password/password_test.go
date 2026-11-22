package password

import (
	"testing"
)

func TestCheck(t *testing.T) {
	hashedPwd, err := Hash("secret123")
	if err != nil {
		t.Fatalf("Hash('secret123') = %v; attendu secret123 hashé", err)
	}

	tests := []struct {
		name     string
		password string
		stored   string
		expected bool
	}{
		{
			name:     "Mot de passe correct",
			password: "secret123",
			stored:   hashedPwd,
			expected: true,
		},
		{
			name:     "Mot de passe incorrect",
			password: "mauvaisMotDePasse",
			stored:   hashedPwd,
			expected: false,
		},
		{
			name:     "Mot de passe vide",
			password: "",
			stored:   hashedPwd,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Check(tt.password, tt.stored)
			if result != tt.expected {
				t.Errorf("Check('%s', '%s') = %t; attendu %t", tt.password, tt.stored, result, tt.expected)
			}
		})
	}
}

func TestHash(t *testing.T) {
	firstHash, err := Hash("secret123")
	if err != nil {
		t.Fatalf("Hash('secret123') = %v; attendu secret123 hashé", err)
	}

	secondHash, err := Hash("secret123")
	if err != nil {
		t.Fatalf("Hash('secret123') = %v; attendu secret123 hashé", err)
	}

	if firstHash == secondHash {
		t.Errorf("même mot de passe haché deux fois donne le même résultat : %s", firstHash)
	}
}
