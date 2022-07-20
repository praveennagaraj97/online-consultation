import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { FormGroup } from '@angular/forms';
import { adminCMSRoutes } from 'src/app/api-routes/routes';
import { BaseAPiResponse } from 'src/app/types/api.response.types';
import { SpecialityEntity } from 'src/app/types/cms.response.types';
import { SpecialityFormDTO } from 'src/app/types/dto.types';

@Injectable({ providedIn: 'any' })
export class SpecialityFormService {
  constructor(private http: HttpClient) {}

  addNewSpeciality(fd: FormGroup<SpecialityFormDTO>, thumbnail: File) {
    const formData = new FormData();

    formData.append('title', fd.controls.title.value || '');
    formData.append('thumbnail', thumbnail);
    formData.append('thumbnail_height', '192');
    formData.append('thumbnail_width', '192');

    return this.http.post<BaseAPiResponse<SpecialityEntity>>(
      adminCMSRoutes.Speciality,
      formData
    );
  }
}
