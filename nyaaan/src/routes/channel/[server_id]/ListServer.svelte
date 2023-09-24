<script lang="ts">
	import type { Server } from './type';
	import { showCreateServerModal } from './store';
	import { RiAddLine } from 'svelte-remixicon';
	import { goto } from '$app/navigation';

	export let selectedServer: Server;
	export let servers: Array<Server>;
</script>

<div>
	{#each servers as server}
		{#if server.id == selectedServer.id}
			<div class="mt-2 avatar h-12 w-12 ring rounded-full ring-info">
				<button class="rounded-full">
					<img src={server.avatar} alt="nekonako" class="rounded-full" />
				</button>
			</div>
		{:else}
			<div class="mt-2 avatar rounded-full">
				<button
					class="rounded-full"
					on:click={() => {
						selectedServer = server;
						goto(`/channel/${server.id}`);
					}}
				>
					<img src={server.avatar} alt="nekonako" class="rounded-full" />
				</button>
			</div>
		{/if}
	{/each}
</div>
<div class="mb-2">
	<button
		class="rounded-full bg-base-100 justify-center p-2"
		on:click={() => showCreateServerModal.set(true)}
	>
		<RiAddLine size="2em" class="text-success text-center" />
	</button>
</div>
