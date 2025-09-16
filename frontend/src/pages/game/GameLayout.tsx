import AuthDialog from '@/components/AuthDialog';
import { useContext } from 'react';
import { Outlet } from 'react-router';
import { AuthContext } from '../AuthContext';
import { getLSToken } from '@/lib/utils';

function GameLayout() {
  const { user } = useContext(AuthContext);
  const token = getLSToken();

  console.log('User', user);
  return (
    <div className="min-h-screen bg-black py-12 px-4 sm:px-6 lg:px-8">
      <AuthDialog open={!user && !token} />
      <Outlet />
    </div>
  );
}

export default GameLayout;
