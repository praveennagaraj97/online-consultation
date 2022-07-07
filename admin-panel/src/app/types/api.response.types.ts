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
