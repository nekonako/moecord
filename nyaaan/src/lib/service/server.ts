import type { ApiResponse, CreateServeerRequest, Fetch, Server, ServerMember } from "./type";


export async function getListServer(fetch: Fetch) {
	try {
		const response = await fetch("/api/servers");
		const result: ApiResponse<Array<Server>> = await response.json();
		return result
	} catch (error) {
		throw error
	}
}

export async function getServerMember(fetch: Fetch, serverID: string) {
	try {
		const response = await fetch("/api/servers/" + serverID + "/member",);
		const result: ApiResponse<Array<ServerMember>> = await response.json();
		return result
	} catch (error) {
		throw error
	}
}

export async function createServer(fetch: Fetch, server: CreateServeerRequest) {
	try {
		const response = await fetch('/api/servers', {
			body: JSON.stringify(server),
			method: 'POST'
		});
		const result = await response.json();
		return result
	} catch (error) {
		throw error
	}
}

export async function updateServer(fetch: Fetch, payload: FormData) {
	try {

		const response = await fetch(`/api/servers`, {
			method: 'PUT',
			body: payload
		});
		const result: ApiResponse<Server> = await response.json();
		return result
	} catch (error) {
		throw error
	}

}
