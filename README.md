![Security Audit](https://github.com/WAZULU503/pr1mal-vault/actions/workflows/security.yml/badge.svg)

# PR1MAL-VAULT

Deterministic local secret storage for developers.

PR1MAL-VAULT (`pv`) is a minimal command-line vault written in Go that stores API keys, tokens, and credentials securely on your machine using modern cryptography.

Secrets are encrypted locally and never stored in plaintext.

Everything stays on your machine.

No cloud. No telemetry.

---

## Features

- Argon2id password hardening
- AES-256-GCM authenticated encryption
- Versioned vault file format
- Atomic vault writes (temp file → rename)
- Hidden password input (prevents shell history leaks)
- Vault file permission isolation (`0600`)

---

## Installation

Clone the repository:

```bash
git clone https://github.com/WAZULU503/pr1mal-vault.git
cd pr1mal-vault
```

Build the binary:

```bash
go build -o pv ./cmd/pv
```

---

## Usage

Initialize the vault:

```bash
pv init
```

Store a secret:

```bash
pv set API_KEY
```

Retrieve a secret:

```bash
pv get API_KEY
```

List stored secrets:

```bash
pv ls
```

Delete a secret:

```bash
pv delete API_KEY
```

---

## Example Workflow

```bash
pv init
pv set GITHUB_TOKEN
pv ls
pv get GITHUB_TOKEN
```

---

## Vault Location

The encrypted vault file is stored at:

```
~/.pr1mal_vault
```

Example vault structure:

```json
{
  "v": 1,
  "c": "AES-256-GCM",
  "kdf": {
    "mem": 65536,
    "iter": 3,
    "par": 2,
    "salt": "base64..."
  },
  "p": "encrypted payload"
}
```

Secrets themselves never appear in plaintext.

---

## Security Design

PR1MAL-VAULT uses modern cryptographic primitives designed to protect secrets stored on a local machine.

### Key Derivation

Passwords are hardened using **Argon2id**.

Parameters:

- Memory: 64 MB
- Iterations: 3
- Parallelism: 2

This makes brute-force attacks significantly more expensive.

### Encryption

Secrets are encrypted using **AES-256-GCM** authenticated encryption.

AES-GCM provides:

- confidentiality
- integrity
- tamper detection

If any bit of the vault file changes, decryption fails.

### Secure Input

Secrets and passwords are entered using masked terminal input:

```
term.ReadPassword()
```

This prevents secrets from appearing in:

- `.bash_history`
- `.zsh_history`
- process arguments
- terminal logs

### Atomic Writes

Vault updates are written using a safe update pattern:

```
write temp file → rename
```

This prevents corrupted vault files during crashes or interruptions.

### File Permissions

The vault file is created with permission mode:

```
0600
```

Only the current user can read or modify the vault.

---

## Project Structure

```
cmd/pv
internal/crypto
internal/store
```

| Module | Purpose |
|------|------|
| crypto | encryption and key derivation |
| store | vault persistence |
| cmd/pv | CLI interface |

---

## Threat Model

PR1MAL-VAULT protects secrets **at rest on a local machine**.

It protects against:

- accidental secret exposure
- disk inspection
- vault file tampering

It does **not** protect against:

- malware on the host system
- keyloggers
- compromised operating systems

If the machine is compromised, the vault cannot guarantee secrecy.

---

## Verification

You can verify the encryption pipeline locally.

Build the binary:

```bash
go build -o pv ./cmd/pv
```

Initialize a vault:

```bash
./pv init
```

Store a test secret:

```bash
./pv set TEST_KEY
```

Retrieve it:

```bash
./pv get TEST_KEY
```

The stored vault file should contain only encrypted data:

```bash
cat ~/.pr1mal_vault
```

The output should contain encrypted payload data rather than plaintext secrets.

---

## Author

Wazulu the Ill Dravidian

GitHub  
https://github.com/WAZULU503

---

## License

MIT License
