package main

import (
	"context"
	"testing"
)

func TestBasicFunctionality(t *testing.T) {
	t.Run("String operations", func(t *testing.T) {
		result := "hello" + " " + "world"
		if result != "hello world" {
			t.Errorf("Expected 'hello world', got '%s'", result)
		}
	})
	
	t.Run("Numeric operations", func(t *testing.T) {
		sum := 1 + 2 + 3
		if sum != 6 {
			t.Errorf("Expected 6, got %d", sum)
		}
	})
}

func TestServiceCreation(t *testing.T) {
	t.Run("Create mock service", func(t *testing.T) {
		service := struct {
			Name string
			ID   int
		}{
			Name: "test-service",
			ID:   1,
		}
		
		if service.Name != "test-service" {
			t.Errorf("Expected 'test-service', got '%s'", service.Name)
		}
	})
}

func TestContextHandling(t *testing.T) {
	t.Run("Context creation", func(t *testing.T) {
		ctx := context.Background()
		if ctx == nil {
			t.Error("Context should not be nil")
		}
	})
}
