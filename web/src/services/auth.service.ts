import { api } from '@/api/client'

const CLIENT_ID = 'c85e6304-7f65-49f9-8145-823bd71a5a83'
const PROVIDER_AUTHORIZE = 'https://id.smarthome.hipahopa.ru'
const REDIRECT_URI = `${window.location.origin}/auth/callback`

interface User {
  ID: string
  name: string
  email?: string
}

function randStr(len = 43) {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-._~'
  let s = ''
  const rnd = crypto.getRandomValues(new Uint8Array(len))
  for (let i = 0; i < len; i++) s += chars[rnd[i]! % chars.length]
  return s
}

function base64UrlEncode(buf: ArrayBuffer) {
  const bytes = new Uint8Array(buf)
  let str = ''
  for (let i = 0; i < bytes.byteLength; i++) str += String.fromCharCode(bytes[i]!)
  return btoa(str).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '')
}

async function sha256(input: string) {
  const enc = new TextEncoder()
  return await crypto.subtle.digest('SHA-256', enc.encode(input))
}

async function generateCodeChallenge(verifier: string) {
  const hash = await sha256(verifier)
  return base64UrlEncode(hash)
}

export const authService = {
  async initiateLogin() {
    const state = randStr(16)
    const codeVerifier = randStr(64)

    sessionStorage.setItem('oauth_state', state)
    sessionStorage.setItem('pkce_verifier', codeVerifier)

    const params = new URLSearchParams({
      client_id: CLIENT_ID,
      redirect_uri: REDIRECT_URI,
      response_type: 'code',
      scope: 'openid profile email',
      state,
      code_challenge: await generateCodeChallenge(codeVerifier),
      code_challenge_method: 'S256',
    })

    window.location.href = `${PROVIDER_AUTHORIZE}?${params.toString()}`
  },

  async handleCallback(code: string, state: string) {
    const savedState = sessionStorage.getItem('oauth_state')
    const codeVerifier = sessionStorage.getItem('pkce_verifier')

    if (!savedState || savedState !== state) {
      throw new Error('Invalid state parameter')
    }
    if (!codeVerifier) {
      throw new Error('Missing code verifier')
    }

    sessionStorage.removeItem('oauth_state')
    sessionStorage.removeItem('pkce_verifier')

    await api.post('https://api.smarthome.hipahopa.ru/panel/v1/users/exchange-code', {
      code: code,
      client_id: CLIENT_ID,
      redirect_uri: REDIRECT_URI,
    })
  },

  async getMe(): Promise<User> {
    const { data } = await api.get('https://api.smarthome.hipahopa.ru/panel/v1/users/me')
    return data
  },

  async logout() {
    try {
      await api.post('https://api.smarthome.hipahopa.ru/panel/v1/users/logout')
    } catch (e) {
      console.warn('Logout error', e)
    } finally {
      sessionStorage.clear()
    }
  },
}
