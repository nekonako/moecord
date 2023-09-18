import { redirect } from '@sveltejs/kit';

type ApiResponse<T> = {
  code: number;
  message: string;
  data: Array<T>;
};

type Server = {
  id: string;
  name: string;
};

type Channel = {
  id: string;
  name: string;
  channel_type: string;
};

type Message = {
  id: string;
  channel_id: string;
  sender_id: string;
  content: string;
};

export async function load({ fetch, setHeaders, cookies }) {
  const responseServer = await fetch('/api/servers');
  const servers: ApiResponse<Server> = await responseServer.json();
  if (servers.code == 401) {
    throw redirect(307, '/oauth');
  }
  const accessToken = cookies.get('access_token');
  setHeaders({
    Authorization: 'bearer ' + accessToken
  });

  const firstServerID = servers.data[0].id;
  const responseChannel = await fetch('/api/channels/' + firstServerID);
  const channels: ApiResponse<Channel> = await responseChannel.json();
  console.log(channels)

  const firstChannelID = channels.data[0].id;
  const responseMessage = await fetch(`/api/messages/channels/${firstChannelID}`);
  const messages: ApiResponse<Message> = await responseMessage.json();

  return {
    servers: servers.data,
    channels: channels.data,
    messages: messages.data
  };
}
