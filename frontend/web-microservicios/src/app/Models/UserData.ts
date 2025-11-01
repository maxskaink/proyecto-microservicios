import { UserProfile } from "firebase/auth";

export interface UserData {
  id: string;
  firebaseUID: string;  
  email: string;
  name: string;
  rol: 'admin' | 'cliente' | 'vendedor' | string; 
  profile: UserProfile;
}