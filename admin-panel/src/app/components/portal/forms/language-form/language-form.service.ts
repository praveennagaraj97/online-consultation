import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { FormGroup } from '@angular/forms';
import { adminCMSRoutes } from 'src/app/api-routes/routes';
import { BaseAPiResponse } from 'src/app/types/api.response.types';
import { LanguageEntity } from 'src/app/types/cms.response.types';
import { LanguageFormDTO } from 'src/app/types/dto.types';

@Injectable({ providedIn: 'any' })
export class LanguageFormService {
  constructor(private http: HttpClient) {}

  addNewLanguage(fd: FormGroup<LanguageFormDTO>) {
    const formData = new FormData();

    formData.append('name', fd.controls.name.value || '');
    formData.append('locale_name', fd.controls.locale_name.value || '');

    return this.http.post<BaseAPiResponse<LanguageEntity>>(
      adminCMSRoutes.Language,
      formData
    );
  }
}
