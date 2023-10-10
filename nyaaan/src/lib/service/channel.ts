import type { ApiResponse, Channel, CreateChannelCategoryRequest, CreateChannelRequest, Fetch } from "./type";


export async function getListChannel(fetch: Fetch, serverID: string) {
	try {
		const response = await fetch("/api/channels/" + serverID);
		const result: ApiResponse<Array<Channel>> = await response.json();
		return result
	} catch (error) {
		throw error
	}
}

export async function createChannel(fetch: Fetch, channel: CreateChannelRequest) {
	try {
		const response = await fetch('/api/channels', {
			body: JSON.stringify(channel),
			method: 'POST'
		});
		const result = await response.json();
		return result
	} catch (error) {
		throw error
	}
}

export async function createChannelCategory(fetch: Fetch, category: CreateChannelCategoryRequest) {
	try {
		const response = await fetch('/api/channels/categories', {
			body: JSON.stringify(category),
			method: 'POST'
		});
		const result = await response.json();
		return result
	} catch (error) {
		throw error
	}
}
