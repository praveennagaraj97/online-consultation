import { transition, trigger, useAnimation } from '@angular/animations';
import { HttpClient } from '@angular/common/http';
import { Component } from '@angular/core';
import { fadeInTransformAnimation } from 'src/app/animations';
import { authRoutes } from 'src/app/api-routes/routes';
import type {
  ErrorResponse,
  LoginErrors,
} from 'src/app/types/api.response.types';
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
  errorMessage = '';

  constructor(private http: HttpClient) {}

  handleLogin() {
    const formData = new FormData();

    formData.append(
      validateEmail(this.email) ? 'email' : 'user_name',
      this.email
    );
    formData.append('password', this.password);

    this.isLoading = true;
    this.http.post(authRoutes.Login, formData).subscribe({
      next: (res) => {
        console.log(res);
      },
      error: (err: ErrorResponse<LoginErrors>) => {
        this.errors = { ...err.error.errors };
        this.errorMessage = err.error.message;
        this.isLoading = false;
        setTimeout(() => {
          this.errors = null;
          this.errorMessage = '';
        }, 3000);
      },
    });

    console.log(this.email, this.password);
  }
}
