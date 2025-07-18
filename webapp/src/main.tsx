import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import './index.css'; // Global styles including Tailwind and daisyUI
import App from './App'; // .tsx extension is usually omitted in imports
import { AuthProvider } from './contexts/AuthContext';

// Create a client
const queryClient = new QueryClient();

const rootElement = document.getElementById('root');

if (rootElement) {
  createRoot(rootElement).render(
    <StrictMode>
      <QueryClientProvider client={queryClient}>
        <AuthProvider>
          <App />
        </AuthProvider>
      </QueryClientProvider>
    </StrictMode>,
  );
} else {
  console.error("Failed to find the root element");
}
