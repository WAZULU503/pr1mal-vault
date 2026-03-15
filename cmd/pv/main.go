package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/wazulu503/pr1mal-vault/internal/crypto"
	"github.com/wazulu503/pr1mal-vault/internal/store"
	"golang.org/x/term"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	switch os.Args[1] {
	case "init":
		initializeVault()
	case "set":
		setSecret()
	case "get":
		getSecret()
	case "ls":
		listSecrets()
	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println("PR1MAL-VAULT")
	fmt.Println("Usage:")
	fmt.Println("  pv init")
	fmt.Println("  pv set <key>")
	fmt.Println("  pv get <key>")
	fmt.Println("  pv ls")
}

func getPassword(prompt string) string {
	fmt.Print(prompt)
	pass, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return ""
	}
	return string(pass)
}

func initializeVault() {

	path := store.GetVaultPath()

	if _, err := os.Stat(path); err == nil {
		fmt.Println("Vault already exists. Aborting.")
		return
	}

	password := getPassword("Create Master Password (min 12 chars): ")

	if len(password) < 12 {
		fmt.Println("Password too short.")
		return
	}

	confirm := getPassword("Confirm Master Password: ")

	if password != confirm {
		fmt.Println("Password mismatch.")
		return
	}

	kdf, err := crypto.GenerateKDFParams()
	if err != nil {
		fmt.Println("KDF error:", err)
		return
	}

	key := crypto.DeriveKey(password, kdf)

	secrets := make(map[string]string)
	plaintext, _ := json.Marshal(secrets)

	ciphertext, err := crypto.Encrypt(plaintext, key)
	if err != nil {
		fmt.Println("Encryption error:", err)
		return
	}

	vault := &store.VaultFile{
		Version: 1,
		Cipher:  "AES-256-GCM",
		KDF:     *kdf,
		Payload: ciphertext,
	}

	if err := store.Save(vault); err != nil {
		fmt.Println("Save error:", err)
		return
	}

	fmt.Println("Vault initialized:", path)
}

func unlockVault(password string) (*store.VaultFile, map[string]string, []byte, error) {

	vault, err := store.Load()
	if err != nil {
		return nil, nil, nil, err
	}

	key := crypto.DeriveKey(password, &vault.KDF)

	plaintext, err := crypto.Decrypt(vault.Payload, key)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("invalid password or corrupted vault")
	}

	var secrets map[string]string

	if err := json.Unmarshal(plaintext, &secrets); err != nil {
		return nil, nil, nil, err
	}

	return vault, secrets, key, nil
}

func setSecret() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: pv set <key>")
		return
	}

	keyName := os.Args[2]

	password := getPassword("Master Password: ")

	vault, secrets, key, err := unlockVault(password)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Enter secret value for [%s]: ", keyName)
	value, _ := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()

	secrets[keyName] = string(value)

	plaintext, _ := json.Marshal(secrets)

	ciphertext, err := crypto.Encrypt(plaintext, key)
	if err != nil {
		fmt.Println(err)
		return
	}

	vault.Payload = ciphertext

	if err := store.Save(vault); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Secret stored.")
}

func getSecret() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: pv get <key>")
		return
	}

	keyName := os.Args[2]

	password := getPassword("Master Password: ")

	_, secrets, _, err := unlockVault(password)
	if err != nil {
		fmt.Println(err)
		return
	}

	value, ok := secrets[keyName]
	if !ok {
		fmt.Println("Secret not found")
		return
	}

	fmt.Println(value)
}

func listSecrets() {

	password := getPassword("Master Password: ")

	_, secrets, _, err := unlockVault(password)
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(secrets) == 0 {
		fmt.Println("Vault empty")
		return
	}

	fmt.Println("Stored secrets:")

	for k := range secrets {
		fmt.Println(" -", k)
	}
}
