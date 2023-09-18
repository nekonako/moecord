import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
// import http from ''

export default defineConfig({
	plugins: [sveltekit()],
	server: {
		proxy: {
			'/api': {
				target: 'http://localhost:4000/v1',
				rewrite: (path) => path.replace(/^\/api/, ''),
				changeOrigin: true,
				configure: (proxy, options) => {
					proxy.on('proxyReq', (proxyReq, req, res) => {
						if (req.headers.cookie) {
							const cookies: Array<string> = req.headers.cookie.split('; ');
							const mapCookies: Map<string, string> = new Map();
							cookies.forEach((val) => {
								const cookie = val.split('=');
								if (cookie.length >= 2) {
									mapCookies.set(cookie[0], cookie[1]);
								}
							});
							proxyReq.setHeader('Authorization', 'Bearer ' + mapCookies.get('access_token'));
						}

					});
				}
			}
		}
	}
});
