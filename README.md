# subseek

Early Go CLI for listing / seeking Azure subscriptions (reads `creds.config` for service principal fields and can call `az`). Superseded in practice by **`zxsubs`**, which uses an interactive terminal picker on top of an existing `az login` session.

## Note

A compiled `subseek` binary is committed in this tree (large; candidate to remove from git and ignore). Prefer building from source:

```bash
go build -o subseek .
```

Do not commit real `creds.config` values.
