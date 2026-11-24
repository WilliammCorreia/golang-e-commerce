package ids

import (
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name   string
		prefix string
	}{
		{"produit", Product},
		{"panier", Basket},
		{"commande", Order},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := New(tt.prefix)

			if len(id) != 10 {
				t.Errorf("New(%q) = %q; longueur %d, attendu 10", tt.prefix, id, len(id))
			}

			if !strings.HasPrefix(id, tt.prefix+"-") {
				t.Errorf("New(%q) = %q; attendu %q", tt.prefix, id, tt.prefix+"-")
			}

			for _, c := range id[4:] {
				if !strings.Contains(alphabet, string(c)) {
					t.Errorf("New(%q) = %q; caractère %c non valide", tt.prefix, id, c)
				}
			}
		})
	}
}

func TestCode(t *testing.T) {
	result := Code()
	if len(result) != 6 {
		t.Error("longueur ne fait pas 6")
	}

	if result == Code() {
		t.Error("codes identiques")
	}
}
