import { FormArray, FormControl } from '@angular/forms';

export interface DoctorFormDTO {
  name: FormControl<string | null>;
  email: FormControl<string | null>;
  professional_title: FormControl<string | null>;
  phone_code: FormControl<string | null>;
  phone_number: FormControl<string | null>;
  education: FormControl<string | null>;
  experience: FormControl<string | null>;
  hospital_id: FormControl<string | null>;
  consultation_type_id: FormControl<string | null>;
  speciality_id: FormControl<string | null>;
  spoken_language_id: FormArray<FormControl<string>>;
  is_active: FormControl<string | null>;
}

export interface HospitalFormDTO {
  name: FormControl<string | null>;
  country: FormControl<string | null>;
  city: FormControl<string | null>;
  address: FormControl<string | null>;
}
