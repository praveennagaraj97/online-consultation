// Shared API Response or basic api interface

import type { HttpErrorResponse } from '@angular/common/http';

export interface ErrorResponse<T = any> extends HttpErrorResponse {
  error: {
    errors: T;
    message: string;
    status_code: number;
  };
}

export interface LoginErrors {
  email: string;
  password: string;
  user_name: string;
}

export interface BaseAPiResponse<T> {
  result: T;
  status_code: number;
  message: string;
}

export interface PaginatedBaseAPiResponse<T> {
  count: number;
  next: boolean;
  prev: boolean;
  paginate_id: string | null;
  result: T | null;
  status_code: number;
  message: string;
}

export type PhoneType = {
  code: string;
  number: string;
};

export type ImageType = {
  image_src: string;
  blur_data_url: string;
  width: number;
  height: number;
};
