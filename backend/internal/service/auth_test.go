package service

import (
	"context"
	"strings"
	"testing"
)

func TestValidateSIWEMessageStrictFields(t *testing.T) {
	ctx := context.Background()
	svc := NewAuthService(nil, nil, 11155111, "", "localhost:5173", "http://localhost:5173")
	address := "0x1111111111111111111111111111111111111111"

	_, message, err := svc.GenerateNonce(ctx, address)
	if err != nil {
		t.Fatalf("GenerateNonce() error = %v", err)
	}
	if _, err := svc.ValidateSIWEMessage(ctx, address, message); err != nil {
		t.Fatalf("ValidateSIWEMessage() error = %v", err)
	}

	badChain := strings.Replace(message, "Chain ID: 11155111", "Chain ID: 1", 1)
	if _, err := svc.ValidateSIWEMessage(ctx, address, badChain); err == nil {
		t.Fatal("ValidateSIWEMessage() accepted a message with the wrong chain id")
	}
}

func TestConsumeNonceAfterValidation(t *testing.T) {
	ctx := context.Background()
	svc := NewAuthService(nil, nil, 11155111, "", "localhost:5173", "http://localhost:5173")
	address := "0x2222222222222222222222222222222222222222"

	_, message, err := svc.GenerateNonce(ctx, address)
	if err != nil {
		t.Fatalf("GenerateNonce() error = %v", err)
	}
	if _, err := svc.ValidateSIWEMessage(ctx, address, message); err != nil {
		t.Fatalf("ValidateSIWEMessage() error = %v", err)
	}
	if _, err := svc.ValidateSIWEMessage(ctx, address, message); err != nil {
		t.Fatalf("ValidateSIWEMessage() should not consume nonce before signature verification: %v", err)
	}
	if err := svc.ConsumeNonce(ctx, address); err != nil {
		t.Fatalf("ConsumeNonce() error = %v", err)
	}
	if _, err := svc.ValidateSIWEMessage(ctx, address, message); err == nil {
		t.Fatal("ValidateSIWEMessage() accepted a consumed nonce")
	}
}
