# Security Policy

PR1MAL-VAULT is designed to provide secure local secret storage using
modern cryptographic primitives.

## Cryptography

PR1MAL-VAULT uses:

Argon2id — password key derivation  
AES-256-GCM — authenticated encryption

Argon2 parameters:

Memory: 64 MB  
Iterations: 3  
Parallelism: 2  

AES-GCM ensures that any modification to encrypted vault data
causes decryption failure.

## Threat Model

PR1MAL-VAULT protects secrets **at rest on a local machine**.

It protects against:

- accidental secret exposure
- disk inspection
- vault file tampering

It does NOT protect against:

- malware on the host system
- keyloggers
- compromised operating systems

If the machine is compromised, the vault cannot guarantee secrecy.

## Reporting Vulnerabilities

If you discover a security vulnerability, please open a GitHub issue
or contact the maintainer.

Responsible disclosure is appreciated.
