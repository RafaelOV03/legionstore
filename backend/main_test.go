package main

import "testing"

func TestSimple(t *testing.T) {
	resultado := 1 + 1
	if resultado != 2 {
		t.Errorf("Matemática básica falló: se esperaba 2, se obtuvo %d", resultado)
	}
}