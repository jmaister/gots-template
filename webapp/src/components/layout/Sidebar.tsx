/**
 * Sidebar Component - Uses DaisyUI 5 components
 * 
 * DaisyUI Components used:
 * - drawer: For sidebar layout
 * - menu: For navigation items
 * - collapse: For accordion functionality (Users submenu)
 * - btn: For buttons
 * - card: For user info section
 * 
 * Reference: https://daisyui.com/llms.txt
 * Follow DaisyUI best practices:
 * - Use semantic component classes
 * - Combine with Tailwind utilities
 * - Use daisyUI color names (primary, secondary, etc.)
 */
import { useAuth } from '../../contexts/AuthContext';
import { Link } from 'react-router-dom';

interface SidebarProps {
  isOpenOnMobile: boolean; // Keep for consistency but not used internally
  toggleMobileSidebar: () => void;
  isCollapsed?: boolean;
  toggleCollapse?: () => void;
}

const Sidebar = ({
  toggleMobileSidebar,
  isCollapsed = false,
  toggleCollapse,
}: SidebarProps) => {
  const { currentUser, logout, isAdmin, isAuthenticated } = useAuth();

  const navItems = [
    { name: 'Home', icon: '🏠', path: '/home' },
    // Only show Profile link if user is authenticated
    ...(isAuthenticated ? [{ name: 'Profile', icon: '👤', path: '/profile' }] : []),
  ];

  return (
    <div className="drawer-side z-40">
      {/* Mobile overlay - for DaisyUI drawer functionality */}
      <label 
        htmlFor="sidebar-drawer" 
        aria-label="close sidebar" 
        className="drawer-overlay lg:hidden"
      ></label>
      
      {/* Sidebar - conditionally render based on collapse state */}
      {!isCollapsed && (
        <aside className="min-h-full bg-base-200 flex flex-col w-64 transition-all duration-300 ease-in-out">
          {/* Header */}
          <div className="navbar bg-base-300 min-h-16">
            <div className="navbar-start">
              <h1 className="text-sm font-bold">GOTS Template</h1>
            </div>
            <div className="navbar-end">
              {/* Desktop hide button */}
              {toggleCollapse && (
                <button
                  onClick={toggleCollapse}
                  className="btn btn-ghost btn-sm hidden lg:flex"
                  aria-label="Hide sidebar"
                >
                  <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
                  </svg>
                </button>
              )}
              {/* Mobile close button */}
              <button
                onClick={toggleMobileSidebar}
                className="btn btn-ghost btn-sm lg:hidden"
                aria-label="Close sidebar"
              >
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
          </div>

          {/* Navigation Menu */}
          <div className="flex-1 p-4">
            <ul className="menu menu-vertical w-full">
              {navItems.map((item) => (
                <li key={item.name}>
                  <Link to={item.path}>
                    <span className="text-lg">{item.icon}</span>
                    <span>{item.name}</span>
                  </Link>
                </li>
              ))}
              
              {/* Admin Dashboard - only show if user is admin */}
              {isAdmin && (
                <li>
                  <a href="/_/admin" target="_blank" rel="noopener noreferrer">
                    <span className="text-lg">⚙️</span>
                    <span>Admin Dashboard</span>
                  </a>
                </li>
              )}
            </ul>
          </div>

          {/* User Info Footer */}
          <div className="p-4 border-t border-base-300">
            {currentUser ? (
              <button
                onClick={logout}
                className="btn btn-outline btn-sm w-full"
              >
                Logout
              </button>
            ) : (
              <a
                href="/_/login"
                className="btn btn-primary btn-sm w-full"
              >
                Login
              </a>
            )}
          </div>
        </aside>
      )}
    </div>
  );
};

export default Sidebar;
