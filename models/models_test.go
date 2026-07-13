package models

import (
	"testing"
	"time"
)

func TestUserCreation(t *testing.T) {
	user := User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	if user.Name != "Test User" {
		t.Fatalf("Expected 'Test User', got '%s'", user.Name)
	}
	if user.Email != "test@example.com" {
		t.Fatalf("Expected 'test@example.com', got '%s'", user.Email)
	}
}

func TestTransactionValidation(t *testing.T) {
	tx := Transaction{
		UserID:   1,
		Type:     "expense",
		Amount:   100.50,
		Category: "Food",
		Date:     time.Now(),
	}

	if tx.Type != "expense" {
		t.Fatalf("Expected 'expense', got '%s'", tx.Type)
	}
	if tx.Amount <= 0 {
		t.Fatal("Expected positive amount")
	}
}

func TestBudgetResponse(t *testing.T) {
	resp := BudgetResponse{
		ID:        1,
		Category:  "Groceries",
		Budget:    500,
		Spent:     300,
		Remaining: 200,
		Progress:  60,
		Status:    "On Track",
	}

	if resp.Remaining != resp.Budget-resp.Spent {
		t.Fatal("Remaining should equal Budget minus Spent")
	}
}

func TestSavingsGoalResponse(t *testing.T) {
	resp := SavingsGoalResponse{
		ID:           1,
		Title:        "Vacation",
		TargetAmount: 5000,
		SavedAmount:  2500,
		Progress:     50,
		Status:       "On Track",
	}

	if resp.Progress != (resp.SavedAmount/resp.TargetAmount)*100 {
		t.Fatal("Progress should be (SavedAmount/TargetAmount)*100")
	}
}
