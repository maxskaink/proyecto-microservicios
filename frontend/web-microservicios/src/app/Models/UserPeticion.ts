export interface UserProfilePeticion {
  address: string;
  phone: string;
  avatar_url: string;
  description: string;
}

export interface userPeticion {
  email: string;
  name: string;
  profile: UserProfilePeticion;
}