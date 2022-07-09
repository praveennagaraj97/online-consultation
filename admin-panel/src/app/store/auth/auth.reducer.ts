import { createReducer, on } from '@ngrx/store';
import { loginAction, refreshAuthStateAction } from './auth.actions';
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
  on(refreshAuthStateAction, (state, props) => ({ ...state }))
);
