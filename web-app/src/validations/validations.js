import {ref, watch} from "vue";

export function useEmailValidation(email) {
    const error = ref('')
    watch(email, (newVal) => {
        if (!newVal) {
            error.value = 'Введите email'
        } else if (!/^\S+@\S+\.\S+$/.test(newVal)) {
            error.value = 'Некорректный email'
        } else {
            error.value = ''
        }
    })
    return error
}
