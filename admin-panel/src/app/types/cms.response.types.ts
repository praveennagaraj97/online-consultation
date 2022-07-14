import { ImageType } from './api.response.types';
import { ConsultationType } from './app.types';

export interface HospitalEntity {
  id: string;
  name: string;
  city: string;
  country: string;
  address: string;
  location: LocationEntity;
}

export interface LocationEntity {
  type: string;
  coordinates?: number[] | null;
}

export interface SpokenLanguagesEntity {
  id: string;
  name: string;
  locale_name: string;
}

export interface SpecialityEntity {
  id: string;
  title: string;
  slug: string;
  thumbnail: ImageType | null;
}

export interface ConsultationEntity {
  id: string;
  title: string;
  icon: ImageType;
  description: string;
  price: number;
  action_name: string;
  type: ConsultationType;
}

export interface LanguageEntity {
  id: string;
  name: string;
  locale_name: string;
}
