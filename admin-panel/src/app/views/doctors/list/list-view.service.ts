import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { map } from 'rxjs';
import { doctorRoutes, sharedRoutes } from 'src/app/api-routes/routes';
import { PaginatedBaseAPiResponse } from 'src/app/types/api.response.types';
import {
  ConsultationEntity,
  SpecialityEntity,
} from 'src/app/types/cms.response.types';
import { DoctorListResponse } from 'src/app/types/doctor.response.types';

@Injectable()
export class DoctorsListViewService {
  constructor(private http: HttpClient) {}

  getDoctorsList(
    perPage = 10,
    speciality: string,
    consultationType: string,
    activeState: string,
    paginateId: string,
    search: string
  ) {
    const params: { [key: string]: string } = {};

    params['per_page'] = `${perPage}`;
    if (speciality) {
      params['speciality_id[eq]'] = speciality;
    }

    if (consultationType) {
      params['consultation_type_id[eq]'] = consultationType;
    }

    if (activeState) {
      params['is_active[eq]'] = activeState;
    }

    if (paginateId) {
      params['paginate_id'] = paginateId;
    }

    if (search) {
      params['name'] = search;
    }

    // Disable next available populate
    params['populate_next_available'] = 'false';

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
          const data =
            res.result?.map((res) => ({
              title: res.title,
              value: res.id,
            })) || [];

          return [{ title: 'All', value: '' }, ...data];
        })
      );
  }

  getConsultationTypes() {
    return this.http
      .get<PaginatedBaseAPiResponse<ConsultationEntity[]>>(
        sharedRoutes.ConsultationTypes
      )
      .pipe(
        map((res) => {
          const data =
            res.result?.map((res) => ({
              title: res.type,
              value: res.id,
            })) || [];

          return [{ title: 'All', value: '' }, ...data];
        })
      );
  }
}
