import { ReactNode, useState, useEffect } from 'react';
import Sidebar from './Sidebar';
import Header from './Header';
import { useLocation } from 'react-router-dom';

interface MainLayoutProps {
  children: ReactNode;
}

// Helper to derive a title from the path
const getPageTitleFromPath = (path: string): string => {
  if (path.startsWith('/profile')) return 'Profile Settings';
  if (path.startsWith('/home')) return 'Home';
  return 'Home'; // Default title
};

const MainLayout = ({ children }: MainLayoutProps) => {
  const [isMobileSidebarOpen, setIsMobileSidebarOpen] = useState(false);
  const [isCollapsed, setIsCollapsed] = useState(false);
  const location = useLocation();
  const [pageTitle, setPageTitle] = useState(getPageTitleFromPath(location.pathname));

  useEffect(() => {
    setPageTitle(getPageTitleFromPath(location.pathname));
  }, [location.pathname]);

  const toggleMobileSidebar = () => {
    setIsMobileSidebarOpen(!isMobileSidebarOpen);
  };

  const toggleCollapse = () => {
    setIsCollapsed(!isCollapsed);
  };

  return (
    <div className={`drawer ${isCollapsed ? '' : 'lg:drawer-open'}`}>
      <input 
        id="sidebar-drawer" 
        type="checkbox" 
        className="drawer-toggle" 
        checked={isMobileSidebarOpen}
        onChange={toggleMobileSidebar}
      />
      
      {/* Page content */}
      <div className="drawer-content flex flex-col">
        <Header
          pageTitle={pageTitle}
          toggleMobileSidebar={toggleMobileSidebar}
          isCollapsed={isCollapsed}
          toggleCollapse={toggleCollapse}
        />
        <main className="flex-1 p-6 bg-base-100 overflow-x-auto">
          {children}
        </main>
      </div>

      {/* Sidebar */}
      <Sidebar
        isOpenOnMobile={isMobileSidebarOpen}
        toggleMobileSidebar={toggleMobileSidebar}
        isCollapsed={isCollapsed}
        toggleCollapse={toggleCollapse}
      />
    </div>
  );
};

export default MainLayout;
