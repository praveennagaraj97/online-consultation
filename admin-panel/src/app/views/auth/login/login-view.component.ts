import { transition, trigger, useAnimation } from '@angular/animations';
import { HttpClient } from '@angular/common/http';
import { Component } from '@angular/core';
import { Subscription } from 'rxjs';
import { fadeInTransformAnimation } from 'src/app/animations';
import { authRoutes } from 'src/app/api-routes/routes';
import type {
  ErrorResponse,
  LoginErrors,
} from 'src/app/types/api.response.types';
import { LoginResponse } from 'src/app/types/auth.response.types';
import { clearSubscriptions } from 'src/app/utils/helpers';
import { validateEmail } from 'src/app/utils/validators';

@Component({
  selector: 'app-login-view-component',
  templateUrl: 'login-view.component.html',
  animations: [
    trigger('fadeInOut', [
      transition('void <=> *', [useAnimation(fadeInTransformAnimation(600))]),
    ]),
    trigger('fadeIn', [
      transition('void => *', [useAnimation(fadeInTransformAnimation(600))]),
    ]),
  ],
})
export class LoginViewComponent {
  // Subscriptions
  private subs: Subscription[] = [];

  // Form Data
  email: string = '';
  password: string = '';
  rememberMe: boolean = true;

  // APi Call State
  isLoading = false;
  errors: {
    email: string;
    password: string;
    user_name: string;
  } | null = null;
  rspMsg: { type: 'error' | 'success'; message: string } | null = null;

  constructor(private http: HttpClient) {}

  handleLogin() {
    this.isLoading = true;

    const formData = new FormData();
    formData.append(
      validateEmail(this.email) ? 'email' : 'user_name',
      this.email
    );
    formData.append('password', this.password);

    this.subs.push(
      this.http
        .post<LoginResponse>(authRoutes.Login, formData, {
          params: { remember_me: this.rememberMe },
        })
        .subscribe({
          next: (res) => {
            this.isLoading = false;
            this.setMessage({ message: res.message, type: 'success' });
          },
          error: (err: ErrorResponse<LoginErrors>) => {
            this.isLoading = false;
            this.errors = { ...err.error.errors };
            this.setMessage({ message: err.error.message, type: 'error' });
          },
        })
    );
  }

  private setMessage(
    msg: { type: 'error' | 'success'; message: string } | null
  ) {
    this.rspMsg = msg;
    setTimeout(() => {
      this.errors = null;
      this.rspMsg = null;
    }, 3000);
  }

  ngOnDestry() {
    clearSubscriptions(this.subs);
  }
}
