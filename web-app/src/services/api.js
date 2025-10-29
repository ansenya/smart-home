import {baseUrl} from '@/composables/useAuth'

export const authService = {
    async login(payload, params) {
        const url = params
            ? `${baseUrl}/auth/authorize?${new URLSearchParams(params)}`
            : `${baseUrl}/auth/authorize`

        return fetch(url, {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify(payload),
            credentials: 'include'
        })
    },

    async register(payload) {
        return fetch(`${baseUrl}/auth/register`, {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify(payload)
        })
    },

    async resetPasswordRequest(payload) {
        return fetch(`${baseUrl}/auth/reset-password`, {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify(payload)
        })
    },

    async resetPasswordConfirm(payload) {
        return fetch(`${baseUrl}/auth/reset-password/change-password`, {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify(payload)
        })
    }
}