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
  const responseChannel = await fetch('/api/servers/' + serverID + "/member/invite");
  const result = await responseChannel.json()
  if (result.code == 200) {
    throw redirect(307, "/channel/" + serverID)
  }
}
