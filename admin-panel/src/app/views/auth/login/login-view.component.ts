import { transition, trigger, useAnimation } from '@angular/animations';
import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { Subscription } from 'rxjs';
import { fadeInTransformAnimation } from 'src/app/animations';
import { ErrorResponse, LoginErrors } from 'src/app/types/api.response.types';
import { clearSubscriptions } from 'src/app/utils/helpers';
import { validateEmail } from 'src/app/utils/validators';
import { LoginService } from './login-view.service';

@Component({
  selector: 'app-login-view',
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

  constructor(private loginService: LoginService, private router: Router) {}

  handleLogin() {
    this.isLoading = true;

    const formData = new FormData();
    formData.append(
      validateEmail(this.email) ? 'email' : 'user_name',
      this.email
    );
    formData.append('password', this.password);

    this.subs.push(
      this.loginService
        .login(this.email, this.password, this.rememberMe)

        .subscribe({
          next: (res) => {
            this.setMessage({ message: res.message, type: 'success' }, () => {
              this.isLoading = false;

              const parsedURL = this.router.parseUrl(this.router.url);
              this.router.navigate([
                parsedURL.queryParams?.['redirectTo'] || '/',
              ]);
            });
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
    msg: { type: 'error' | 'success'; message: string } | null,
    callback?: () => void
  ) {
    this.rspMsg = msg;
    setTimeout(() => {
      this.errors = null;
      this.rspMsg = null;
      if (callback) {
        callback();
      }
    }, 3000);
  }

  ngOnDestry() {
    clearSubscriptions(this.subs);
  }
}
