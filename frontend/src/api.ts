import type {
    ApiError,
    AuthResponse,
    ForgotPasswordRequest,
    LoginRequest,
    RegisterRequest,
    ResendVerificationRequest,
    ResetPasswordRequest,
    User,
    VerifyEmailRequest,
} from '~/types';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

class ApiServiceError extends Error {
    constructor(public message: string, public status: number) {
        super(message);
        this.name = 'ApiServiceError';
    }
}

class ApiService {
    private baseURL: string;
    private csrfToken: string | null = null;

    constructor(baseURL: string) {
        this.baseURL = baseURL;
    }

    private async getCSRFToken(): Promise<string> {
        if (this.csrfToken) {
            return this.csrfToken;
        }

        try {
            const response = await fetch(`${this.baseURL}/csrf-token`, {
                method: 'GET',
                credentials: 'include',
            });

            if (response.ok) {
                const data = await response.json();
                this.csrfToken = data.csrf_token;
                return this.csrfToken!;
            }
        } catch (error) {
            console.error('Failed to get CSRF token:', error);
        }

        throw new ApiServiceError('Failed to get CSRF token', 500);
    }

    private async request<T>(
        endpoint: string,
        options: RequestInit = {}
    ): Promise<T> {
        const url = `${this.baseURL}${endpoint}`;

        const config: RequestInit = {
            headers: {
                'Content-Type': 'application/json',
                ...options.headers,
            },
            credentials: 'include',
            ...options,
        };

        if (options.method && options.method !== 'GET') {
            try {
                const csrfToken = await this.getCSRFToken();
                config.headers = {
                    ...config.headers,
                    'X-CSRF-Token': csrfToken,
                };
            } catch (error) {
                console.error('Failed to add CSRF token:', error);
            }
        }

        try {
            const response = await fetch(url, config);

            if (response.status === 401) {
                const refreshed = await this.tryRefreshToken();
                if (refreshed) {
                    if (options.method && options.method !== 'GET') {
                        const csrfToken = await this.getCSRFToken();
                        config.headers = {
                            ...config.headers,
                            'X-CSRF-Token': csrfToken,
                        };
                    }
                    const retryResponse = await fetch(url, config);
                    if (retryResponse.ok) {
                        const contentType = retryResponse.headers.get('content-type');
                        if (contentType && contentType.includes('application/json')) {
                            return retryResponse.json();
                        }
                        return retryResponse.text() as T;
                    }
                }
                throw new ApiServiceError('Session expired', 401);
            }

            if (response.status === 403) {
                this.csrfToken = null;
                if (options.method && options.method !== 'GET') {
                    try {
                        const csrfToken = await this.getCSRFToken();
                        config.headers = {
                            ...config.headers,
                            'X-CSRF-Token': csrfToken,
                        };
                        const retryResponse = await fetch(url, config);
                        if (retryResponse.ok) {
                            const contentType = retryResponse.headers.get('content-type');
                            if (contentType && contentType.includes('application/json')) {
                                return retryResponse.json();
                            }
                            return retryResponse.text() as T;
                        }
                    } catch (retryError) {
                        // If retry fails, continue with original error handling
                    }
                }
            }

            if (!response.ok) {
                const errorData = await response
                    .json()
                    .catch(() => ({ error: 'Request failed' }));
                const error = errorData as ApiError;
                throw new ApiServiceError(
                    error.error || 'Request failed',
                    response.status
                );
            }

            const contentType = response.headers.get('content-type');
            if (contentType && contentType.includes('application/json')) {
                return response.json();
            }
            return response.text() as T;
        } catch (error) {
            if (error instanceof ApiServiceError) {
                throw error;
            }
            throw new ApiServiceError('Network error', 0);
        }
    }

    private async tryRefreshToken(): Promise<boolean> {
        try {
            const csrfToken = await this.getCSRFToken();

            const response = await fetch(`${this.baseURL}/auth/refresh`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'X-CSRF-Token': csrfToken,
                },
                credentials: 'include',
            });

            return response.ok;
        } catch (error) {
            console.error('Token refresh failed:', error);
            return false;
        }
    }

    async login(credentials: LoginRequest): Promise<AuthResponse> {
        return this.request<AuthResponse>('/auth/login', {
            method: 'POST',
            body: JSON.stringify(credentials),
        });
    }

    async register(userData: RegisterRequest): Promise<AuthResponse> {
        return this.request<AuthResponse>('/auth/register', {
            method: 'POST',
            body: JSON.stringify(userData),
        });
    }

    async verifyEmail(data: VerifyEmailRequest): Promise<{ message: string }> {
        return this.request<{ message: string }>('/auth/verify-email', {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async resendVerification(
        data: ResendVerificationRequest
    ): Promise<{ message: string }> {
        return this.request<{ message: string }>('/auth/resend-verification', {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async forgotPassword(
        data: ForgotPasswordRequest
    ): Promise<{ message: string }> {
        return this.request<{ message: string }>('/auth/forgot-password', {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async resetPassword(
        data: ResetPasswordRequest
    ): Promise<{ message: string }> {
        return this.request<{ message: string }>('/auth/reset-password', {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async validateResetToken(token: string): Promise<{ valid: boolean }> {
        return this.request<{ valid: boolean }>(
            `/auth/validate-reset-token?token=${token}`
        );
    }

    async logout(): Promise<void> {
        try {
            await this.request('/auth/logout', {
                method: 'POST',
            });
        } catch (error) {
            console.warn('Server logout failed:', error);
        }
    }

    async logoutAll(): Promise<void> {
        try {
            await this.request('/auth/logout-all', {
                method: 'POST',
            });
        } catch (error) {
            console.warn('Server logout failed:', error);
        }
    }

    async getGoogleAuthUrl(): Promise<{ auth_url: string }> {
        return this.request<{ auth_url: string }>('/auth/google');
    }

    async getProfile(): Promise<User> {
        return this.request<User>('/auth/profile');
    }

    async getServices(): Promise<any> {
        return this.request<any>('/api/services');
    }

    async healthCheck(): Promise<{ status: string }> {
        return this.request<{ status: string }>('/health');
    }
}

export const apiService = new ApiService(API_BASE_URL);
export { ApiServiceError };
