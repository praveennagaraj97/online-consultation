export const doctorFormErrors = {
  name: {
    required: "Doctor's name cannot be empty",
  },
  email: {
    required: "Doctor's email cannot be empty",
    email: 'Entered email is not valid',
  },
  professional_title: {
    required: "Doctor's professional title cannot be empty",
  },
  phone_number: {
    required: "Doctor's phone number cannot be empty",
    pattern: 'Provided phone number is not valid',
  },
  education: {
    required: "Doctor's education cannot be empty",
  },
  hospital_id: {
    required: "Doctor's work place is required",
  },
  experience: {
    required: "Doctor's experience cannot be empty",
    pattern: "Doctor's experience must be number",
  },
  consultation_type_id: {
    required: "Doctor's Consultation type cannot be empty",
  },
  speciality_id: {
    required: 'Speciality is required for Scheduled consultation',
  },
  spoken_language_id: {
    required: 'Select atleast one language which doctor speaks',
  },
};
