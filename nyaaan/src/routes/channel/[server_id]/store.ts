import { writable } from 'svelte/store';

export const count = writable(0);
export const showCreateServerModal = writable(false);
export const ShowServerSettingModal = writable(false);
export const currentChannel = writable({
  channel_type: '',
  channel_id: ''
});
export const ShowUserSettingModal = writable(false)
export const MapServerMember = writable<{ [key: string]: boolean }>({});
