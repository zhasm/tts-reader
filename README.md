# tts-reader

## install





```sh
tools=(
  github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  golang.org/x/tools/cmd/gopls@latest
  golang.org/x/tools/cmd/goimports@latest
  golang.org/x/tools/cmd/gorename@latest
  golang.org/x/tools/cmd/goimports@latest
)

for i in ${tools[@]}; do
	go install "$i"
done
```
```
```
