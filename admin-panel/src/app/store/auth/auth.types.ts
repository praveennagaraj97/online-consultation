export enum AuthActions {
  Login = '[LoginViewComponent] Login',
  Refresh = '[RefreshTokenResolver] Refresh',
}

export interface AuthState {
  rememberMe: boolean;
  expiresAt: string;
  isLogged: boolean;
}
