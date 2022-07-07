import { transition, trigger, useAnimation } from '@angular/animations';
import { HttpClient } from '@angular/common/http';
import { Component } from '@angular/core';
import { Subscription } from 'rxjs';
import { fadeInTransformAnimation } from 'src/app/animations';
import { authRoutes } from 'src/app/api-routes/routes';
import {
  BaseAPiResponse,
  ErrorResponse,
} from 'src/app/types/api.response.types';
import { clearSubscriptions } from 'src/app/utils/helpers';

@Component({
  selector: 'app-forgot-password-view-component',
  templateUrl: 'forgot-password.component.html',
  animations: [
    trigger('fadeInOut', [
      transition('void<=>*', [useAnimation(fadeInTransformAnimation())]),
    ]),
    trigger('fadeIn', [
      transition('void=>*', [useAnimation(fadeInTransformAnimation())]),
    ]),
  ],
})
export class ForgotPasswordViewComponent {
  // Subscription
  private subs: Subscription[] = [];

  // Template Form Data
  email: string = '';
  isLoading = false;
  rspMsg: { type: 'error' | 'success'; message: string } | null = null;

  constructor(private http: HttpClient) {}

  handleRequestResetLink() {
    const formData = new FormData();
    formData.append('email', this.email);

    this.isLoading = true;
    this.subs.push(
      this.http
        .post<BaseAPiResponse<null>>(authRoutes.ForgotPassword, formData)
        .subscribe({
          next: (res) => {
            this.isLoading = false;
            this.setMessage({
              message: res.message,
              type: 'error',
            });
          },
          error: (err: ErrorResponse<null>) => {
            this.isLoading = false;
            this.setMessage({
              message: err.error.message,
              type: 'error',
            });
          },
        })
    );
  }

  private cancelId: any;
  private setMessage(
    msg: { type: 'error' | 'success'; message: string } | null
  ) {
    clearTimeout(this.cancelId);
    this.rspMsg = msg;

    this.cancelId = setTimeout(() => {
      this.rspMsg = null;
    }, 3000);
  }

  ngOnDestroy() {
    clearSubscriptions(this.subs);
  }
}
