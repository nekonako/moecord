<script lang="ts">
	import { RiAddLine, RiVolumeUpFill, RiChat1Fill, RiSendPlaneFill } from 'svelte-remixicon';
	import { browser } from '$app/environment';
	export let data;
	let ws: WebSocket | undefined;
	let message: string;

	if (browser) {
		ws = new WebSocket('ws://localhost:4001');
		console.log(ws);
		handleMessage(ws);
	}

	type Message = {
		channel_id: string;
		content: string;
	};

	async function sendMessage() {
		const payload: Message = {
			channel_id: data.channels[0].id,
			content: message
		};
		const response = await fetch('/api/messages', {
			body: JSON.stringify(payload),
			method: 'POST'
		});
		const result = await response.json();
		console.log(result);
	}

	function handleMessage(ws: WebSocket) {
		ws.onmessage = (m) => {
			console.log(m.data);
		};
	}
</script>

<div
	class="flex flex-row min-w-screen min-h-screen bg-stone-950 text-white font-sans text-md antialiased font-light"
>
	<div class="w-20 flex flex-col justify-between items-center border-r border-stone-800">
		<div>
			{#each data.servers as server}
				<button class="mx-3 my-3">
					<img
						src="https://avatars.githubusercontent.com/u/46141275?v=4"
						alt="nekonako"
						class="relative inline-block rounded-full"
					/>
				</button>
			{/each}
		</div>
		<div class="mb-2">
			<button class="w-full p-1 rounded-lg bg-stone-800/40">
				<RiAddLine class="text-orange-500" size="2em" />
			</button>
		</div>
	</div>
	<div class="w-1/5 flex flex-col">
		{#each data.channels as channel}
			<div class="py-1 px-4 text-sm flex flex-row items-center lowercase">
				{#if channel.channel_type === 'text'}
					<RiChat1Fill size="1.2em" />
				{:else}
					<RiVolumeUpFill size="1.2em" />
				{/if}
				<span class="ml-2">{channel.name}</span>
			</div>
		{/each}
	</div>
	<div class="border border-stone-800 w-3/5">
		<div class="flex flex-col justify-between min-h-screen">
			<div class="flex flex-col">
				{#each data.messages as message}
					<div class="m-4 text-white">
						{message.content}
					</div>
				{/each}
			</div>
			<div class="flex flex-row m-4">
				<input bind:value={message} class="w-full bg-stone-800/40 h-12 rounded-l-lg" />
				<button
					class="text-orange-600 bg-stone-800/40 rounded-r-lg px-4 border-l border-stone-800"
					on:click={() => sendMessage()}
				>
					<RiSendPlaneFill />
				</button>
			</div>
		</div>
	</div>
	<div class="w-1/5">
		<div class="flex flex-col justify-between min-h-screen">
			<div>members</div>
		</div>
	</div>
</div>
