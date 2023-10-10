<script lang="ts">
	import {
		RiChatVoiceLine,
		RiMessage2Line,
		RiCircleFill,
		RiArrowDownSLine,
		RiSettings4Line,
		RiCloseLine
	} from 'svelte-remixicon';
	import type {
		Channel,
		CreateChannelCategoryRequest,
		CreateChannelRequest,
		Profile,
		Server
	} from '$lib/service/type';
	import { ShowServerSettingModal, ShowUserSettingModal, currentChannel } from './store';
	import { createChannel, createChannelCategory, getListChannel } from '$lib/service/channel';

	export let selectedServer: Server;
	export let channels: Array<Channel>;
	export let profile: Profile;

	let showPopupServer = false;
	let inviteModal = false;
	let categoryModal = false;
	let createChannelModal = false;

	let category: CreateChannelCategoryRequest = {
		name: '',
		is_private: false,
		server_id: selectedServer.id
	};

	let channel: CreateChannelRequest = {
		name: '',
		server_id: selectedServer.id,
		category_id: '',
		is_private: false,
		type: ''
	};

	function selectChannel(channel_id: string, channel_type: string) {
		currentChannel.set({
			channel_type: channel_type,
			channel_id: channel_id
		});
	}

	function copy() {
		let input = document.getElementById('invite-friend') as HTMLInputElement;
		input.select();
		document.execCommand('copy');
		navigator.clipboard.writeText(input.value);
	}

	async function handleCreateChannelCategory() {
		const response = await createChannelCategory(fetch, category);
		if (response.code == 200) {
			categoryModal = false;
			const response = await getListChannel(fetch, selectedServer.id);
			channels = response.data;
		}
	}

	async function handleCreateChannel() {
		const response = await createChannel(fetch, channel);
		if (response.code == 200) {
			createChannelModal = false;
			const response = await getListChannel(fetch, selectedServer.id);
			channels = response.data;
		}
	}
</script>

{#if inviteModal}
	<div class="absolute inset-0 z-20 min-h-screen min-w-screen overflow-disable">
		<div class="flex flex-row w-full h-full justify-center items-center bg-base-300 bg-opacity-95">
			<div class="bg-base-100 p-8 flex flex-col w-1/4">
				<div class="text-right w-full flex flex-row justify-between items-center mb-8">
					<div>
						<span class="font-bold mb-8">Invite firends to {selectedServer.name}</span>
					</div>
					<button on:click={() => (inviteModal = false)}>
						<RiCloseLine />
					</button>
				</div>
				<span class="mb-2">Share Server Invite link to a friend</span>
				<div class="flex flex-row select-none w-full">
					<input
						id="invite-friend"
						type="text"
						class="bg-base-200 px-4 p-3 w-full select-none focus:outline-none"
						value={'http://localhost:3000/invite/' + selectedServer.id}
					/>
					<button class="bg-success p-3 text-base-100 px-2" on:click={copy}>copy</button>
				</div>
			</div>
		</div>
	</div>
{/if}

{#if categoryModal}
	<div class="absolute inset-0 z-20 min-h-screen min-w-screen overflow-disable">
		<div class="flex flex-row w-full h-full justify-center items-center bg-base-300 bg-opacity-95">
			<div class="bg-base-100 p-8 flex flex-col w-1/4 rounded-md">
				<div class="text-right w-full flex flex-row justify-between items-center mb-8">
					<div>
						<span class="font-bold mb-8 text-lg">Create Category</span>
					</div>
					<button on:click={() => (categoryModal = false)}>
						<RiCloseLine />
					</button>
				</div>
				<span class="mb-2 font-semibold">Category Name</span>
				<div class="flex flex-row select-none w-full">
					<input
						id="invite-friend"
						type="text"
						placeholder="New Category"
						class="bg-base-200 px-4 p-3 w-full text-sm select-none focus:outline-none"
						bind:value={category.name}
					/>
				</div>
				<div class="flex flex-row justify-between items-center mt-6">
					<div class="font-semibold">Private Category</div>
					<input type="checkbox" class="toggle toggle-success" bind:checked={category.is_private} />
				</div>
				<span class="text-xs mt-2">
					By making category private, only selected member will able to view this category.
				</span>
				<div class="flex flex-row justify-between mt-8">
					<div />
					<button
						class="bg-success text-base-100 px-3 py-2 rounded"
						on:click={handleCreateChannelCategory}>Save</button
					>
				</div>
			</div>
		</div>
	</div>
{/if}

{#if createChannelModal}
	<div class="absolute inset-0 z-20 min-h-screen min-w-screen overflow-disable">
		<div class="flex flex-row w-full h-full justify-center items-center bg-base-300 bg-opacity-95">
			<div class="bg-base-100 p-8 flex flex-col w-1/4 rounded-md">
				<div class="text-right w-full flex flex-row justify-between items-center mb-8">
					<div>
						<span class="font-bold text-lg">Create Channel</span>
					</div>
					<button on:click={() => (createChannelModal = false)}>
						<RiCloseLine />
					</button>
				</div>

				<span class="mb-2 font-bold uppercase text-xs">Channel Type</span>
				<div class="form-control">
					<label class="label cursor-pointer">
						<div class="flex flex-row gap-x-2">
							<RiMessage2Line />
							<span class="label-text">Text</span>
						</div>
						<input
							type="radio"
							name="channel_type"
							value="text"
							class="radio checked:bg-success radio-sm"
							bind:group={channel.type}
						/>
					</label>
				</div>
				<div class="form-control">
					<label class="label cursor-pointer">
						<div class="flex flex-row gap-x-2">
							<RiChatVoiceLine />
							<span class="label-text">Voice</span>
						</div>
						<input
							type="radio"
							name="channel_type"
							value="voice"
							class="radio checked:bg-success radio-sm"
							bind:group={channel.type}
						/>
					</label>
				</div>

				<span class="mb-2 font-bold mt-6 uppercase text-xs">Channel Name</span>
				<div class="flex flex-row select-none w-full">
					<input
						id="invite-friend"
						type="text"
						placeholder="New Channel"
						class="bg-base-200 px-4 p-2 rounded w-full text-sm select-none focus:outline-none"
						bind:value={channel.name}
					/>
				</div>

				<span class="mb-2 font-bold mt-6 uppercase text-xs mt-6">Category</span>
				<select
					bind:value={channel.category_id}
					class="text-sm rounded w-full py-2 bg-base-200 px-4 focus:outline-none border-r-8 border-rounded border-r-base-200"
				>
					{#each channels as channel}
						<option value={channel.category_id}>{channel.category_name}</option>
					{/each}
				</select>

				<div class="flex flex-row justify-between items-center mt-6">
					<div class="font-bold text-xs uppercase">Private channel</div>
					<input type="checkbox" class="toggle toggle-success" bind:checked={channel.is_private} />
				</div>
				<span class="text-xs mt-2">
					By making category private, only selected member will able to view this category.
				</span>
				<div class="flex flex-row justify-between mt-8">
					<div />
					<button class="bg-success text-base-100 px-3 py-2 rounded" on:click={handleCreateChannel}
						>Save</button
					>
				</div>
			</div>
		</div>
	</div>
{/if}

<div class="avatar top-0">
	<button on:click={() => (showPopupServer = !showPopupServer)}>
		<div
			class="absolute bottom-0 px-4 py-3 bg-base-100/70 w-full flex justify-between flex-row items-center"
		>
			<div class="flex flex-row w-full items-center">
				<RiCircleFill class="mr-2 text-info" size="1.2em" />
				<div>
					<span class="text-white font-bold"> {selectedServer.name}</span>
				</div>
			</div>
			<RiArrowDownSLine class="text-white text-left" />
		</div>
	</button>
	<img src="/kato.jpg" class="inline" alt="kana" />
</div>
{#if showPopupServer}
	<div class="rounded-lg bg-base-300 flex m-2">
		<ul class="menu menu-md w-full">
			<li on:click={() => ShowServerSettingModal.set(true)}>
				<a>Server Setting</a>
			</li>
			<li on:click={() => (inviteModal = true)}>
				<a>Invite People</a>
			</li>
			<li on:click={() => (createChannelModal = true)}>
				<a>Create Channel</a>
			</li>
			<li on:click={() => (categoryModal = true)}>
				<a>Create Category</a>
			</li>
			<li class="text-error">
				<a>Leave Server</a>
			</li>
		</ul>
	</div>
{/if}
<ul class="menu">
	{#each channels as channelCategory}
		<li class="menu-title">
			<a href="#">
				{channelCategory.category_name.toUpperCase()}
			</a>
		</li>
		{#each channelCategory.channels as channel}
			{#if $currentChannel.channel_id === channel.id}
				<li class="bg-warning/20">
					<a href="#" on:click={() => selectChannel(channel.id, channel.channel_type)}>
						{#if channel.channel_type == 'text'}
							<RiMessage2Line size="1.5em" />
						{:else}
							<RiChatVoiceLine size="1.5em" />
						{/if}
						{channel.name}
					</a>
				</li>
			{:else}
				<li>
					<a href="#" on:click={() => selectChannel(channel.id, channel.channel_type)}>
						{#if channel.channel_type == 'text'}
							<RiMessage2Line size="1.5em" />
						{:else}
							<RiChatVoiceLine size="1.5em" />
						{/if}
						{channel.name}
					</a>
				</li>
			{/if}
		{/each}
	{/each}
</ul>
<div
	class="sticky bottom-0 top-[100vh] p-4 bg-base-300 flex flex-row items-center justify-between hover:cursor-pointer"
	on:click={() => ShowUserSettingModal.set(true)}
>
	<div class="flex flex-row items-center group">
		<div class="avatar w-8 rounded-full online">
			<img src={profile.avatar} alt="profile" class="rounded-full" />
		</div>

		<div class="ml-4 font-bold">{profile.username}</div>
	</div>
	<div><RiSettings4Line size="1.2em" class="justify-end" /></div>
</div>
