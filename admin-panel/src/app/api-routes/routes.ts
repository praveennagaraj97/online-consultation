export const authRoutes = {
  Login: '/admin/login',
  ForgotPassword: '/admin/forgot_password',
  ResetPassword: '/admin/reset_password/',
  CheckTokenStatus: '/additional/jwt/status',
  RefreshToken: `/admin/refresh_token`,
};

export const adminRoutes = { ProfileDetails: '/admin/me' };

export const doctorRoutes = {
  DoctorsList: `/admin/doctor`,
  AddDoctor: `/admin/doctor`,
  UpdateDoctorStatus: (id: string, status: boolean) =>
    `/admin/doctor/status/${id}?status=${status}`,
};

export const adminCMSRoutes = {
  Hospital: `/admin/hospital`,
  Speciality: `/admin/speciality`,
  Language: `/admin/language`,
};

export const sharedRoutes = {
  Specialities: `/speciality`,
  ConsultationTypes: `/consultation_type`,
  Hospitals: `/hospital`,
  Languages: `/language`,
};

export const additionalRoutes = {
  GetCountries: `/additional/countries.json`,
};
