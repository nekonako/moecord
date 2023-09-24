<script lang="ts">
	import { RiCloseLine } from 'svelte-remixicon';
	import { ShowUserSettingModal } from './store';
	import type { Profile, Server } from './type';

	export let profile: Profile;

	let file: File;
	let serverImg: HTMLElement;
	async function updateProfile() {
		const rBody = new FormData();
		rBody.append('id', profile.id);
		rBody.append('username', profile.username);
		rBody.append('email', profile.email);
		rBody.append('avatar', file);

		const response = await fetch(`/api/profile`, {
			method: 'PUT',
			body: rBody
		});
		const result = await response.json();
		if (result.code == 200) {
			ShowUserSettingModal.set(false);
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
	<div class="bg-base-300 w-1/3 flex flex-row justify-between text-left pt-8">
		<div />
		<div class="w-36 mx-4">
			<ul class="w-full">
				<li class="font-semibold uppercase text-sm mb-4 px-4">User Setting</li>
				<li class=" hover:cursor-pointer hover:bg-base-200 py-2 px-4">Profile</li>
			</ul>
		</div>
	</div>
	<div class="flex flex-col gap-4 bg-base-200 w-2/3 p-8 pt-8 pl-12">
		<div class=" flex flex-row justify-between">
			<span class="font-semibol text-xl text-white/90">Profile</span>
			<button
				class="ring ring-base-content rounded-full p-2"
				on:click={() => ShowUserSettingModal.set(false)}
			>
				<RiCloseLine />
			</button>
			<div />
		</div>
		<div class="flex flex-row gap-24 w-3/6">
			<div class="flex flex-col">
				<div class="mt-2 avatar ring rounded-full ring-info w-32 h-32 text-center">
					<button class="rounded-full">
						<img src={profile.avatar} bind:this={serverImg} alt="nekonako" class="rounded-full" />
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
			<div class="w-full">
				<label class="label">
					<span class="font-semibold text-xs uppercase">Username</span>
				</label>
				<input
					class="input input-md focus:outline-none w-full"
					type="text"
					bind:value={profile.username}
				/>
				<label class="label mt-4">
					<span class="text-xs font-semibold uppercase">Email</span>
				</label>
				<input
					class="input input-md focus:outline-none w-full"
					type="text"
					bind:value={profile.email}
				/>
			</div>
		</div>
		<div class="inline-block text-center mt-8 rounded">
			<button on:click={() => updateProfile()} class="bg-success text-base-100 px-4 rounded py-2"
				>Save</button
			>
		</div>
		<div />
	</div>
</div>
