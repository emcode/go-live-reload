serve:
  watchexec \
    --watch ./templates \
    --watch ./website.go \
    --watch ./live_reload.go \
    --debounce 60ms \
    --restart \
    -- go run .

build:
  go build .