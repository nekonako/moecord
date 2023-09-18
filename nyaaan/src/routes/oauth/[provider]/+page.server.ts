import { redirect } from '@sveltejs/kit';

export async function load({ params, url, cookies, fetch }) {
	const provider = params.provider;
	const payload = {
		authorization_code: url.searchParams.get('code'),
		state: url.searchParams.get('state'),
		provider: provider
	};

	const response = await fetch('/api/login/oauth/callback/' + provider, {
		body: JSON.stringify(payload),
		method: 'POST'
	});

	const result = await response.json();
	if (result.code != 200) {
		throw redirect(307, '/oauth');
	}

	cookies.set('access_token', result.data.access_token, {
		httpOnly: false,
		secure: false,
		path: '/'
	});

	cookies.set('refresh_token', result.data.refresh_token, {
		httpOnly: false,
		secure: false,
		path: '/'
	});

	throw redirect(307, '/channel');
}
