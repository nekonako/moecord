import type { MaybePromise, RequestEvent, ResolveOptions } from '@sveltejs/kit';

export async function handle({
	event,
	resolve
}: {
	event: RequestEvent;
	resolve(event: RequestEvent, opts?: ResolveOptions): MaybePromise<Response>;
}) {
	let cookies = event.cookies;
	console.log('cok 123');
	let access_token = cookies.get('access_token');
	const path = event.url.pathname;
	if (!access_token && !path.startsWith('/oauth')) {
		return await Response.redirect('/oauth', 307);
	}
	event.request.headers.set('Authorization', 'Bearer ' + access_token);
	const response = await resolve(event);
	return response;
}
