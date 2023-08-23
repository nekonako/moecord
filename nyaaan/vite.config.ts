import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [sveltekit()],
	server: {
		proxy: {
			'/api': {
				target: 'http://localhost:4000/v1',
				rewrite: (path) => path.replace(/^\/api/, ''),
				changeOrigin: true,
				configure: (proxy, options) => {
					// proxy will be an instance of 'http-proxy'
				}
			}
		}
	}
});
