import { ThemeProvider } from '@/components/ThemeContext';
import { Outlet } from 'react-router';

function RootLayout() {
  return (
    <ThemeProvider>
      <Outlet />
    </ThemeProvider>
  );
}

export default RootLayout;
