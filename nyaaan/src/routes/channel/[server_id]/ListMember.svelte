<script lang="ts">
	import { MapServerMember } from './store';
	import type { ServerMember, UserConnectionState, WebsocketMessage } from '$lib/service/type';
	import { getColor } from '$lib/util';

	export let members: Array<ServerMember>;
	export let wsMessageHandler: Set<(e: MessageEvent) => void> = new Set();
	let filteredMembers = [...members];
	let keyword = '';

	function search() {
		if (keyword != '') {
			const result = filteredMembers.filter((value) => value.username.includes(keyword));
			members = [...result];
			return;
		}
		if (!keyword) {
			members = filteredMembers;
		}
	}

	wsMessageHandler.add((e) => {
		let message: WebsocketMessage<any> = JSON.parse(e.data);
		if (message.event_id == 'USER_DISCONNECTED' || message.event_id == 'NEW_CONNECTION') {
			console.log('connection state changes');
			let data = message.data as UserConnectionState;
			members.forEach((value) => {
				if (value.user_id == data.user_id) {
					console.log(value.user_id, data.user_id, value);
					value.online = data.status == 'online';
				}
			});
			members = members;
		}
	});
</script>

<input
	type="text"
	bind:value={keyword}
	on:input={search}
	class="bg-base-content/20 w-full px-2 mb-8 py-1 rounded-md focus:outline-none"
	placeholder="search..."
/>
{#each members as member}
	<div class="flex flex-row items-center mb-4">
		{#if member.online}
			<div class="avatar online mr-4">
				<div class="w-8 rounded-full">
					<img src={member.avatar} alt="profile" />
				</div>
			</div>
		{:else}
			<div class="avatar mr-4">
				<div class="w-8 rounded-full">
					<img src={member.avatar} alt="profile" />
				</div>
			</div>
		{/if}

		<div class="fle flex-col">
			<span class={getColor(member.username.toUpperCase().charCodeAt(0))}>{member.username}</span>
			{#if $MapServerMember[member.user_id]}
				<div class="-mt-2">
					<span class="text-xs leading-none text-info">typing...</span>
				</div>
			{/if}
		</div>
	</div>
{/each}
