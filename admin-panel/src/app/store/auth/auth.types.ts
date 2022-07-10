export enum AuthActions {
  Login = '[LoginViewComponent] Login',
  Refresh = '[RefreshTokenResolver] Refresh',
  RehydrateAuthState = '[AppComponent] RehydrateAuthState',
}

export interface AuthState {
  rememberMe: boolean;
  expiresAt: string;
  isLogged: boolean;
}
