import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { map } from 'rxjs';
import { adminRoutes } from 'src/app/api-routes/routes';
import { ProfileResponse } from 'src/app/types/auth.response.types';

@Injectable({ providedIn: 'any' })
export class UserOptionsService {
  constructor(private http: HttpClient) {}

  getProfileDetails() {
    return this.http
      .get<ProfileResponse>(adminRoutes.ProfileDetails)
      .pipe(map((res) => res.result));
  }
}
