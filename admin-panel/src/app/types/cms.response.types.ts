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