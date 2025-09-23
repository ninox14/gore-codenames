import { createRoot } from 'react-dom/client';
import './index.css';
import App from './App.tsx';
import { Toaster } from '@/components/ui/sonner';
import { BrowserRouter } from 'react-router';

createRoot(document.getElementById('root')!).render(
  // <StrictMode>
  <>
    <BrowserRouter>
      <App />
    </BrowserRouter>
    <Toaster richColors />
  </>
  // </StrictMode>
);
