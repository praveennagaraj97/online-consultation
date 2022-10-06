import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Store } from '@ngrx/store';
import { tap } from 'rxjs';
import { authRoutes } from 'src/app/api-routes/routes';
import { APP_STORAGE_NAMES } from 'src/app/constants';
import { loginAction } from 'src/app/store/auth/auth.actions';
import { AuthState } from 'src/app/store/auth/auth.types';
import { LoginResponse } from 'src/app/types/auth.response.types';
import { validateEmail } from 'src/app/utils/validators';

@Injectable({ providedIn: 'any' })
export class LoginService {
  constructor(
    private store: Store<{ auth: AuthState }>,
    private http: HttpClient
  ) {}

  // Login user using email/username and password
  login(email: string, password: string, rememberMe: boolean) {
    const formData = new FormData();
    formData.append(validateEmail(email) ? 'email' : 'user_name', email);
    formData.append('password', password);

    return this.http
      .post<LoginResponse>(authRoutes.Login, formData, {
        params: { remember_me: rememberMe },
      })
      .pipe(
        tap((res) => {
          if (res) {
            this.setAuthTokenIfCookieDisabled(
              res.access_token,
              res.refresh_token,
              rememberMe
            );

            this.store.dispatch(
              loginAction({
                rememberMe,
                token: { access: res.access_token, refresh: res.refresh_token },
              })
            );
          }
        })
      );
  }

  private setAuthTokenIfCookieDisabled(
    access: string,
    refresh: string,
    rememberMe: boolean
  ) {
    if (!navigator.cookieEnabled) {
      if (rememberMe) {
        localStorage.setItem(APP_STORAGE_NAMES.AUTH_ACCESS_TOKEN, access);
      } else {
        sessionStorage.setItem(APP_STORAGE_NAMES.AUTH_REFRESH_TOKEN, refresh);
      }
    }
  }
}
