<script lang="ts">
	import { RiAddLine, RiCloseFill, RiSendPlaneFill } from 'svelte-remixicon';
	import { browser } from '$app/environment';
	export let data;
	let ws: WebSocket | undefined;
	let message: string;
	let messages = data.messages;
	let showModal = false;
	let serverName = '';
	let selectedServer = data.servers[0].name;

	if (browser) {
		ws = new WebSocket('ws://localhost:4000/ws');
		handleMessage(ws);
	}

	type Message = {
		channel_id: string;
		content: string;
	};

	type CreateServer = {
		name: string;
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
	}

	async function createServer() {
		const payload: CreateServer = {
			name: serverName
		};
		const response = await fetch('/api/servers', {
			body: JSON.stringify(payload),
			method: 'POST'
		});
		const result = await response.json();
		console.log(result);
		if (result.code == 200) {
			showModal = !showModal;
		}
	}

	function handleMessage(ws: WebSocket) {
		ws.onmessage = (m) => {
			messages = [...messages, JSON.parse(m.data)];
		};
	}
</script>

{#if showModal}
	<div class="min-h-screen min-w-screen flex justify-center items-center bg-neutral">
		<div class="modal-box">
			<div class="flex w-ful align-right justify-end">
				<button on:click={() => (showModal = !showModal)}> <RiCloseFill size={'1.5em'} /></button>
			</div>
			<h3 class="font-bold text-lg mb-4">SERVER NAME</h3>
			<input
				bind:value={serverName}
				type="text"
				placeholder="moe..."
				class="input input-lg bg-neutral opacity-75 focus:outline-none w-full input-bordered text-white"
			/>
			<div class="modal-action">
				<form method="dialog">
					<button class="btn btn-info" on:click={() => createServer()}>Create</button>
				</form>
			</div>
		</div>
	</div>
{/if}

{#if !showModal}
	<div class="flex flex-row min-w-screen min-h-screen">
		<div class="w-20 flex px-4 flex-col bg-base-300 justify-between items-center">
			<div>
				{#each data.servers as server}
					<div class="mt-2 avatar ring rounded-full">
						<button class="rounded-full">
							<img
								src="https://avatars.githubusercontent.com/u/46141275?v=4"
								alt="nekonako"
								class="rounded-full"
							/>
						</button>
					</div>
				{/each}
			</div>
			<div class="mb-2">
				<button
					class="rounded-full bg-base-100 justify-center p-2"
					on:click={() => (showModal = !showModal)}
				>
					<RiAddLine size="2em" class="text-success text-center" />
				</button>
			</div>
		</div>
		<div class="w-1/6 flex flex-col bg-base-200">
			{selectedServer}
			<div class="avatar">
				<div class="w-full h-48 rounded">
					<img src="kato.jpg" class="inline" alt="kana" />
				</div>
			</div>
			<ul class="menu">
				{#each data.channels as channel}
					<li>
						<a href="#">
							{channel.name}
						</a>
					</li>
				{/each}
			</ul>
		</div>
		<div class=" w-3/4 min-h-screen">
			<div class="flex flex-col overflow-y-auto max-h-screen justify-between scrollbar-hide">
				<div class="flex flex-col h-full mb-4 px-6">
					{#each messages as message}
						<div class="chat chat-start">
							<div class="chat-image avatar">
								<div class="w-10 rounded-full">
									<img src="https://avatars.githubusercontent.com/u/46141275?v=4" />
								</div>
							</div>
							<div class="chat-header mb-2 mt-4">
								from: {message.sender_id}
								<time class="text-xs opacity-50">2 hours ago</time>
							</div>
							<div class="chat-bubble chat-bubble-primary">{message.content}</div>
						</div>
					{/each}
				</div>
				<div class="sticky bottom-0 bg-base-200">
					<div class="flex flex-row items-center m-4 rounded-lg bg-base-100">
						<input
							bind:value={message}
							type="text"
							placeholder="Send message"
							class="input input-lg focus:outline-none w-full placeholder-base-content text-base"
						/>
						<div class="h-full flex bg-base-100 min-h-full p-4">
							<button
								class="bg-base-100 h-full min-h-full text-white"
								on:click={() => sendMessage()}
							>
								<RiSendPlaneFill size={'2em'} class="text-success" />
							</button>
						</div>
					</div>
				</div>
			</div>
		</div>
		<div class="w-1/6 bg-base-200">
			<div class="flex flex-col justify-between min-h-screen">
				<div>members</div>
			</div>
		</div>
	</div>
{/if}
