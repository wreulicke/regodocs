# regodocs

Document Generator from Rego. 
This generates markdown files from [annotations](https://www.openpolicyagent.org/docs/latest/policy-language/#metadata) in rego files.

This tool aims to help documentation for conftest.

## Install

```bash
# linux/macOS
curl -o regodocs -sL "https://github.com/wreulicke/regodocs/releases/download/v0.0.2/regodocs_0.0.2_$(uname -s | tr "[:upper:]" "[:lower:]")_$(uname -m)"
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

## License

MIT License
