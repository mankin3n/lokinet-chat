package user

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"os"
)

type Profile struct {
	Username   string
	PublicKey  string
	PrivateKey string
}

// GenerateRSAKeyPair generates a new RSA key pair.
func GenerateRSAKeyPair(bits int) (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", err
	}

	// Export private key
	privKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privKeyBytes,
	})

	// Export public key
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", err
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubKeyBytes,
	})

	return string(publicKeyPEM), string(privateKeyPEM), nil
}

func SaveProfile(profile Profile) error {
	file, err := os.Create("profile.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(profile)
}

func LoadProfile() (Profile, error) {
	file, err := os.Open("profile.json")
	if err != nil {
		return Profile{}, err
	}
	defer file.Close()

	var profile Profile
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&profile)
	return profile, err
}

func CreateProfile(username string) (Profile, error) {
	publicKey, privateKey, err := GenerateRSAKeyPair(2048)
	if err != nil {
		return Profile{}, err
	}

	profile := Profile{
		Username:   username,
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}

	return profile, nil
}
