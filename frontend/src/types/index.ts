export type UserResponse = { name: string; id: string; created_at: string };

export type TokenResponse = {
  AuthenticationToken: string;
  AuthenticationTokenExpiry: string;
};
