import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { doctorRoutes } from 'src/app/api-routes/routes';
import { DoctorListResponse } from 'src/app/types/doctor.response.types';

@Injectable()
export class DoctorsListViewService {
  constructor(private http: HttpClient) {}

  getDoctorsList() {
    return this.http.get<DoctorListResponse>(doctorRoutes.DoctorsList, {
      params: { next_available_slot: true },
    });
  }
}
