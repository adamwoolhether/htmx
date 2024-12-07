# curl -fsSL https://bun.sh/install | bash
# bunx create-hono
# bun install
# bun run dev
# npm init @eslint/config
# bun run lint
# bun add -d prettier
# bun run format

templ:
	templ generate

go:
	./run.sh

bun:
	cd ts/htmx-demo && bun dev