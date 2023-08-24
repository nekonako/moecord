import { redirect } from '@sveltejs/kit';

type ApiResponse<T> = {
	code: number;
	message: string;
	data: Array<T>;
};

type Server = {
	id: string;
	name: string;
};

type Channel = {
	id: string;
	name: string;
};

export async function load({ fetch }) {
	const responseServer = await fetch('/api/servers');
	const servers: ApiResponse<Server> = await responseServer.json();
	if (servers.code == 401) {
		throw redirect(307, '/oauth');
	}
	const firstServerID = servers.data[0].id;
	const responseChannel = await fetch('/api/channels/' + firstServerID);
	const channels: ApiResponse<Channel> = await responseChannel.json();

	return {
		servers: servers.data,
		channels: channels.data
	};
}
