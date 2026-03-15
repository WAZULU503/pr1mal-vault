package store

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/wazulu503/pr1mal-vault/internal/crypto"
)

type VaultFile struct {
	Version int              `json:"v"`
	Cipher  string           `json:"c"`
	KDF     crypto.KDFParams `json:"kdf"`
	Payload []byte           `json:"p"`
}

func GetVaultPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".pr1mal_vault")
}

func Save(v *VaultFile) error {
	path := GetVaultPath()
	tmp := path + ".tmp"

	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(tmp, data, 0600); err != nil {
		return err
	}

	return os.Rename(tmp, path)
}

func Load() (*VaultFile, error) {
	path := GetVaultPath()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var v VaultFile
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, err
	}

	return &v, nil
}
