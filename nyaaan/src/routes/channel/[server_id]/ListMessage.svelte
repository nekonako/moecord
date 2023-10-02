<script lang="ts">
	import type { ApiResponse, Message, Profile, Server, WebsocketMessage } from '$lib/service/type';
	import { getColor } from '$lib/util';
	import { MapServerMember, currentChannel } from './store';
	import { RiSendPlane2Line } from 'svelte-remixicon';
	import { afterUpdate, onMount } from 'svelte';
	import {
		Room,
		createLocalTracks,
		RoomEvent,
		RemoteTrack,
		RemoteParticipant
	} from 'livekit-client';

	let localTrack: HTMLMediaElement;
	let remoteTrack: HTMLMediaElement;

	export let messages: Array<Message> = [];
	export let ws: WebSocket;
	export let wsMessageHandler: Set<(e: MessageEvent) => void> = new Set();

	export let profile: Profile;
	export let server: Server;
	let newMessage = Array<Message>();
	let message: string;
	let messageEl: HTMLElement;
	let timer: number | undefined;

	$: typing(message);

	function handleMessage() {
		wsMessageHandler.add((e) => {
			let message: WebsocketMessage<any> = JSON.parse(e.data);
			if (message.event_id == 'NEW_CHANNEL_MESSAGE') {
				let data = message.data as SaveMessageResponse;
				if ($currentChannel.channel_id === data.channel_id) {
					newMessage = [...newMessage, data];
				}
			}
			if (message.event_id == 'TYPING') {
				let data = message.data as typingMessage;
				let o: { [key: string]: boolean } = {};
				o[data.user_id] = true;
				MapServerMember.set({ ...$MapServerMember, ...o });
				return;
			}
			if (message.event_id == 'STOP_TYPING') {
				let data = message.data as typingMessage;
				let o: { [key: string]: boolean } = {};
				o[data.user_id] = false;
				MapServerMember.set({ ...$MapServerMember, ...o });
			}
		});
	}

	type typingMessage = {
		server_id: string;
		user_id: string;
		channel_id: string;
	};

	type SendMessage = {
		channel_id: string;
		content: string;
	};

	type SaveMessageResponse = {
		id: string;
		channel_id: string;
		sender_id: string;
		content: string;
		created_at: string;
		updated_at: string;
		username: string;
		avatar: string;
	};

	async function sendMessage() {
		if (!message) return;
		stopTyping();

		const payload: SendMessage = {
			channel_id: $currentChannel.channel_id,
			content: message
		};

		const response = await fetch('/api/messages', {
			body: JSON.stringify(payload),
			method: 'POST'
		});

		const result: ApiResponse<SaveMessageResponse> = await response.json();

		if (result.code == 200) {
			message = '';
			let data = result.data;
			newMessage = [...newMessage, result.data];
			scrollBottom(messageEl);
			let o: { [key: string]: boolean } = {};
			o[data.sender_id] = false;
			MapServerMember.set({ ...$MapServerMember, ...o });
		}
	}

	currentChannel.subscribe(async (channel) => {
		if (channel.channel_type == 'text') {
			const response = await fetch(`/api/messages/channels/${channel.channel_id}`);
			const result = await response.json();
			if (result.code == 200) {
				messages = result.data;
				newMessage = [];
			}
			await scrollBottom(messageEl);
		}
		if (channel.channel_type == 'voice') {
			await createVoiceRoomToken();
		}
	});

	async function scrollBottom(el: HTMLElement) {
		el.scroll({ top: el.scrollHeight, behavior: 'instant' });
	}

	onMount(() => {
		scrollBottom(messageEl);
		handleMessage();
	});

	afterUpdate(() => {
		scrollBottom(messageEl);
	});

	async function createVoiceRoomToken() {
		const response = await fetch('/api/room/token');
		const result = await response.json();

		const wsURL = 'ws://localhost:7880';
		const token = result.data.token;

		const room = new Room();
		await room.connect(wsURL, token);

		room.on(RoomEvent.TrackSubscribed, (track, publication, participant) => {
			attachTrack(track, participant);
		});

		const tracks = await createLocalTracks({
			audio: true,
			video: true
		});
		for (let track of tracks) {
			await room.localParticipant.publishTrack(track);
			track.attach(localTrack);
		}
	}

	function attachTrack(track: RemoteTrack, participant: RemoteParticipant) {
		track.attach(remoteTrack);
	}

	function stopTyping() {
		let stopTyping = {
			event_id: 'STOP_TYPING',
			data: {
				user_id: profile.id,
				channel_id: $currentChannel.channel_id,
				server_id: server.id
			}
		};
		ws.send(JSON.stringify(stopTyping));
	}

	function typing(x: string) {
		if (!x || x === '') {
			if (ws && ws.readyState === ws.OPEN) {
				stopTyping();
			}
			return;
		}
		if (!timer) {
			let message: WebsocketMessage<typingMessage> = {
				event_id: 'TYPING',
				data: {
					channel_id: $currentChannel.channel_id,
					user_id: profile.id,
					server_id: server.id
				}
			};
			ws.send(JSON.stringify(message));
			timer = setTimeout(() => {
				timer = undefined;
			}, 2000);
		}
	}
</script>

<div
	class="flex flex-col overflow-y-auto max-h-screen min-h-screen justify-between scrollbar-hide"
	bind:this={messageEl}
>
	<div class="flex flex-col mb-4 px-6">
		{#each messages as message}
			<div class="chat chat-start">
				<div class="chat-image avatar">
					<div class="w-10 rounded-full">
						<img src={message.avatar} alt="avatar" />
					</div>
				</div>
				<div class="chat-header mb-2 mt-4">
					<span class={getColor(message.username.toUpperCase().charCodeAt(0))}
						>{message.username}</span
					>
					<time class="ml-2 text-xs opacity-50"
						>{new Date(message.created_at).toLocaleString()}</time
					>
				</div>
				{#if message.sender_id == profile.id}
					<div class="chat-bubble chat-bubble-primary">{message.content}</div>
				{:else}
					<div class="chat-bubble chat-bubble-secondary">{message.content}</div>
				{/if}
			</div>
		{/each}
		{#if newMessage.length > 0}
			<div class="fle flex row justify-between items-center mt-4">
				<div class="w-full h-px bg-warning/10 w-full" />
				<div class=" px-2 text-warning w-72 text-center items-center rounded-lg">new message</div>
				<div class="w-full h-px bg-warning/10 w-full" />
			</div>
		{/if}
		{#each newMessage as message}
			<div class="chat chat-start">
				<div class="chat-image avatar">
					<div class="w-10 rounded-full">
						<img src={message.avatar} alt="avatar" />
					</div>
				</div>
				<div class="chat-header mb-2">
					<span class={getColor(message.username.toUpperCase().charCodeAt(0))}
						>{message.username}</span
					>
					<time class="text-xs opacity-50 ml-2"
						>{new Date(message.created_at).toLocaleString()}</time
					>
				</div>
				{#if message.sender_id == profile.id}
					<div class="chat-bubble chat-bubble-primary">{message.content}</div>
				{:else}
					<div class="chat-bubble chat-bubble-secondary">{message.content}</div>
				{/if}
			</div>
		{/each}
	</div>
	<div class="sticky bottom-0 bg-base-200">
		<div class="flex flex-row items-center m-4 bg-base-100 rounded-lg">
			<input
				bind:value={message}
				type="text"
				placeholder="Send message"
				class="input input-lg rounded-lg focus:outline-none w-full placeholder-base-content text-base"
			/>
			<div
				class="h-full flex bg-base-100 min-h-full p-4 border-l-2 border-l-base-200 rounded-lg rounded-l-none"
			>
				<button class="bg-base-100 h-full min-h-full text-white" on:click={() => sendMessage()}>
					<RiSendPlane2Line size={'2em'} class="text-success" />
				</button>
			</div>
		</div>
	</div>
</div>
