export interface UserProfile {
  id: string;
  user_id: string;
  address: string;
  phone: string;
  avatar_url: string;
  description: string;
  created_at: string; 
  updated_at: string; 
}

export interface UserResponseBack {
  id: string;
  firebaseUID: string;
  email: string;
  name: string;
  rol: string;
  profile: UserProfile;
}
