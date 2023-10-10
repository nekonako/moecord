import type { ApiResponse, Fetch, Profile } from "./type";

export async function getUserProfile(fetch: Fetch) {
	try {
		let response = await fetch("/api/profile");
		let result : ApiResponse<Profile> = await response.json();
		return result
	} catch (error) {
		throw error
	}
}

export async function updateUserProfile(fetch: Fetch, payload: FormData) {
	try {
		let response = await fetch("/api/profile", {
			body: payload,
			method: 'PUT'
		});
		let result : ApiResponse<any> = await response.json();
		return result
	} catch (error) {
		throw error
	}
}

