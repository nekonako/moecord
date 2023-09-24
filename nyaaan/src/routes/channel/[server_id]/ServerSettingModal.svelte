<script lang="ts">
	import type { Server } from './type';
	import { ShowServerSettingModal } from './store';
	import { RiCloseLine } from 'svelte-remixicon';

	export let selectedServer: Server;
	let file: File;
	let updateServerName = selectedServer.name;
	let serverImg: HTMLElement;

	async function updateServer() {
		const rBody = new FormData();
		rBody.append('id', selectedServer.id);
		rBody.append('name', selectedServer.name);
		rBody.append('avatar', file);

		const response = await fetch(`/api/servers`, {
			method: 'PUT',
			body: rBody
		});
		const result = await response.json();
		if (result.code == 200) {
			ShowServerSettingModal.set(false);
		}
	}

	function onFileSelected(event: Event) {
		const target = event.target as unknown as { files: File[] };
		file = target?.files[0];
		const reader = new FileReader();
		reader.addEventListener('load', function () {
			serverImg.setAttribute('src', reader.result?.toString()!);
		});
		reader.readAsDataURL(file);
	}
</script>

<div class="min-h-screen min-w-screen flex flex-row bg-base-200">
	<div class="bg-base-300 w-1/3 flex flex-row justify-between text-right pt-4">
		<div />
		<div class="w-24 mx-4">
			<ul class="menu">
				<li class="menu-title text-lg">{selectedServer.name}</li>
				<li><a>Overview</a></li>
			</ul>
		</div>
	</div>
	<div class="flex flex-col gap-4 bg-base-200 w-2/3 p-8 pt-8">
		<div class=" flex flex-row justify-between">
			<span class="font-semibold">Server Overview</span>
			<button
				class="ring ring-base-content rounded-full p-2"
				on:click={() => ShowServerSettingModal.set(false)}
			>
				<RiCloseLine />
			</button>
			<div />
		</div>
		<div class="flex flex-row gap-6">
			<div class="flex flex-col basis-1">
				<div class="mt-2 avatar ring rounded-full ring-info w-20 h-20">
					<button class="rounded-full">
						<img
							src={selectedServer.avatar}
							bind:this={serverImg}
							alt="nekonako"
							class="rounded-full"
						/>
					</button>
				</div>
				<div class="bg-base-300 text-center mt-4 px-4 py-1 border-base-300 border">
					<input
						on:change={onFileSelected}
						type="file"
						id="upload-server-avatar"
						class="file-input hidden display-none file-input-xs"
					/>
					<button on:click={() => document?.getElementById('upload-server-avatar')?.click()}>
						upload</button
					>
				</div>
			</div>
			<div class="basis-2">
				<label class="label">
					<span class="label-text">Server name</span>
				</label>
				<input class="input focus:outline-none" type="text" bind:value={updateServerName} />
			</div>
		</div>
		<div class="inline-block text-center">
			<button on:click={() => updateServer()} class="bg-success text-base-100 px-4">Save</button>
		</div>
		<div />
	</div>
</div>
