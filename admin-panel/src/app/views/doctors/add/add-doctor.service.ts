import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { map } from 'rxjs';
import { sharedRoutes } from 'src/app/api-routes/routes';
import { PaginatedBaseAPiResponse } from 'src/app/types/api.response.types';
import { SelectOption } from 'src/app/types/app.types';
import {
  ConsultationEntity,
  HospitalEntity,
} from 'src/app/types/cms.response.types';

@Injectable()
export class AddDoctorService {
  private hospitals: HospitalEntity[] = [];

  constructor(private http: HttpClient) {}

  getHospitals(
    paginateId: string | null,
    search: string,
    shouldReset: boolean
  ) {
    let params: { [key: string]: string } = {};
    params['per_page'] = '50';
    if (paginateId) {
      params['paginate_id'] = paginateId;
    }

    if (search.trim().length) {
      params['name[search]'] = search;
    }
    if (shouldReset) {
      this.hospitals = [];
    }

    return this.http
      .get<PaginatedBaseAPiResponse<HospitalEntity[]>>(sharedRoutes.Hospitals, {
        params,
      })
      .pipe(
        map((res) => {
          if (res.result) {
            this.hospitals = [...this.hospitals, ...res.result];
          }

          const hospitals: SelectOption[] =
            this.hospitals?.map((hospital) => ({
              title: hospital.name,
              value: hospital.id,
            })) || [];

          return {
            hospitals: [
              { title: 'Add new hospital', value: 'add_new' },
              ...hospitals,
            ],
            nextId: res.paginate_id,
          };
        })
      );
  }

  getConsultaionTypeOptions() {
    return this.http
      .get<PaginatedBaseAPiResponse<ConsultationEntity[]>>(
        sharedRoutes.ConsultationTypes
      )
      .pipe(
        map((res) => {
          return res.result?.map((type) => ({
            title: type.title,
            value: type.id,
          }));
        })
      );
  }
}
