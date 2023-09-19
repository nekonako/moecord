import { redirect } from '@sveltejs/kit';
import type { ApiResponse, Channel, Message, Server } from './type';

export async function load({ fetch, setHeaders, cookies, params }) {
  const responseServer = await fetch('/api/servers');
  const servers: ApiResponse<Server> = await responseServer.json();
  if (servers.code == 401) {
    throw redirect(307, '/oauth');
  }

  console.log(params.server_id, params.channel_id)

  const accessToken = cookies.get('access_token');
  setHeaders({
    Authorization: 'bearer ' + accessToken
  });

  const serverID = params.server_id;
  const responseChannel = await fetch('/api/channels/' + serverID);
  const channels: ApiResponse<Channel> = await responseChannel.json();
  console.log(channels)

  const firstChannelID = channels.data[0].channels[0].id;
  const responseMessage = await fetch(`/api/messages/channels/${firstChannelID}`);
  const messages: ApiResponse<Message> = await responseMessage.json();
  let selected_server: Server;

  servers.data.forEach((val) => {
    if (serverID == "@me" && val.name == "@me") {
      selected_server = val
    }
    if (serverID != "@me" && serverID == val.id) {
      selected_server = val
    }
  })

  return {
    servers: servers.data,
    channels: channels.data,
    messages: messages.data,
    selected_server: selected_server!
  };
}
