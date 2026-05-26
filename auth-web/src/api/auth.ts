import {api} from './axios'

interface AuthPayload {
    email: string
    password: string
}

export const me = () => {
    return api.get('/auth/me')
}

export const login = (payload: AuthPayload) => {
    return api.post('/auth/login', payload)
}

export const logout = () => {
    return api.post('/auth/logout')
}

export const authorize = (queries: Record<string, string>) => {
    const formData = new URLSearchParams(queries)
    return api.post('/oauth/authorize', formData, {
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    })
}

export const register = (payload: AuthPayload) => {
    return api.post('/auth/register', payload)
}

export const requestPasswordReset = (email: string) => {
    return api.post('/auth/password-reset/request', { email })
}

export const confirmPasswordReset = (token: string, newPassword: string) => {
    return api.post('/auth/password-reset/confirm', { token, new_password: newPassword })
}