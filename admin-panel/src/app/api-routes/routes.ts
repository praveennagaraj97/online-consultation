export const authRoutes = {
  Login: '/admin/login',
  ForgotPassword: '/admin/forgot_password',
  ResetPassword: '/admin/reset_password/',
  CheckTokenStatus: '/additional/jwt/status',
};

export const adminRoutes = { ProfileDetails: '/admin/me' };

export const doctorRoutes = {
  DoctorsList: `/admin/doctor`,
};

export const sharedRoutes = {
  Specialities: `/speciality`,
  ConsultationTypes: `/consultation/`,
  Hospitals: `/hospital`,
};

export const additionalRoutes = {
  GetCountryCodes: `/additional/countries.json`,
};
