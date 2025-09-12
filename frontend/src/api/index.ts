import { up } from 'up-fetch';
import { getLSToken } from '../lib/utils';
import type { TokenResponse, UserResponse } from '../types';

if (!import.meta.env.VITE_API_BASE_URL) {
  throw new Error('API base url is not defined in env');
}

export const api = up(fetch, () => ({
  baseUrl: import.meta.env.VITE_API_BASE_URL,
  timeout: 30000,
  headers: { Authorization: getLSToken() },
}));

type CreateUserOpt = {
  name: string;
};

export async function createUser({ name }: CreateUserOpt) {
  try {
    const user = await api<UserResponse>('/user', {
      method: 'POST',
      body: { name },
    });

    return user;
  } catch (err) {
    console.error('Failed to create user', err);
  }
}
export async function getCurrentUserData() {
  try {
    const user = await api<UserResponse>('/user/me');

    return user;
  } catch (err) {
    console.error('Failed to fetch current user', err);
  }
}

export async function getToken({ name, id }: UserResponse) {
  try {
    const tokenResponse = await api<TokenResponse>('/token', {
      method: 'POST',
      body: { name, id },
    });

    return tokenResponse;
  } catch (err) {
    console.error('Failed to get user token', err);
  }
}
