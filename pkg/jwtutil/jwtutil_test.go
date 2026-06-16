package jwtutil_test

import (
	"testing"
	"time"

	"github.com/denidarta/dauth/pkg/jwtutil"
)

const (
	testPrivKeyPath = "../../keys/private.pem"
	testPubKeyPath  = "../../keys/public.pem"
)

func newTestManager(t *testing.T) *jwtutil.Manager {
	t.Helper()
	m, err := jwtutil.NewManager(testPrivKeyPath, testPubKeyPath)
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	return m
}

func TestSignAndVerify(t *testing.T) {
	m := newTestManager(t)

	token, err := m.Sign("user-123", "product-456", "admin", time.Minute)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}

	claims, err := m.Verify(token)
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}

	if claims.Subject != "user-123" {
		t.Errorf("subject = %q, want user-123", claims.Subject)
	}
	if claims.ProductID != "product-456" {
		t.Errorf("product_id = %q, want product-456", claims.ProductID)
	}
	if claims.Role != "admin" {
		t.Errorf("role = %q, want admin", claims.Role)
	}
}

func TestVerifyExpiredToken(t *testing.T) {
	m := newTestManager(t)

	token, err := m.Sign("user-123", "product-456", "member", -time.Minute)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}

	_, err = m.Verify(token)
	if err == nil {
		t.Fatal("expected error for expired token, got nil")
	}
}

func TestJWKS(t *testing.T) {
	m := newTestManager(t)
	jwks := m.JWKS()

	if len(jwks.Keys) != 1 {
		t.Fatalf("expected 1 key, got %d", len(jwks.Keys))
	}
	k := jwks.Keys[0]
	if k.Kty != "RSA" || k.Alg != "RS256" || k.Use != "sig" {
		t.Errorf("unexpected JWK fields: %+v", k)
	}
	if k.N == "" || k.E == "" {
		t.Error("JWK N or E is empty")
	}
}
