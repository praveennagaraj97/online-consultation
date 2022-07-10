import { Injectable } from '@angular/core';
import {
  ActivatedRouteSnapshot,
  CanActivate,
  Router,
  RouterStateSnapshot,
} from '@angular/router';
import { Store } from '@ngrx/store';
import { map } from 'rxjs';
import { StoreState } from '../store/store.module';

@Injectable()
export class NotAuthorizedGuard implements CanActivate {
  constructor(private router: Router, private store: Store<StoreState>) {}

  // Get the Auth State from Local Storage Or Session Storage
  canActivate(_: ActivatedRouteSnapshot, { url }: RouterStateSnapshot) {
    return this.store.select('auth').pipe(
      map(({ isLogged }) => {
        if (isLogged) {
          return this.router.createUrlTree(['/dashboard'], {});
        }

        return true;
      })
    );
  }
}
