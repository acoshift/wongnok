package auth

import "testing"

func TestPassword(t *testing.T) {
	// Test Table
	cases := []struct {
		Name            string
		Password        string
		ComparePassword string
		Equal           bool
	}{
		{"Sucess", "superman", "superman", true},
		{"Compare with empty password", "superman", "", false},
		{"Compare with invalid password", "superman", "batman", false},
	}

	for _, tC := range cases {
		t.Run(tC.Name, func(t *testing.T) {
			hashed := hashPassword(tC.Password)
			if hashed == "" {
				t.Errorf("expected non-empty string from hashPassword; got empty")
			}
			if compareHashAndPassword(hashed, tC.ComparePassword) != tC.Equal {
				if tC.Equal {
					t.Errorf("expected hashed and password are equal; got not equal")
				} else {
					t.Errorf("expected hashed and password are not equal; got equal")
				}
			}
		})
	}
}
