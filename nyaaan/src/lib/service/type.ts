export type ApiResponse<T> = {
  code: number;
  message: string;
  data: T;
};

export type Server = {
  id: string;
  name: string;
  avatar: string;
};

export type CreateServeerRequest = {
  name: string,
}

export type Channel = {
  category_id: string;
  category_name: string;
  channels: Array<{
    id: string;
    name: string;
    channel_type: string;
  }>;
};

export type CreateChannelRequest = {
  name: string,
  server_id: string,
  category_id: string
  is_private: boolean,
  type: string
}

export type CreateChannelCategoryRequest = {
  name: string,
  is_private: boolean,
  server_id: string
};



export type Message = {
  id: string;
  channel_id: string;
  sender_id: string;
  content: string;
  username: string;
  avatar: string;
  created_at: string;
};

export type ServerMember = {
  id: string;
  user_id: string;
  server_id: string;
  avatar: string;
  online: boolean;
  username: string;
};

export type Profile = {
  id: string;
  username: string;
  email: string;
  avatar: string;
};

export type WebsocketMessage<T> = {
  event_id: string;
  data: T;
};

export type UserConnectionState = {
  user_id: string;
  status: string;
};

export type Fetch = (input: RequestInfo | URL, init?: RequestInit) => Promise<Response>;

