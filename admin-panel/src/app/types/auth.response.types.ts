import type { BaseAPiResponse } from './api.response.types';

export interface LoginResponse extends BaseAPiResponse<UserEntity> {
  access_token: string;
  refresh_token: string;
}

export interface UserEntity {
  id: string;
  name: string;
  user_name: string;
  email: string;
  joined_on: string;
}
