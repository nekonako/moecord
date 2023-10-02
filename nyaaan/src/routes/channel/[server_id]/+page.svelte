<script lang="ts">
	import { browser } from '$app/environment';
	import type { Message } from '$lib/service/type';
	import {
		showCreateServerModal,
		ShowServerSettingModal,
		ShowUserSettingModal,
		MapServerMember
	} from './store';
	import CreateServerModal from './CreateServerModal.svelte';
	import ServerSettingModal from './ServerSettingModal.svelte';
	import ListServer from './ListServer.svelte';
	import ListChannel from './ListChannel.svelte';
	import ListMessage from './ListMessage.svelte';
	import ListMember from './ListMember.svelte';
	import UserSettingModal from './UserSetting.svelte';
	import { onMount } from 'svelte';

	export let data;
	let ws: WebSocket;
	let wsEvent: ((this: WebSocket, ev: MessageEvent) => any) | null;
	let wsMessageHandler: Set<(e: MessageEvent) => void> = new Set();

	onMount(() => {
		for (let i = 0; i < data.server_member.length; i++) {
			var o: { [key: string]: boolean } = {};
			o[data.server_member[i].user_id] = false;
			MapServerMember.set({ ...$MapServerMember, ...o });
		}
	});

	let messages: Array<Message>;
	let selectedServer = data.selected_server;

	function connect_ws() {
		ws = new WebSocket('ws://localhost:4000/ws');
		ws.onclose = function () {
			setTimeout(function () {
				connect_ws();
			}, 1000);
		};
		ws.onmessage = (e: MessageEvent) => {
			wsMessageHandler.forEach((handler) => handler(e));
		};
		wsEvent = ws.onmessage;
	}

	if (browser) {
		connect_ws();
	}
</script>

{#if $showCreateServerModal}
	<CreateServerModal on:getListServer={(e) => (data.servers = e.detail)} />
{/if}

{#if $ShowServerSettingModal}
	<ServerSettingModal
		selectedServer={data.selected_server}
		on:updateServer={(e) => (data.selected_server = e.detail)}
	/>
{/if}

{#if $ShowUserSettingModal}
	<UserSettingModal profile={data.profile} />
{/if}

{#if !$showCreateServerModal && !$ShowServerSettingModal}
	<div class="flex flex-row min-w-screen min-h-screen">
		<div class="w-20 flex px-4 flex-col bg-base-300 justify-between items-center">
			<ListServer servers={data.servers} {selectedServer} />
		</div>
		<div class="w-1/6 flex flex-col bg-base-200">
			<ListChannel
				selectedServer={data.selected_server}
				channels={data.channels}
				profile={data.profile}
			/>
		</div>
		<div class=" w-3/4 min-h-screen">
			<ListMessage
				{messages}
				{ws}
				profile={data.profile}
				server={data.selected_server}
				{wsMessageHandler}
			/>
		</div>
		<div class="w-1/6 bg-base-200 p-4 min-h-screen">
			<ListMember members={data.server_member} {wsMessageHandler} />
		</div>
	</div>
{/if}
