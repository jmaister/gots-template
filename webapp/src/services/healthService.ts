import { useQuery } from '@tanstack/react-query';
import { healthCheck } from '../apiclient/sdk.gen';
import { client } from '../apiclient/client.gen';

/**
 * Configure the API client for health service
 */
const configureHealthClient = () => {
    client.setConfig({
        baseUrl: '/',
    });
};

/**
 * Custom hook for fetching health status using TanStack Query
 */
export const useHealthStatus = () => {
    // Configure the client
    configureHealthClient();
    
    return useQuery({
        queryKey: ['health'],
        queryFn: async () => {
            const response = await healthCheck();
            return response.data;
        },
        refetchInterval: 30000, // Refetch every 30 seconds
        staleTime: 10000, // Consider data stale after 10 seconds
        retry: 3, // Retry failed requests 3 times
        retryDelay: attemptIndex => Math.min(1000 * 2 ** attemptIndex, 30000), // Exponential backoff
    });
};

/**
 * Health status query options for manual usage
 */
export const healthStatusQueryOptions = {
    queryKey: ['health'],
    queryFn: async () => {
        configureHealthClient();
        const response = await healthCheck();
        return response.data;
    },
    refetchInterval: 30000,
    staleTime: 10000,
    retry: 3,
    retryDelay: (attemptIndex: number) => Math.min(1000 * 2 ** attemptIndex, 30000),
};
