import { createAction, INIT, props } from '@ngrx/store';
import type { AuthState } from './auth.types';
import { AuthActions } from './auth.types';

export const loginAction = createAction(
  AuthActions.Login,
  props<{
    rememberMe: boolean | undefined;
    token: {
      access: string;
      refresh: string;
    };
  }>()
);

export const rehydrateAuthState = createAction(INIT);

export const refreshAuthStateAction = createAction(
  AuthActions.Refresh,
  props<Omit<AuthState, ''>>()
);
