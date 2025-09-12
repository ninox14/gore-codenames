import type { TokenResponse } from '@/types';
import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

const TOKEN_LS_KEY = 'bearer-token';
const TOKEN_EXPIRY_LS_KEY = 'bearer-token-expiry';

export function saveToken(tokenResponse?: TokenResponse) {
  if (!tokenResponse) {
    return;
  }

  localStorage.setItem(
    TOKEN_LS_KEY,
    `Bearer ${tokenResponse.AuthenticationToken}`
  );
  localStorage.setItem(
    TOKEN_EXPIRY_LS_KEY,
    tokenResponse.AuthenticationTokenExpiry
  );
}

export function getLSToken(): string | undefined {
  const token = localStorage.getItem(TOKEN_LS_KEY);
  if (!token) return;

  return token;
}
