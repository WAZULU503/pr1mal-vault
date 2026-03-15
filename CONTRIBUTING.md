# Contributing

Contributions are welcome.

## Development

Build the project:

go build -o pv ./cmd/pv

Run the tool:

./pv init

## Code Style

PR1MAL-VAULT follows standard Go conventions.

Please ensure:

- code builds successfully
- no plaintext secrets are written to disk
- cryptographic primitives are not modified without review

## Pull Requests

Submit pull requests with clear commit messages describing the change.
