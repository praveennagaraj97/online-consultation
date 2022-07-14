export const authRoutes = {
  Login: '/admin/login',
  ForgotPassword: '/admin/forgot_password',
  ResetPassword: '/admin/reset_password/',
  CheckTokenStatus: '/additional/jwt/status',
};

export const adminRoutes = { ProfileDetails: '/admin/me' };

export const doctorRoutes = {
  DoctorsList: `/admin/doctor`,
  AddDoctor: `/admin/doctor`,
};

export const sharedRoutes = {
  Specialities: `/speciality`,
  ConsultationTypes: `/consultation_type`,
  Hospitals: `/hospital`,
  Languages: `/language`,
};

export const additionalRoutes = {
  GetCountryCodes: `/additional/countries.json`,
};
