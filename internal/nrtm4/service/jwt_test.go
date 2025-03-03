package service

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"io"
	"log"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/petchells/nrtm4tools/internal/nrtm4/persist"
	"github.com/petchells/nrtm4tools/internal/nrtm4/testresources"
)

func TestJWT(t *testing.T) {
	//ntfyURL := "https://nrtm.db.ripe.net/nrtmv4/RIPE/update-notification-file.jose"
	// keyURL:="https://ftp.ripe.net/ripe/dbase/nrtmv4/nrtmv4_public_key.txt"
	keyTxt := `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEOkzpjobirEcqoR6zLXnPkm4cCTEY
Xi2rLlCSXc5EZ3L3PycAdDmWQtGHD8GF++RqWgrdKv+9l+InalmiCGkpRQ==
-----END PUBLIC KEY-----`
	file := testresources.OpenFile(t, "update-notification-file.jose")
	if file == nil {
		t.Fatal("Error opening file")
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		t.Fatal("Error reading file")
	}

	block, _ := pem.Decode([]byte(keyTxt))
	if block == nil || block.Type != "PUBLIC KEY" {
		log.Fatal("failed to decode PEM block containing public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}
	tokenString := string(bytes)
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return pub, nil
	})
	if err != nil {
		t.Error("Parse failed", err)
	}
	t.Log("token.Header", token.Header)
	// do something with decoded claims
	cljson, err := json.Marshal(claims)
	if err != nil {
		t.Error("Marshal failed", err)
	}
	notification := new(persist.NotificationJSON)
	err = json.Unmarshal(cljson, notification)
	if err != nil {
		t.Error("Unmarshal failed", err)
	}
	t.Log("unf", *notification)
}
