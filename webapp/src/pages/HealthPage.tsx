import { useHealthStatus } from '../services';

/**
 * HealthPage component that displays the application health status
 */
export const HealthPage = () => {
    const {
        data: healthData,
        isLoading: loading,
        error,
        refetch,
        dataUpdatedAt
    } = useHealthStatus();

    const lastChecked = new Date(dataUpdatedAt);

    const fetchHealthStatus = () => {
        refetch();
    };

    const getStatusColor = (status: string) => {
        switch (status.toLowerCase()) {
            case 'ok':
            case 'healthy':
                return 'text-green-600 bg-green-100';
            case 'warning':
                return 'text-yellow-600 bg-yellow-100';
            case 'error':
            case 'unhealthy':
                return 'text-red-600 bg-red-100';
            default:
                return 'text-gray-600 bg-gray-100';
        }
    };

    const getStatusIcon = (status: string) => {
        switch (status.toLowerCase()) {
            case 'ok':
            case 'healthy':
                return '✅';
            case 'warning':
                return '⚠️';
            case 'error':
            case 'unhealthy':
                return '❌';
            default:
                return 'ℹ️';
        }
    };

    const formatUptime = (uptime: string) => {
        // If uptime is in a duration format like "1h30m45s", return as is
        // Otherwise, try to parse and format it
        if (uptime.match(/^\d+[hms]/)) {
            return uptime;
        }
        
        // Try to parse as seconds and convert to human readable
        const seconds = parseInt(uptime);
        if (!isNaN(seconds)) {
            const hours = Math.floor(seconds / 3600);
            const minutes = Math.floor((seconds % 3600) / 60);
            const secs = seconds % 60;
            
            if (hours > 0) {
                return `${hours}h ${minutes}m ${secs}s`;
            } else if (minutes > 0) {
                return `${minutes}m ${secs}s`;
            } else {
                return `${secs}s`;
            }
        }
        
        return uptime;
    };

    return (
        <div className="w-full p-6 max-w-4xl mx-auto">
            {/* Page Header */}
            <div className="mb-8">
                <div className="flex items-center justify-between">
                    <div>
                        <h1 className="text-4xl font-bold text-gray-900 mb-2">System Health</h1>
                        <p className="text-lg text-gray-600">Monitor the application health status</p>
                    </div>
                    <button
                        onClick={fetchHealthStatus}
                        disabled={loading}
                        className="bg-blue-500 hover:bg-blue-700 disabled:bg-blue-300 text-white font-bold py-2 px-4 rounded transition-colors"
                    >
                        {loading ? '🔄 Checking...' : '🔄 Refresh'}
                    </button>
                </div>
            </div>

            {/* Error State */}
            {error && (
                <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-6">
                    <div className="flex items-center">
                        <span className="text-xl mr-2">❌</span>
                        <div>
                            <strong className="font-bold">Error!</strong>
                            <span className="block sm:inline"> {error instanceof Error ? error.message : 'Failed to fetch health status'}</span>
                        </div>
                    </div>
                </div>
            )}

            {/* Loading State */}
            {loading && !healthData && (
                <div className="bg-white rounded-lg shadow-lg p-8 text-center">
                    <div className="animate-spin w-8 h-8 border-4 border-blue-500 border-t-transparent rounded-full mx-auto mb-4"></div>
                    <p className="text-gray-600">Checking system health...</p>
                </div>
            )}

            {/* Health Status Display */}
            {healthData && (
                <div className="space-y-6">
                    {/* Main Status Card */}
                    <div className="bg-white rounded-lg shadow-lg p-6">
                        <div className="flex items-center justify-between mb-4">
                            <h2 className="text-2xl font-semibold text-gray-900">Current Status</h2>
                            <div className={`px-3 py-1 rounded-full text-sm font-medium ${getStatusColor(healthData.status)}`}>
                                <span className="mr-1">{getStatusIcon(healthData.status)}</span>
                                {healthData.status.toUpperCase()}
                            </div>
                        </div>
                        
                        <div className="text-sm text-gray-500">
                            Last checked: {lastChecked.toLocaleString()}
                        </div>
                    </div>

                    {/* System Information */}
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                        {/* Timestamp */}
                        <div className="bg-white rounded-lg shadow-lg p-6">
                            <h3 className="text-lg font-semibold text-gray-900 mb-2 flex items-center">
                                <span className="text-xl mr-2">🕐</span>
                                Server Time
                            </h3>
                            <p className="text-gray-700">
                                {new Date(healthData.timestamp).toLocaleString()}
                            </p>
                        </div>

                        {/* Version */}
                        {healthData.version && (
                            <div className="bg-white rounded-lg shadow-lg p-6">
                                <h3 className="text-lg font-semibold text-gray-900 mb-2 flex items-center">
                                    <span className="text-xl mr-2">📦</span>
                                    Version
                                </h3>
                                <p className="text-gray-700 font-mono">
                                    {healthData.version}
                                </p>
                            </div>
                        )}

                        {/* Uptime */}
                        {healthData.uptime && (
                            <div className="bg-white rounded-lg shadow-lg p-6">
                                <h3 className="text-lg font-semibold text-gray-900 mb-2 flex items-center">
                                    <span className="text-xl mr-2">⏱️</span>
                                    Uptime
                                </h3>
                                <p className="text-gray-700 font-mono">
                                    {formatUptime(healthData.uptime)}
                                </p>
                            </div>
                        )}
                    </div>

                    {/* Raw Response (for debugging) */}
                    <div className="bg-white rounded-lg shadow-lg p-6">
                        <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
                            <span className="text-xl mr-2">🔍</span>
                            Raw Health Response
                        </h3>
                        <pre className="bg-gray-100 p-4 rounded text-sm overflow-x-auto text-gray-700">
                            {JSON.stringify(healthData, null, 2)}
                        </pre>
                    </div>
                </div>
            )}

            {/* Info Section */}
            <div className="mt-8 bg-blue-50 border border-blue-200 rounded-lg p-6">
                <h3 className="text-lg font-semibold text-blue-900 mb-2 flex items-center">
                    <span className="text-xl mr-2">ℹ️</span>
                    About Health Monitoring
                </h3>
                <p className="text-blue-800">
                    This page monitors the backend API health endpoint to ensure the system is running properly. 
                    The health check includes server status, uptime information, and version details.
                </p>
            </div>
        </div>
    );
};
