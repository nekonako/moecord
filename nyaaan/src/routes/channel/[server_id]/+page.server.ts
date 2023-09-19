import { redirect } from '@sveltejs/kit';
import type { ApiResponse, Channel, Message, Server, Servermember } from './type';

export async function load({ fetch, setHeaders, cookies, params }) {
  const responseServer = await fetch('/api/servers');
  const servers: ApiResponse<Server> = await responseServer.json();
  if (servers.code == 401) {
    throw redirect(307, '/oauth');
  }
  const accessToken = cookies.get('access_token');
  setHeaders({
    Authorization: 'bearer ' + accessToken
  });

  let serverID = params.server_id
  const responseChannel = await fetch('/api/channels/' + serverID);
  const channels: ApiResponse<Channel> = await responseChannel.json();
  let selected_server: Server;


  servers.data.forEach((val) => {
    if (serverID == "@me" && val.name == "@me") {
      selected_server = val
    }
    if (serverID != "@me" && serverID == val.id) {
      selected_server = val
    }
  })


  const serverMember = await fetch('/api/servers/member/' + selected_server!.id)
  const servermemberResult: ApiResponse<Servermember> = await serverMember.json()
  console.log(servermemberResult)

  return {
    servers: servers.data,
    channels: channels.data,
    selected_server: selected_server!,
    server_member: servermemberResult.data
  };
}
