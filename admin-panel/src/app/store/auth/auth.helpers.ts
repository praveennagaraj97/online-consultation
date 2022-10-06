import { authRoutes } from 'src/app/api-routes/routes';
import { APP_STORAGE_NAMES } from 'src/app/constants';
import { RefreshTokenResponse } from 'src/app/types/auth.response.types';
import {
  addDaysToDate,
  addMinutesToDate,
  subtractMinutes,
} from 'src/app/utils/date.utils';
import {
  _isMobileBrowser,
  _isSafari,
  _localStorage,
  _sessionStorage,
} from 'src/app/utils/web.api';
import { environment } from 'src/environments/environment';
import type { AuthState } from './auth.types';

// Store the auth session to storage
export function setAuthSession(
  rememberMe: boolean,
  access: string,
  refresh: string
) {
  const sessionData: AuthState = {
    loggedAt: new Date(),
    rememberMe,
    // If remember 30 days else 30 mins
    expiresAt: rememberMe ? addDaysToDate(29) : addMinutesToDate(29),
  };

  // Store in local Storage with expiry of 30 days
  if (rememberMe) {
    _localStorage()?.setItem(
      APP_STORAGE_NAMES.AUTH_STATE,
      JSON.stringify(sessionData)
    );
  } else {
    // Store in session storage and keep refreshing for every 30 mins
    _sessionStorage()?.setItem(
      APP_STORAGE_NAMES.AUTH_STATE,
      JSON.stringify(sessionData)
    );
  }

  // Save token to localStorage or session if safari or mobile browser
  if (_isMobileBrowser() || _isSafari()) {
    if (rememberMe) {
      _localStorage()?.setItem(APP_STORAGE_NAMES.AUTH_ACCESS_TOKEN, access);
      _localStorage()?.setItem(APP_STORAGE_NAMES.AUTH_REFRESH_TOKEN, refresh);
    } else {
      _sessionStorage()?.setItem(APP_STORAGE_NAMES.AUTH_ACCESS_TOKEN, access);
      _sessionStorage()?.setItem(APP_STORAGE_NAMES.AUTH_REFRESH_TOKEN, refresh);
    }
  }
}

let _cancelTimeOutId: any;
export function _nextAutoRefreshToken() {
  clearTimeout(_cancelTimeOutId);
  const session = getAuthSession();
  if (session) {
    const currentDate = new Date();
    const expireDate = subtractMinutes(1, new Date(session.expiresAt));

    const nextRefreshTime = +expireDate - +currentDate;

    console.log(
      'NEXT REFRESH TOKEN AFTER : ',
      nextRefreshTime / 1000 / 60,
      ' minutes'
    );
    if (nextRefreshTime >= 432000000) {
      return;
    }
    _cancelTimeOutId = setTimeout(async () => {
      try {
        await refreshAuthToken();
        _nextAutoRefreshToken();
      } catch (error) {
        //   logout();
      }
    }, nextRefreshTime);
  }
}

export function getAuthSession() {
  let session: string;
  session =
    _localStorage()?.getItem(APP_STORAGE_NAMES.AUTH_STATE) ||
    _sessionStorage()?.getItem(APP_STORAGE_NAMES.AUTH_STATE) ||
    '';

  return session ? (JSON.parse(session) as AuthState) : null;
}

// Refresh Auth Token
export async function refreshAuthToken() {
  try {
    const response = await fetch(
      environment.baseURL +
        authRoutes.RefreshToken +
        `?force=true&refresh_token=${
          _localStorage()?.getItem(APP_STORAGE_NAMES.AUTH_REFRESH_TOKEN) ||
          _sessionStorage()?.getItem(APP_STORAGE_NAMES.AUTH_REFRESH_TOKEN)
        }`,
      {
        method: 'GET',
        body: null,
        credentials: 'include',
      }
    );

    const data = (await response.json()) as RefreshTokenResponse;

    setAuthSession(false, data.access_token, data.refresh_token);
    console.log('refreshed');
    return data;
  } catch (error) {
    console.log(error);
    throw error;
  }
}
