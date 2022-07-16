import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { FormGroup } from '@angular/forms';
import { map } from 'rxjs';
import { additionalRoutes, adminCMSRoutes } from 'src/app/api-routes/routes';
import { BaseAPiResponse } from 'src/app/types/api.response.types';
import { Country } from 'src/app/types/app.types';
import { HospitalEntity } from 'src/app/types/cms.response.types';
import { HospitalFormDTO } from 'src/app/types/dto.types';

@Injectable({ providedIn: 'any' })
export class HospitalFormService {
  constructor(private http: HttpClient) {}

  getCountries() {
    return this.http.get<Country[]>(additionalRoutes.GetCountries).pipe(
      map((country) => {
        return country.map(({ name }) => ({ title: name, value: name }));
      })
    );
  }

  addNewHospital(fd: FormGroup<HospitalFormDTO>) {
    const formKeys = Object.keys(fd.value);

    const formData = new FormData();

    formKeys.forEach((formKey) => {
      formData.append(formKey, fd.get(formKey)?.value);
    });

    return this.http.post<BaseAPiResponse<HospitalEntity>>(
      adminCMSRoutes.Hospital,
      formData
    );
  }
}
