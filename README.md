# FV Generator

### Summary

FV Generator is a simple CLI application that generates invoices in PDF format.

### How to use

0. Install [Go compiler](https://go.dev/dl/)

1. Create `recipients.json` file - use `recipients-sample.json` as a reference. Content of this file will be used by HTML parser so you can use special symbols like non-breaking white space (`&nbsp;`).

2. (Optional) Create `.fvrc` file to prefill the default entry in the invoice and source it in your shell configuration.

```sh
export FV_DEFAULT_PROMPT="Description of the&nbsp;entry" # Use can use special symbols there as well.
export FV_DEFAULT_PRICE="1"
export FV_DEFAULT_UNIT="hour"
export FV_DEFAULT_AMOUNT="1"
export FV_NOW="2006-01-02" # time.DateOnly
export FV_USE_DEFAULT="1" # Disable/Enable default entry.
```

3. Build the binary and execute it.

```
go build main.go
chmod +x fv-generator
./fv-generator
```

4. Follow terminal prompts for further instructions.

Generated PDF file will be saved to `output` directory.
