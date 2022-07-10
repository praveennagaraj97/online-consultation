export type APiResponseStatus = { type: 'error' | 'success'; message: string };

export enum UserRoles {
  SuperAdmin = 'super_admin',
  Admin = 'admin',
  Editor = 'editor',
}
