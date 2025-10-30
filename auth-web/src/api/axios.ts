import axios from "axios";
import type {AxiosInstance} from 'axios'

export const api: AxiosInstance = axios.create({
    baseURL: 'https://api.smarthome.hipahopa.ru',
    withCredentials: true,
    headers: {
        'Content-Type': 'application/json',
    },
})

api.interceptors.response.use(
    response => response,
    error => {
        console.error('API error:', error.response?.data || error.message)
        return Promise.reject(error)
    }
)
