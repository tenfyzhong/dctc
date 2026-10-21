# dctc
dctc is a tool to generate a docker-compose file from containers. 

# Usage
![](https://ws1.sinaimg.cn/large/69472223gy1g1pd2sql4cj20og0dpjt7.jpg)

# Install
## Get from releases
[](https://github.com/tenfyzhong/dctc/releases)
[![GitHub tag](https://img.shields.io/github/release/tenfyzhong/dctc.svg)](https://github.com/tenfyzhong/dctc/releases)

## Building
Build with Go 1.26.10 (or Go 1.27.3 or newer on that release line), which
includes the `crypto/x509` fix for [GO-2026-6809](https://pkg.go.dev/vuln/GO-2026-6809).
The module keeps its Go 1.14 language version; the toolchain directive selects
the patched compiler when building this checkout with modern Go tooling.

```sh
git clone https://github.com/tenfyzhong/dctc.git
cd dctc
GOTOOLCHAIN=go1.26.10 go install .
```

CI pins Go 1.26.10 and verifies the compiler used for the built executable.
Previously downloaded binaries must be rebuilt or replaced; this source change
does not patch existing releases.
