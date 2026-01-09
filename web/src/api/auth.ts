import {api} from './axios'

interface CallbackPayload {
    code: string
    client_id: string
}

export const exchangeCode = (payload: CallbackPayload) => {
    return api.post('/panel/v1/users/exchange-code', payload)
}


export const me = () => {
    return api.get('/panel/v1/users/me')
}

export const logout = () => {
    return api.post('/panel/v1/users/logout')
}