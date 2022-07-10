import { Injectable } from '@angular/core';
import {
  ActivatedRouteSnapshot,
  CanActivate,
  CanActivateChild,
  Router,
  RouterStateSnapshot,
  UrlTree,
} from '@angular/router';
import { Store } from '@ngrx/store';
import { map, Observable } from 'rxjs';
import { StoreState } from '../store/store.module';

@Injectable()
export class AuthorizedGuard implements CanActivate, CanActivateChild {
  constructor(private router: Router, private store: Store<StoreState>) {}

  // Get the Auth State from Local Storage Or Session Storage
  canActivate(_: ActivatedRouteSnapshot, { url }: RouterStateSnapshot) {
    return this.store.select('auth').pipe(
      map(({ isLogged }) => {
        if (!isLogged) {
          return this.router.createUrlTree(['/auth/login'], {
            queryParams: { redirectTo: url || '/dashboard' },
          });
        }

        return isLogged;
      })
    );
  }

  canActivateChild(
    _: ActivatedRouteSnapshot,
    { url }: RouterStateSnapshot
  ):
    | boolean
    | UrlTree
    | Observable<boolean | UrlTree>
    | Promise<boolean | UrlTree> {
    return this.store.select('auth').pipe(
      map(({ isLogged }) => {
        if (!isLogged) {
          return this.router.createUrlTree(['/auth/login'], {
            queryParams: { redirectTo: url || '/dashboard' },
          });
        }

        return isLogged;
      })
    );
  }
}
