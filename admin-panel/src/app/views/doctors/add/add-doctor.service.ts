import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { FormGroup } from '@angular/forms';
import { map } from 'rxjs';
import { doctorRoutes, sharedRoutes } from 'src/app/api-routes/routes';
import {
  BaseAPiResponse,
  PaginatedBaseAPiResponse,
} from 'src/app/types/api.response.types';
import { SelectOption } from 'src/app/types/app.types';
import {
  ConsultationEntity,
  HospitalEntity,
  LanguageEntity,
  SpecialityEntity,
} from 'src/app/types/cms.response.types';
import { DoctorsEntity } from 'src/app/types/doctor.response.types';
import { DoctorFormDTO } from 'src/app/types/dto.types';

@Injectable()
export class AddDoctorService {
  private hospitals: HospitalEntity[] = [];
  private specialities: SpecialityEntity[] = [];
  private languages: LanguageEntity[] = [];

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

    return this.http
      .get<PaginatedBaseAPiResponse<HospitalEntity[]>>(sharedRoutes.Hospitals, {
        params,
      })
      .pipe(
        map((res) => {
          if (shouldReset) {
            this.hospitals = [];
          }
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
            title: type.type,
            value: type.id,
          }));
        })
      );
  }

  getSpecialitiesOptions(
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
      params['title[search]'] = search;
    }

    return this.http
      .get<PaginatedBaseAPiResponse<SpecialityEntity[]>>(
        sharedRoutes.Specialities,
        { params }
      )
      .pipe(
        map((res) => {
          if (shouldReset) {
            this.specialities = [];
          }
          if (res.result) {
            this.specialities = [...this.specialities, ...res.result];
          }

          const specialities: SelectOption[] =
            this.specialities?.map((speciality) => ({
              title: speciality.title,
              value: speciality.id,
            })) || [];

          return {
            specialities: [
              { title: 'Add new speciality', value: 'add_new' },
              ...specialities,
            ],
            nextId: res.paginate_id,
          };
        })
      );
  }

  getLangaugesOptions(
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

    return this.http
      .get<PaginatedBaseAPiResponse<LanguageEntity[]>>(sharedRoutes.Languages, {
        params,
      })
      .pipe(
        map((res) => {
          if (shouldReset) {
            this.languages = [];
          }

          if (res.result) {
            this.languages = [...this.languages, ...res.result];
          }

          const langauges: SelectOption[] =
            this.languages?.map((speciality) => ({
              title: speciality.name,
              value: speciality.id,
            })) || [];

          return {
            langauges: [
              { title: 'Add new language', value: 'add_new' },
              ...langauges,
            ],
            nextId: res.paginate_id,
          };
        })
      );
  }

  handleAddDoctor(form: FormGroup<DoctorFormDTO>, profilePic: File | null) {
    const formData = new FormData();

    const keys = Object.keys(form.value);

    keys.forEach((key) => {
      const value = form.get(key)?.value;

      if (value) {
        if (key == 'spoken_language_id') {
          value.forEach((val: string) => {
            formData.append('spoken_language_id', val);
          });
        } else {
          formData.append(key, value);
        }
      }
    });
    if (profilePic) {
      formData.append('profile_pic', profilePic);
      formData.append('profile_pic_width', '160');
      formData.append('profile_pic_height', '160');
    }

    return this.http.post<BaseAPiResponse<DoctorsEntity>>(
      doctorRoutes.AddDoctor,
      formData
    );
  }
}
