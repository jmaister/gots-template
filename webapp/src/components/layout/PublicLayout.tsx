import { ReactNode, useState, useEffect } from 'react';
import { useLocation } from 'react-router-dom';

interface PublicLayoutProps {
  children: ReactNode;
}

// Simple public header without authentication
const PublicHeader = ({ pageTitle }: { pageTitle: string }) => {
  return (
    <div className="navbar bg-base-100 shadow-lg">
      <div className="navbar-start">
        <h1 className="text-xl font-semibold ml-2">{pageTitle}</h1>
      </div>
      <div className="navbar-end">
        <a
          href="/_/login"
          className="btn btn-primary btn-sm"
        >
          Login
        </a>
      </div>
    </div>
  );
};

// Simple public sidebar without authentication
const PublicSidebar = () => {
  return (
    <div className="drawer-side z-40">
      <label 
        htmlFor="sidebar-drawer" 
        aria-label="close sidebar" 
        className="drawer-overlay lg:hidden"
      ></label>
      
      <aside className="min-h-full bg-base-200 flex flex-col w-64">
        {/* Header */}
        <div className="navbar bg-base-300 min-h-16">
          <div className="navbar-start">
            <h1 className="text-sm font-bold">GOTS Template</h1>
          </div>
        </div>

        {/* Navigation Menu */}
        <div className="flex-1 p-4">
          <ul className="menu menu-vertical w-full">
            <li>
              <a href="/home">
                <span className="text-lg">🏠</span>
                <span>Home</span>
              </a>
            </li>
          </ul>
        </div>

        {/* Login Footer */}
        <div className="p-4 border-t border-base-300">
          <a
            href="/_/login"
            className="btn btn-primary btn-sm w-full"
          >
            Login
          </a>
        </div>
      </aside>
    </div>
  );
};

// Helper to derive a title from the path
const getPageTitleFromPath = (path: string): string => {
  if (path.startsWith('/home')) return 'Home';
  return 'Home'; // Default title
};

const PublicLayout = ({ children }: PublicLayoutProps) => {
  const [isMobileSidebarOpen, setIsMobileSidebarOpen] = useState(false);
  const location = useLocation();
  const [pageTitle, setPageTitle] = useState(getPageTitleFromPath(location.pathname));

  useEffect(() => {
    setPageTitle(getPageTitleFromPath(location.pathname));
  }, [location.pathname]);

  const toggleMobileSidebar = () => {
    setIsMobileSidebarOpen(!isMobileSidebarOpen);
  };

  return (
    <div className="drawer lg:drawer-open">
      <input 
        id="sidebar-drawer" 
        type="checkbox" 
        className="drawer-toggle" 
        checked={isMobileSidebarOpen}
        onChange={toggleMobileSidebar}
      />
      
      {/* Page content */}
      <div className="drawer-content flex flex-col">
        <PublicHeader pageTitle={pageTitle} />
        
        {/* Main content area */}
        <main className="flex-1 p-4 bg-base-100">
          {children}
        </main>
      </div>

      {/* Sidebar */}
      <PublicSidebar />
    </div>
  );
};

export default PublicLayout;
