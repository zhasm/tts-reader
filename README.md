# tts-reader

## install


```sh
tools=(
  github.com/cweill/gotests/gotests@v1.6.0
  github.com/go-delve/delve/cmd/dlv@latest
  github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  github.com/haya14busa/goplay/cmd/goplay@v1.0.0
  github.com/josharian/impl@v1.4.0
  github.com/rakyll/gotest@latest
  golang.org/x/tools/cmd/goimports@latest
  golang.org/x/tools/cmd/gopls@latest
  golang.org/x/tools/cmd/gorename@latest
  honnef.co/go/tools/cmd/staticcheck@latest
)

for i in ${tools[@]}; do
	GO111MODULE=on go install -x -v "$i"
done
```
