<script lang="ts">
	import { RiCloseFill } from 'svelte-remixicon';
	import { showCreateServerModal } from './store';
	import { createServer, getListServer } from '$lib/service/server';
	import { createEventDispatcher } from 'svelte';

	let serverName: string;
	const dispatch = createEventDispatcher();

	async function handleCreateServer() {
		const response = await createServer(fetch, { name: serverName });
		if (response.code == 200) {
			showCreateServerModal.set(false);
			const response = await getListServer(fetch);
			dispatch('getListServer', response.data);
		}
	}
</script>

<div class="min-h-screen min-w-screen flex justify-center items-center bg-neutral">
	<div class="modal-box">
		<div class="flex w-ful align-right justify-end">
			<button on:click={() => showCreateServerModal.set(false)}>
				<RiCloseFill size={'1.5em'} /></button
			>
		</div>
		<h3 class="font-bold text-lg mb-4">SERVER NAME</h3>
		<input
			bind:value={serverName}
			type="text"
			placeholder="Moe..."
			class="input input-lg bg-neutral opacity-75 focus:outline-none w-full input-bordered text-white"
		/>
		<div class="modal-action">
			<form method="dialog">
				<button class="btn btn-info" on:click={() => handleCreateServer()}>Create</button>
			</form>
		</div>
	</div>
</div>
