# SingXClient

A cross-platform graphical/CLI client built on the sing-box core and designed to connect seamlessly with the XBoard 2025 panel.  It automatically fetches subscription data from your XBoard server, generates a full **sing-box** configuration (supporting the latest protocols such as **Hysteria2**, **AnyTLS**, **VLESS+Reality**, …) and launches the proxy engine.

The repository is shipped with a ready-to-use GitHub Actions workflow & GoReleaser configuration so that, once you push a tag, installers/binaries for **Windows**, **macOS** (universal binary), **Linux**, and **Android** will be built, packaged and uploaded to the release page automatically.

---

## Features

* 🚀 Zero-config runtime: just specify the XBoard subscription URL or token.
* 🔌 Supports all protocols currently implemented by sing-box (`hysteria2`, `anytls`, `vless+reality`, …).
* 🔄 Auto-refresh subscription & hot-reload sing-box without downtime.
* 📦 One-command release workflow powered by _GoReleaser_.
* 🪵 Structured logging (JSON / console) & optional tray icon on desktop (planned).

---

## Quick start (development)

```bash
# Clone your fork, then:
make generate # install tools if necessary
make build    # local build for your OS/arch
./bin/singxclient --subscription "https://example.com/api/v1/client/subscribe?token=..."
```

---

## Building / releasing

Releases are handled automatically by GitHub Actions.  Tag a commit – e.g.

```bash
git tag -a v0.1.0 -m "First public release"
git push origin v0.1.0
```

A workflow will spin up, execute GoReleaser inside a matrix that targets all major operating systems and upload the artefacts to the _v0.1.0_ release page.

---

## Repository layout

```
├── cmd/                # CLI entry points (cobra)
├── internal/
│   ├── xboard/         # Subscription & profile parsing logic
│   └── config/         # sing-box config generator
├── .github/workflows/  # CI/CD pipelines
└── .goreleaser.yaml    # Build & packaging definition
```

---

## License

This project is released under the GPL-3.0 license, following the upstream sing-box licensing terms.
