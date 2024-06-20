# regodocs

Document Generator from Rego. 
This generates markdown from rego files.

## Install

```bash
# linux/macOS
curl -o regodocs -sL "https://github.com/wreulicke/regodocs/releases/download/v0.0.1/regodocs_0.0.1_$(uname -s | tr "[:upper:]" "[:lower:]")_$(uname -m)"
chmod +x regodocs
mv regodocs /usr/local/bin

# windows
TBD
```

## Usage

```
$ regodocs generate
Generate documentation from Rego policy files

Usage:
  regodocs generate POLICY_PATH... [flags]

Flags:
  -h, --help              help for generate
  -o, --output string     output path for generated documentation
  -p, --pattern strings   pattern to filter files (default [deny.*,violation.*,warn.*])
```

## Example

See [testadata](./testdata/) and
also See [generated document](./.snapshot/TestGenerator)

## TODO

- Add description for rules

## License

MIT License
