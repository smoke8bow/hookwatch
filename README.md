# hookwatch

A local webhook relay server with request inspection and replay for testing webhook integrations.

---

## Installation

```bash
go install github.com/yourusername/hookwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/hookwatch.git && cd hookwatch && go build -o hookwatch .
```

---

## Usage

Start the relay server on a local port:

```bash
hookwatch --port 8080
```

hookwatch will listen for incoming webhook requests and display a full inspection of each request in your terminal, including headers, payload, and timestamp.

**Replay a captured request:**

```bash
hookwatch replay --id <request-id> --target http://localhost:3000/webhook
```

**Forward incoming requests to your local service:**

```bash
hookwatch --port 8080 --forward http://localhost:3000/webhook
```

All captured requests are stored in memory during the session and can be listed with:

```bash
hookwatch list
```

---

## Why hookwatch?

- No external dependencies or cloud tunnels required
- Instant request inspection directly in your terminal
- Replay any captured request to debug handler logic
- Lightweight and easy to integrate into local dev workflows

---

## License

MIT © [yourusername](https://github.com/yourusername)