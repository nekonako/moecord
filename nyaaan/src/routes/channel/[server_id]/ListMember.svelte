<script lang="ts">
	import { MapServerMember } from './store';
	import { getColor, type Servermember, type WebsocketMessage } from './type';

	export let members: Array<Servermember>;

	let filteredMembers = [...members];
	let keyword = '';

	function search() {
		if (keyword != '') {
			members = filteredMembers.filter((value) => value.username.startsWith(keyword));
			return;
		}
		if (keyword == '') {
			members = filteredMembers;
		}
	}
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
		<div class="avatar online mr-4">
			<div class="w-8 ring ring-success rounded-full">
				<img src={member.avatar} alt="profile" />
			</div>
		</div>
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
