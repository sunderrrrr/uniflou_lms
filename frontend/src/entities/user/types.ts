export interface UserEntity {
  id: string;
  email: string;
  role: 'student' | 'parent' | 'tutor' | 'admin';
}
