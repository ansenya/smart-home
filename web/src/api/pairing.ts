import {api} from './axios'

interface PairingPayload {
    code: string
}

export const startPairing = () => {
    return api.post('/devices/pairing/start')
}

export const pairingStatus = (payload: PairingPayload) => {
    return api.post('/devices/pairing/status', null, {
        params: payload
    })
}