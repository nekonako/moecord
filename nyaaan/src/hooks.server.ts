import type { MaybePromise, RequestEvent, ResolveOptions } from '@sveltejs/kit';

export async function handle({
	event,
	resolve
}: {
	event: RequestEvent;
	resolve(event: RequestEvent, opts?: ResolveOptions): MaybePromise<Response>;
}) {
	let cookies = event.cookies;
	let access_token = cookies.get('access_token');
	const path = event.url.pathname;
	if (!access_token && !path.startsWith('/oauth')) {
		return Response.redirect(new URL('/oauth', event.url), 307);
	}
	event.request.headers.set('Authorization', 'Bearer ' + access_token);
	const response = await resolve(event);
	return response;
}
