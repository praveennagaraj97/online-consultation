import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { map } from 'rxjs';
import { doctorRoutes, sharedRoutes } from 'src/app/api-routes/routes';
import { PaginatedBaseAPiResponse } from 'src/app/types/api.response.types';
import { SpecialityEntity } from 'src/app/types/cms.response.types';
import { DoctorListResponse } from 'src/app/types/doctor.response.types';

@Injectable()
export class DoctorsListViewService {
  constructor(private http: HttpClient) {}

  getDoctorsList(perPage = 10, speciality: string) {
    const params: { [key: string]: string } = {};

    params['per_page'] = `${perPage}`;
    if (speciality) {
      params['speciality_id[eq]'] = speciality;
    }

    return this.http.get<DoctorListResponse>(doctorRoutes.DoctorsList, {
      params: params,
    });
  }

  // Get All Specialities with Key Set Pagination
  getSpecialities() {
    return this.http
      .get<PaginatedBaseAPiResponse<SpecialityEntity[]>>(
        sharedRoutes.Specialities,
        {
          params: { per_page: 50 },
        }
      )
      .pipe(
        map((res) => {
          return res.result?.map((res) => ({
            title: res.title,
            value: res.id,
          }));
        })
      );
  }
}
