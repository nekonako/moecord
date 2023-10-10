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
	import * as mediasoupClient from 'mediasoup-client';
	import type { DtlsParameters, IceCandidate, IceParameters } from 'mediasoup-client/lib/types';
	import { dev } from '$app/environment';

	type MediaTrack = {
		stream: HTMLMediaElement;
		state?: string;
		user?: Profile;
		id?: string;
	};

	export let messages: Array<Message> = [];
	export let ws: WebSocket;
	export let wsMessageHandler: Set<(e: MessageEvent) => void> = new Set();

	export let profile: Profile;
	export let server: Server;
	let newMessage = Array<Message>();
	let message: string;
	let messageEl: HTMLElement;
	let timer: ReturnType<typeof setTimeout> | undefined;

	let mediaTracks: Array<MediaTrack> = [];
	let ss: HTMLMediaElement;

	$: typing(message);

	$: ref = mediaTracks;

	async function updateStream(tracks: Array<MediaTrack>) {
		console.log(tracks, 'cok 123');
		tracks.forEach((value) => {
			console.log(value.stream.srcObject);
			// setInterval(() => {
			// 	console.log(value.stream.srcObject, "+++++++")
			// }, 1000)
		});
	}

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
			try {
				await createVoiceRoomToken();
				const device = new mediasoupClient.Device();
				const dtls: DtlsParameters = {
					fingerprints: []
				};
				const iceCandidate: IceCandidate[] = [];
				const iceParameter: IceParameters = {
					usernameFragment: '12',
					password: '123',
					iceLite: false
				};
				const transport = await device.createSendTransport({
					dtlsParameters: dtls,
					iceCandidates: iceCandidate,
					iceParameters: iceParameter,
					id: 'i'
				});
			} catch (err) {
				console.log(err);
			}
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
		try {
			console.log('cok 123');
			const media = await navigator.mediaDevices.getUserMedia({
				audio: true,
				video: true
			});

			const ws = new WebSocket('ws://localhost:4000/ws/voice/' + $currentChannel.channel_id);
			console.log(ws);

			console.log('start...');
		} catch (err) {
			console.log(err);
			throw err;
		}
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

{#if $currentChannel.channel_type == 'text'}
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
{/if}

{#if $currentChannel.channel_type == 'voice'}
	<div>
		{#each ref as mediaTrack, i}
			<!-- <mediaTrack.element/> -->
			<!-- {mediaTrack.element.srcObject} -->
			<!-- {console.log(ref[i].stream, '=====')} -->
			<video bind:this={ref[i].stream} width="400" controls autoplay playsinline />
		{/each}
		<!-- <video bind:this={ss} width="400" controls autoplay /> -->
	</div>
{/if}
