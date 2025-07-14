/**
 * HomePage component - Provides information about the GOTS Template
 */
export const HomePage = () => {
    return (
        <div className="w-full p-6 max-w-4xl mx-auto">
            {/* Page Header */}
            <div className="mb-8 text-center">
                <h1 className="text-4xl font-bold text-gray-900 mb-4">GOTS Template</h1>
                <p className="text-xl text-gray-600">A full-stack template for building modern web applications</p>
            </div>

            {/* Main Content */}
            <div className="space-y-8">
                {/* Overview Card */}
                <div className="bg-white rounded-lg shadow-lg p-6">
                    <h2 className="text-2xl font-semibold text-gray-900 mb-4">Overview</h2>
                    <p className="text-gray-700 mb-4">
                        The GOTS Template (Go, OpenAPI, TypeScript, SQLite) is a comprehensive full-stack application template designed to help you 
                        quickly bootstrap modern web applications. It combines a robust Go backend with a responsive React frontend.
                    </p>
                    <p className="text-gray-700">
                        This template includes authentication, modern API development, and a clean, modern user interface built with 
                        industry best practices and modern technologies.
                    </p>
                </div>

                {/* Technology Stack */}
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    {/* Backend */}
                    <div className="bg-white rounded-lg shadow-lg p-6">
                        <h3 className="text-xl font-semibold text-gray-900 mb-4 flex items-center">
                            <span className="text-2xl mr-2">⚙️</span>
                            Backend Technologies
                        </h3>
                        <ul className="space-y-2 text-gray-700">
                            <li className="flex items-center">
                                <span className="w-2 h-2 bg-blue-500 rounded-full mr-3"></span>
                                <strong>Go</strong> - High-performance backend language
                            </li>
                            <li className="flex items-center">
                                <span className="w-2 h-2 bg-blue-500 rounded-full mr-3"></span>
                                <strong>GORM</strong> - Go ORM for database interactions
                            </li>
                            <li className="flex items-center">
                                <span className="w-2 h-2 bg-blue-500 rounded-full mr-3"></span>
                                <strong>OpenAPI</strong> - API specification and code generation
                            </li>
                            <li className="flex items-center">
                                <span className="w-2 h-2 bg-blue-500 rounded-full mr-3"></span>
                                <strong>net/http</strong> - Standard library HTTP server
                            </li>
                        </ul>
                    </div>

                    {/* Frontend */}
                    <div className="bg-white rounded-lg shadow-lg p-6">
                        <h3 className="text-xl font-semibold text-gray-900 mb-4 flex items-center">
                            <span className="text-2xl mr-2">🎨</span>
                            Frontend Technologies
                        </h3>
                        <ul className="space-y-2 text-gray-700">
                            <li className="flex items-center">
                                <span className="w-2 h-2 bg-green-500 rounded-full mr-3"></span>
                                <strong>React</strong> - Modern UI library
                            </li>
                            <li className="flex items-center">
                                <span className="w-2 h-2 bg-green-500 rounded-full mr-3"></span>
                                <strong>TypeScript</strong> - Type-safe JavaScript
                            </li>
                            <li className="flex items-center">
                                <span className="w-2 h-2 bg-green-500 rounded-full mr-3"></span>
                                <strong>Vite</strong> - Fast build tool and dev server
                            </li>
                            <li className="flex items-center">
                                <span className="w-2 h-2 bg-green-500 rounded-full mr-3"></span>
                                <strong>Tailwind CSS</strong> - Utility-first CSS framework
                            </li>
                            <li className="flex items-center">
                                <span className="w-2 h-2 bg-green-500 rounded-full mr-3"></span>
                                <strong>DaisyUI</strong> - Component library for Tailwind
                            </li>
                        </ul>
                    </div>
                </div>

                {/* Features */}
                <div className="bg-white rounded-lg shadow-lg p-6">
                    <h3 className="text-xl font-semibold text-gray-900 mb-4 flex items-center">
                        <span className="text-2xl mr-2">✨</span>
                        Key Features
                    </h3>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div className="space-y-2">
                            <div className="flex items-center text-gray-700">
                                <span className="text-green-500 mr-2">✓</span>
                                Authentication system
                            </div>
                            <div className="flex items-center text-gray-700">
                                <span className="text-green-500 mr-2">✓</span>
                                RESTful API with OpenAPI specification
                            </div>
                            <div className="flex items-center text-gray-700">
                                <span className="text-green-500 mr-2">✓</span>
                                Database integration with GORM
                            </div>
                            <div className="flex items-center text-gray-700">
                                <span className="text-green-500 mr-2">✓</span>
                                Responsive design
                            </div>
                        </div>
                        <div className="space-y-2">
                            <div className="flex items-center text-gray-700">
                                <span className="text-green-500 mr-2">✓</span>
                                Modern UI components
                            </div>
                            <div className="flex items-center text-gray-700">
                                <span className="text-green-500 mr-2">✓</span>
                                Type-safe development
                            </div>
                            <div className="flex items-center text-gray-700">
                                <span className="text-green-500 mr-2">✓</span>
                                Hot reload development
                            </div>
                            <div className="flex items-center text-gray-700">
                                <span className="text-green-500 mr-2">✓</span>
                                Production-ready build system
                            </div>
                        </div>
                    </div>
                </div>

                {/* Getting Started */}
                <div className="bg-gradient-to-r from-blue-50 to-indigo-50 rounded-lg p-6 border border-blue-200">
                    <h3 className="text-xl font-semibold text-gray-900 mb-4 flex items-center">
                        <span className="text-2xl mr-2">🚀</span>
                        Getting Started
                    </h3>
                    <div className="space-y-3 text-gray-700">
                        <p>
                            <strong>1. Backend:</strong> Run <code className="bg-gray-200 px-2 py-1 rounded text-sm">make run</code> to start the Go server
                        </p>
                        <p>
                            <strong>2. Frontend:</strong> Navigate to the webapp directory and run <code className="bg-gray-200 px-2 py-1 rounded text-sm">npm run dev</code>
                        </p>
                        <p>
                            <strong>3. API:</strong> Use <code className="bg-gray-200 px-2 py-1 rounded text-sm">make api-codegen</code> to regenerate API code after OpenAPI changes
                        </p>
                        <p>
                            <strong>4. Tests:</strong> Run <code className="bg-gray-200 px-2 py-1 rounded text-sm">make test</code> to execute the test suite
                        </p>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default HomePage;
