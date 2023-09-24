export type ApiResponse<T> = {
  code: number;
  message: string;
  data: T;
};

export type Server = {
  id: string;
  name: string;
  avatar: string
};

export type Channel = {
  category_id: string,
  category_name: string,
  channels: Array<{
    id: string;
    name: string;
    channel_type: string;
  }>
};

export type Message = {
  id: string;
  channel_id: string;
  sender_id: string;
  content: string;
  username: string;
  avatar: string;
  created_at: string
};

export type Servermember = {
  id: string;
  user_id: string;
  server_id: string;
  avatar: string;
  username: string;
}

export type Profile = {
  id: string;
  username: string;
  email: string;
  avatar: string;
}


export type WebsocketMessage<T> = {
  event_id: string;
  data: T;
};

export function getColor(asciiCode: number): string {
  let style = ' font-semibold'
  switch (true) {
    case asciiCode >= 65 && asciiCode <= 68:
      return 'text-success' + style;
    case asciiCode >= 69 && asciiCode <= 72:
      return 'text-error' + style;
    case asciiCode >= 73 && asciiCode <= 77:
      return 'text-accent' + style;
    case asciiCode >= 78 && asciiCode <= 81:
      return 'text-secondary' + style;
    case asciiCode >= 82 && asciiCode <= 85:
      return 'text-warning' + style;
    case asciiCode >= 86 && asciiCode <= 90:
      return 'text-error' + style
    default:
      return 'text-accent' + style;
  }
}
