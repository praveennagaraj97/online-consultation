import { createReducer, on } from '@ngrx/store';
import { APP_STORAGE_NAMES } from 'src/app/constants';
import {
  loginAction,
  refreshAuthStateAction,
  rehydrateAuthState,
} from './auth.actions';
import { AuthState } from './auth.types';

const initialState: AuthState = {
  expiresAt: '',
  isLogged: false,
  rememberMe: false,
};

export const authReducer = createReducer(
  initialState,
  on(loginAction, (_, props) => {
    return {
      expiresAt: props.expiresAt,
      isLogged: props.isLogged,
      rememberMe: props.rememberMe,
    };
  }),
  on(rehydrateAuthState, () => {
    let data: string;
    data = localStorage.getItem(APP_STORAGE_NAMES.AUTH_STATE) || '';

    if (!data) {
      data = sessionStorage.getItem(APP_STORAGE_NAMES.AUTH_STATE) || '';
    }

    if (data) {
      const authState = JSON.parse(data) as AuthState;
      return authState;
    }

    return initialState;
  }),
  on(refreshAuthStateAction, (state, props) => ({ ...state }))
);
