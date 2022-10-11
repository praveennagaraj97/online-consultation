import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { doctorRoutes } from 'src/app/api-routes/routes';

@Injectable({ providedIn: 'any' })
export class EditDoctorService {
  constructor(private http: HttpClient) {}

  updateDoctorStatus(id: string, status: boolean) {
    return this.http.patch(doctorRoutes.UpdateDoctorStatus(id, status), null);
  }
}
