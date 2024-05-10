server.w:
	@trap 'kill $(jobs -p) 2>/dev/null' EXIT
	devwatcher &
	@sleep 1 # make sure devwatcher is started
	watchexec --watch . --restart -- \
            'npx tailwindcss -i ./public/raw.css -o ./public/styles.css -m; go run ./cmd/server'
