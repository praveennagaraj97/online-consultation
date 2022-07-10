import type { BaseAPiResponse } from './api.response.types';
import { UserRoles } from './app.types';

export interface LoginResponse extends BaseAPiResponse<UserEntity> {
  access_token: string;
  refresh_token: string;
}

export interface UserEntity {
  id: string;
  name: string;
  user_name: string;
  email: string;
  role: UserRoles;
  joined_on: string;
}

export type JWTTokenStatus = {
  expires: string | null;
  is_valid: boolean;
};

export interface ProfileResponse extends BaseAPiResponse<UserEntity> {}
