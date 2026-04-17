# vaultpipe

A CLI tool to inject secrets from HashiCorp Vault into process environments without shell exposure.

---

## Installation

```bash
go install github.com/yourusername/vaultpipe@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/vaultpipe/releases).

---

## Usage

`vaultpipe` runs a command with secrets from Vault injected as environment variables. Secrets are never written to shell history or exposed via `ps`.

```bash
vaultpipe run --path secret/data/myapp -- ./myapp serve
```

### Options

| Flag | Description |
|------|-------------|
| `--path` | Vault secret path to read from |
| `--addr` | Vault server address (default: `$VAULT_ADDR`) |
| `--token` | Vault token (default: `$VAULT_TOKEN`) |
| `--prefix` | Optional env var prefix for injected secrets |

### Example

```bash
export VAULT_ADDR=https://vault.example.com
export VAULT_TOKEN=s.mytoken

vaultpipe run --path secret/data/db --prefix DB_ -- ./myapp migrate
```

Secrets stored at `secret/data/db` (e.g., `username`, `password`) will be available to `myapp` as `DB_USERNAME` and `DB_PASSWORD`.

---

## How It Works

`vaultpipe` authenticates with Vault, fetches the specified secrets, and injects them directly into the child process environment using `exec`. No secrets touch your shell environment or command history.

---

## Contributing

Pull requests and issues are welcome. Please open an issue before submitting large changes.

---

## License

MIT © [yourusername](https://github.com/yourusername)