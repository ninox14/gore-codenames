import { getCurrentUserData, getToken } from '@/api';
import { getLSToken, removeToken, saveToken } from '@/lib/utils';
import type { UserResponse } from '@/types';

import React, { createContext, useEffect, useState } from 'react';
import { useNavigate } from 'react-router';
import { toast } from 'sonner';

type UserState = UserResponse | null;

type CurrentUserContextType = {
  onSuccessfullUserCreate: (
    data: UserResponse,
    onSuccess?: () => void
  ) => Promise<void>;
  user: UserState;
};
export const AuthContext = createContext<CurrentUserContextType>({
  user: null,
  onSuccessfullUserCreate: async () => {},
});

export function AuthContextProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  const [currentUser, setCurrentUser] = useState<UserState>(null);
  const navigate = useNavigate();
  async function fetchUserData() {
    const token = getLSToken();
    if (!token) {
      // TODO: handle no token better
      // Cleanup ls at least
      return;
    }
    const user = await getCurrentUserData();

    if (!user) {
      removeToken();
      navigate('/');
      return;
    }

    setCurrentUser(user);
  }

  async function onSuccessfullUserCreate(
    user: UserResponse,
    onSuccess?: () => void
  ) {
    const tokenResponse = await getToken(user);
    // TODO: think about redirect to somewhere?
    if (!tokenResponse) {
      toast.error('Could not get token');
      return;
    }
    saveToken(tokenResponse);
    const userData = await getCurrentUserData();
    if (!userData) {
      removeToken();
      navigate('/');
      return;
    }

    setCurrentUser(user);
    onSuccess?.();
  }

  useEffect(() => {
    fetchUserData();
  }, []);

  return (
    <AuthContext.Provider
      value={{ user: currentUser, onSuccessfullUserCreate }}
    >
      {children}
    </AuthContext.Provider>
  );
}
