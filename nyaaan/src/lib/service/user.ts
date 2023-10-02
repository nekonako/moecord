import type { ApiResponse, Fetch, Profile } from "./type";

export async function getuserProfile(fetch: Fetch) {
	try {
		let response = await fetch("/api/profile");
		let result : ApiResponse<Profile> = await response.json();
		return result
	} catch (error) {
		throw error
	}
}
