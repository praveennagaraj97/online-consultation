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
  selector: 'app-forgot-password-view',
  templateUrl: 'forgot-password-view.component.html',
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

  // State
  isLoading = false;
  rspMsg: { type: 'error' | 'success'; message: string } | null = null;

  // Properties
  showSentView = false;
  resetSuccessMsgs = {
    title: 'Check your email',
    message: 'We have sent a password recover link to your email address',
  };

  constructor(private http: HttpClient) {}

  handleRequestResetLink() {
    const formData = new FormData();
    formData.append('email', this.email);

    this.isLoading = true;
    this.subs.push(
      this.http
        .post<BaseAPiResponse<null>>(authRoutes.ForgotPassword, formData)
        .subscribe({
          next: () => {
            this.isLoading = false;
            this.showSentView = true;
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
