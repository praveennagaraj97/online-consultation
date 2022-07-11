import {
  ImageType,
  PaginatedBaseAPiResponse,
  PhoneType,
} from './api.response.types';
import { HospitalEntity, SpokenLanguagesEntity } from './cms.response.types';

export interface DoctorsEntity {
  id: string;
  name: string;
  email: string;
  phone: PhoneType;
  professional_title: string;
  education: string;
  experience: number;
  profile_pic: ImageType | null;
  is_active: boolean;
  hospital: HospitalEntity;
  spoken_languages?: SpokenLanguagesEntity[] | null;
  next_available_slot: NextAvailableEntity | null;
}

type NextAvailableEntity = {
  id: string;
  date: string;
  start: string;
  end: string | null;
  is_available: boolean;
};

export interface DoctorListResponse
  extends PaginatedBaseAPiResponse<DoctorsEntity[]> {}
