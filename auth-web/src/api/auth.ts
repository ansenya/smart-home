import {api} from './axios'

interface LoginPayload {
    email: string
    password: string
}

export const me = () => {
    return api.get('/auth/me')
}

export const login = (payload: LoginPayload) => {
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