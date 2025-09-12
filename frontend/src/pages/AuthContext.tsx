import { getCurrentUserData, getToken } from '@/api';
import { getLSToken, saveToken } from '@/lib/utils';
import type { UserResponse } from '@/types';

import React, { createContext, useEffect, useState } from 'react';
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

  async function fetchUserData() {
    const token = getLSToken();
    if (!token) {
      // TODO: handle no token better
      // Cleanup ls at least
      return;
    }
    const user = await getCurrentUserData();

    if (user) {
      setCurrentUser(user);
    }
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
    // TODO: use /user/me instead
    setCurrentUser(user);
    saveToken(tokenResponse);
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
