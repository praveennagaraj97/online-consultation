import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Store } from '@ngrx/store';
import { tap } from 'rxjs';
import { authRoutes } from 'src/app/api-routes/routes';
import { loginAction } from 'src/app/store/auth/auth.actions';
import { AuthState } from 'src/app/store/auth/auth.types';
import { LoginResponse } from 'src/app/types/auth.response.types';
import { addDaysToDate, addMinutesToDate } from 'src/app/utils/date.utils';
import { validateEmail } from 'src/app/utils/validators';

@Injectable({ providedIn: 'any' })
export class LoginService {
  constructor(
    private store: Store<{ auth: AuthState }>,
    private http: HttpClient
  ) {}

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
            this.store.dispatch(
              loginAction({
                expiresAt: rememberMe
                  ? addDaysToDate(30).toISOString()
                  : addMinutesToDate(30).toISOString(),
                isLogged: true,
                rememberMe: rememberMe,
              })
            );
          }
        })
      );
  }
}
