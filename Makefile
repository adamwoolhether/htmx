# curl -fsSL https://bun.sh/install | bash
# bunx create-hono
# bun install
# bun run dev

templ:
	templ generate

go:
	./run.sh

bun:
	cd ts/htmx-demo && bun dev