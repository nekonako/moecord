

export type ApiResponse<T> = {
  code: number;
  message: string;
  data: Array<T>;
};

export type Server = {
  id: string;
  name: string;
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
  created_at: string
};

export type Servermember = {
  id: string;
  user_id: string;
  server_id: string;
  username: string;
}


