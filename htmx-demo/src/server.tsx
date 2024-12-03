import {type Context, Hono} from 'hono'
import {serveStatic} from 'hono/bun'

const app = new Hono()

// Serve static files from public dir
app.use('/', serveStatic({root: './public'}))
app.get('/version', (c: Context) => {
  return c.text(Bun.version)
})

export default app