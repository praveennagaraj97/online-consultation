export type APiResponseStatus = { type: 'error' | 'success'; message: string };

export enum UserRoles {
  SuperAdmin = 'super_admin',
  Admin = 'admin',
  Editor = 'editor',
}

export type SelectOption = { title: string; value: string };

export type PaginateKeySetCache = {
  [pageNum: number]: string;
};

export type BreadcrumbPath = { path: string; title: string };

export interface Country {
  name: string;
  flag: string;
  code: string;
  dial_code: string;
}

export type ConsultationType = 'Instant' | 'Schedule';

export type ResponseMessageType = {
  message: string;
  type: 'error' | 'success' | 'info';
  timeOut?: number;
};

export type ConfirmPortalEventTypes = 'confirm' | 'cancel';
