import { createReducer, on } from '@ngrx/store';
import { addDaysToDate } from 'src/app/utils/date.utils';
import {
  loginAction,
  refreshAuthStateAction,
  rehydrateAuthState,
} from './auth.actions';
import {
  getAuthSession,
  setAuthSession,
  _nextAutoRefreshToken,
} from './auth.helpers';

const initialState = {
  isLogged: false,
};

export const authReducer = createReducer(
  initialState,
  on(loginAction, (_, props) => {
    setAuthSession(
      props.rememberMe || false,
      props.token.access,
      props.token.refresh
    );
    _nextAutoRefreshToken();

    return {
      isLogged: true,
    };
  }),
  on(rehydrateAuthState, (state) => {
    const session = getAuthSession();
    if (session) {
      // Logout if login was 29 days before. | Refresh Token Expired
      if (addDaysToDate(29, new Date(session.loggedAt)) < new Date()) {
        return { ...state, isLogged: false };
      }

      _nextAutoRefreshToken();

      return { ...state, isLogged: true };
    }

    return { ...state, ...initialState };
  }),
  on(refreshAuthStateAction, (state, props) => ({ ...state }))
);
