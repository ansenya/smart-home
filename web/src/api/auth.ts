import {api} from './axios'

interface CallbackPayload {
    code: string
    client_id: string
}

export const exchangeCode = (payload: CallbackPayload) => {
    return api.post('/users/exchange-code', payload)
}


export const me = () => {
    return api.get('/users/me')
}
