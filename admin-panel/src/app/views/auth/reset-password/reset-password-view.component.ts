import { transition, trigger, useAnimation } from '@angular/animations';
import { HttpClient } from '@angular/common/http';
import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { Subscription } from 'rxjs';
import { fadeInTransformAnimation } from 'src/app/animations';
import { authRoutes } from 'src/app/api-routes/routes';
import {
  BaseAPiResponse,
  ErrorResponse,
} from 'src/app/types/api.response.types';
import { APiResponseStatus } from 'src/app/types/app.types';
import { clearSubscriptions } from 'src/app/utils/helpers';

@Component({
  selector: 'app-reset-password-view',
  templateUrl: 'reset-password-view.component.html',
  animations: [
    trigger('fadeInOut', [
      transition('void<=>*', [useAnimation(fadeInTransformAnimation())]),
    ]),
    trigger('fadeIn', [
      transition('void<=>*', [useAnimation(fadeInTransformAnimation())]),
    ]),
  ],
})
export class ResetPasswordViewComponent {
  // Subs
  private subs: Subscription[] = [];

  // Form Data
  password = '';
  confirmPassword = '';

  // State
  isLoading = false;
  rspMsg: APiResponseStatus | null = null;

  errors: {
    password: string;
    confirmPassword: string;
  } | null = null;

  constructor(private http: HttpClient, private router: Router) {}

  handleSubmit() {
    this.isLoading = true;

    const parsedURL = this.router.parseUrl(this.router.url);
    const token = parsedURL.queryParams?.['verifyCode'];
    if (!token) {
      this.setMessage({ message: 'Something went wrong', type: 'error' });
    }

    const formData = new FormData();

    formData.append('password', this.password);
    formData.append('confirm_password', this.confirmPassword);

    this.subs.push(
      this.http
        .post<BaseAPiResponse<null>>(authRoutes.ResetPassword + token, formData)
        .subscribe({
          next: (res) => {
            this.setMessage({ message: res.message, type: 'success' }, () => {
              this.isLoading = false;
              this.router.navigateByUrl('/auth/login');
            });
          },
          error: (err: ErrorResponse<null>) => {
            this.setMessage({ message: err.error.message, type: 'error' });
            this.isLoading = false;
          },
        })
    );
  }

  private setMessage(rsp: APiResponseStatus | null, callback?: () => void) {
    this.rspMsg = rsp;

    setTimeout(() => {
      this.rspMsg = null;
      if (callback) {
        callback();
      }
    }, 3000);
  }

  ngOnDestroy() {
    clearSubscriptions(this.subs);
  }
}
