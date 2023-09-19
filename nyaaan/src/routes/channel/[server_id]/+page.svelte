<script lang="ts">
	import {
		RiAddLine,
		RiChatVoiceLine,
		RiCloseFill,
		RiMessage2Line,
		RiSendPlaneFill,
		RiCircleFill,
		RiArrowUpSLine,
		RiArrowDownSLine,
		RiSendPlane2Line,
		RiSendPlaneLine
	} from 'svelte-remixicon';
	import { browser } from '$app/environment';
	import type { Message, Channel } from './type';
	export let data;
	let ws: WebSocket | undefined;
	let message: string;
	let messages: Array<Message> = [];
	let showModal = false;
	let serverName = '';
	let selectedServer = data.selected_server;
	let newMessage = Array<Message>();
	let members = data.server_member;
	let showPopupServer = false;

	if (browser) {
		ws = new WebSocket('ws://localhost:4000/ws');
		handleMessage(ws);
	}

	type SendMessage = {
		channel_id: string;
		content: string;
	};

	type CreateServer = {
		name: string;
	};

	async function sendMessage() {
		const payload: SendMessage = {
			channel_id: data.channels[0].channels[0].id,
			content: message
		};
		const response = await fetch('/api/messages', {
			body: JSON.stringify(payload),
			method: 'POST'
		});
		const result = await response.json();
		if (result.code == 200) {
			message = '';
		}
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

	async function getMessageByChannel(channelID: string) {
		const response = await fetch(`/api/messages/channels/${channelID}`);
		const result = await response.json();
		messages = result.data;
	}

	function handleMessage(ws: WebSocket) {
		ws.onmessage = (m) => {
			newMessage = [...newMessage, JSON.parse(m.data)];
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
					{#if server.id == selectedServer.id}
						<div class="mt-2 avatar ring rounded-full ring-info">
							<button class="rounded-full">
								<img
									src="https://avatars.githubusercontent.com/u/46141275?v=4"
									alt="nekonako"
									class="rounded-full"
								/>
							</button>
						</div>
					{:else}
						<div class="mt-2 avatar rounded-full">
							<button
								class="rounded-full"
								on:click={() => {
									selectedServer = server;
									window.location.href = `/channel/${server.id}`;
								}}
							>
								<img
									src="https://avatars.githubusercontent.com/u/46141275?v=4"
									alt="nekonako"
									class="rounded-full"
								/>
							</button>
						</div>
					{/if}
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
			<div class="avatar">
				<div class="w-full h-48">
					<button on:click={() => (showPopupServer = !showPopupServer)}>
						<div
							class="absolute bottom-0 px-4 py-3 bg-base-300/70 w-full flex justify-between flex-row items-center"
						>
							<div class="flex flex-row w-full">
								<RiCircleFill class="mr-2 text-info" size="1.5em" />
								{selectedServer.name}
							</div>
							<RiArrowDownSLine class="text-white text-left" />
						</div>
					</button>
					<img src="/kato.jpg" class="inline" alt="kana" />
				</div>
			</div>
			{#if showPopupServer}
				<div class="rounded-lg bg-base-300 flex m-2">
					<ul class="menu menu-md w-full">
						<li class="">
							<a>Edit Server Profile </a>
						</li>
						<li class="text-error">
							<a>Leave Server</a>
						</li>
					</ul>
				</div>
			{/if}
			<ul class="menu">
				{#each data.channels as channelCategory}
					<li class="menu-title">
						<a href="#">
							{channelCategory.category_name.toUpperCase()}
						</a>
					</li>
					{#each channelCategory.channels as channel}
						<li>
							<a href="#" on:click={() => getMessageByChannel(channel.id)}>
								{#if channel.channel_type == 'text'}
									<RiMessage2Line size="1.5em" />
								{:else}
									<RiChatVoiceLine size="1.5em" />
								{/if}
								{channel.name}
							</a>
						</li>
					{/each}
				{/each}
			</ul>
			<div class="sticky bottom-0 top-[100vh] p-4 bg-base-300 flex flex-row items-center">
				<div class="avatar w-8 rounded-full online">
					<img
						src="https://avatars.githubusercontent.com/u/46141275?v=4"
						alt="nekonako"
						class="rounded-full"
					/>
				</div>
				<div class="ml-4 font-bold">yuune</div>
			</div>
		</div>
		<div class=" w-3/4 min-h-screen">
			<div
				class="flex flex-col overflow-y-auto max-h-screen min-h-screen justify-between scrollbar-hide"
			>
				<div class="flex flex-col mb-4 px-6">
					{#each messages as message}
						<div class="chat chat-start">
							<div class="chat-image avatar">
								<div class="w-10 rounded-full">
									<img src="https://avatars.githubusercontent.com/u/46141275?v=4" />
								</div>
							</div>
							<div class="chat-header mb-2 mt-4 opacity-70">
								{message.username}
								<time class="ml-2 text-xs opacity-50"
									>{new Date(message.created_at).toLocaleString()}</time
								>
							</div>
							<div class="chat-bubble chat-bubble-primary">{message.content}</div>
						</div>
					{/each}
					{#if newMessage.length > 0}
						<div class="fle flex row items-center mt-4">
							<div class="bg-error px-2 text-neutral rounded-lg">new</div>
							<div class="w-full h-px bg-error" />
						</div>
					{/if}
					{#each newMessage as message}
						<div class="chat chat-start">
							<div class="chat-image avatar">
								<div class="w-10 rounded-full">
									<img src="https://avatars.githubusercontent.com/u/46141275?v=4" />
								</div>
							</div>
							<div class="chat-header mb-2 text-xs opacity-50">
								{message.username}
								<time>{message.created_at}</time>
							</div>
							<div class="chat-bubble chat-bubble-primary">{message.content}</div>
						</div>
					{/each}
				</div>
				<div class="sticky bottom-0 bg-base-200">
					<div class="flex flex-row items-center m-4 bg-base-100">
						<input
							bind:value={message}
							type="text"
							placeholder="Send message"
							class="input input-lg focus:outline-none w-full placeholder-base-content text-base"
						/>
						<div class="h-full flex bg-base-100 min-h-full p-4 border-l-2 border-l-base-200">
							<button
								class="bg-base-100 h-full min-h-full text-white"
								on:click={() => sendMessage()}
							>
								<RiSendPlane2Line size={'2em'} class="text-primary" />
							</button>
						</div>
					</div>
				</div>
			</div>
		</div>
		<div class="w-1/6 bg-base-200 p-4 min-h-screen">
			{#each members as member}
				<div class="flex flex-row items-center">
					<div class="avatar online mr-4">
						<div class="w-8 ring ring-success rounded-full">
							<img src="https://avatars.githubusercontent.com/u/46141275?v=4" alt="profile" />
						</div>
					</div>
					<span class="font-medium"> {member.username} </span>
				</div>
			{/each}
		</div>
	</div>
{/if}
