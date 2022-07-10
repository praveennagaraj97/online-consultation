import { createAction, props } from '@ngrx/store';
import { AuthActions, AuthState } from './auth.types';

export const loginAction = createAction(AuthActions.Login, props<AuthState>());

export const rehydrateAuthState = createAction(AuthActions.RehydrateAuthState);

export const refreshAuthStateAction = createAction(
  AuthActions.Refresh,
  props<Omit<AuthState, 'isLogged'>>()
);
