import {ref} from 'vue'

export const baseUrl = 'https://api.smarthome.hipahopa.ru'

export function useAuth() {
    const errors = ref({
        email: '',
        password: '',
        confirmedPassword: '',
        form: ''
    })

    const resetErrors = () => {
        errors.value = {
            email: '',
            password: '',
            confirmedPassword: '',
            form: ''
        }
    }

    return {baseUrl, errors, resetErrors}
}

export function useAuthValidation() {
    const validateEmail = (email) => {
        if (!email) return 'Введите email'
        if (!/^\S+@\S+\.\S+$/.test(email)) return 'Некорректный email'
        return ''
    }

    const validatePassword = (password) => {
        if (!password) return 'Введите пароль'
        if (password.length < 8) return 'Минимум 8 символов'
        return ''
    }

    return {validateEmail, validatePassword}
}