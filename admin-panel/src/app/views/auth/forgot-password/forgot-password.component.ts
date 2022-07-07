import { transition, trigger, useAnimation } from '@angular/animations';
import { HttpClient } from '@angular/common/http';
import { Component } from '@angular/core';
import { fadeInTransformAnimation } from 'src/app/animations';
import { authRoutes } from 'src/app/api-routes/routes';
import { ErrorResponse } from 'src/app/types/api.response.types';

@Component({
  selector: 'app-forgot-password-view-component',
  templateUrl: 'forgot-password.component.html',
  animations: [
    trigger('fadeInOut', [
      transition('void<=>*', [useAnimation(fadeInTransformAnimation(600))]),
    ]),
  ],
})
export class ForgotPasswordViewComponent {
  email: string = '';

  isLoading = false;
  rspMsg: { type: 'error' | 'success'; message: string } | null = null;

  constructor(private http: HttpClient) {}

  handleRequestResetLink() {
    const formData = new FormData();
    formData.append('email', this.email);
    this.http.post(authRoutes.ForgotPassword, formData).subscribe({
      next: (res) => {
        console.log(res);
      },
      error: (err: ErrorResponse) => {
        console.log(err);
      },
    });
  }
}
