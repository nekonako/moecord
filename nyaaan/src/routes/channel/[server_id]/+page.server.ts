import { getListServer, getServerMember } from '$lib/service/server'
import { getListChannel } from "$lib/service/channel";
import type { Server } from "$lib/service/type";
import { getUserProfile } from "$lib/service/user";


export async function load({ fetch, params }) {

  const profile = await getUserProfile(fetch);
  const servers = await getListServer(fetch);
  const channels = await getListChannel(fetch, params.server_id);

  let selected_server: Server;

  servers.data.forEach((val: Server) => {
    if ((params.server_id == "@me" && val.name == "@me") || (params.server_id != "@me" && params.server_id == val.id)) {
      selected_server = val;
    }
  });

  const serverMember = await getServerMember(fetch, selected_server!.id)

  return {
    servers: servers.data,
    channels: channels.data,
    server_member: serverMember.data,
    profile: profile.data,
    selected_server: selected_server!,
  };
}
