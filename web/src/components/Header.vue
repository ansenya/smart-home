<script setup lang="ts">
import Button from "./Button.vue";

const CLIENT_ID = "c85e6304-7f65-49f9-8145-823bd71a5a83";
const PROVIDER_AUTHORIZE = "https://id.smarthome.hipahopa.ru";

function randStr(len = 43) {
  const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-._~";
  let s = "";
  const rnd = crypto.getRandomValues(new Uint8Array(len));
  for (let i = 0; i < len; i++) s += chars[rnd[i]! % chars.length];
  return s;
}

function base64UrlEncode(buf: ArrayBuffer) {
  const bytes = new Uint8Array(buf);
  let str = "";
  for (let i = 0; i < bytes.byteLength; i++) str += String.fromCharCode(bytes[i]!);
  return btoa(str).replace(/\+/g, "-").replace(/\//g, "_").replace(/=+$/, "");
}

async function sha256(input: string) {
  const enc = new TextEncoder();
  return await crypto.subtle.digest("SHA-256", enc.encode(input));
}

async function generateCodeChallenge(verifier: string) {
  const hash = await sha256(verifier);
  return base64UrlEncode(hash);
}

async function login() {
  const redirectUri = `${location.origin}/auth/callback`;
  const state = randStr(16);
  const codeVerifier = randStr(64);
  const codeChallenge = await generateCodeChallenge(codeVerifier);

  sessionStorage.setItem("oauth_state", state);
  sessionStorage.setItem("pkce_verifier", codeVerifier);

  const params = new URLSearchParams({
    client_id: CLIENT_ID,
    redirect_uri: redirectUri,
    response_type: "code",
    scope: "openid profile email",
    state,
    code_challenge: codeChallenge,
    code_challenge_method: "S256",
  });

  window.location.href = `${PROVIDER_AUTHORIZE}?${params.toString()}`;
}
</script>


<template>
  <header>
    <nav class="bg-black border-gray-200 px-4 lg:px-6 py-2.5 dark:bg-gray-800">
      <div class="flex flex-wrap justify-between items-center mx-auto max-w-screen-xl">
        <a href="/" class="flex items-center">
          <span class="self-center text-xl font-semibold whitespace-nowrap dark:text-white">Hiphome</span>
        </a>
        <div class="flex items-center lg:order-2">
          <Button
              :onclick="login"
          >
            Login
          </Button>
        </div>
        <div class="hidden justify-between items-center w-full lg:flex lg:w-auto lg:order-1" id="mobile-menu-2">
          <ul class="flex flex-col mt-4 font-medium lg:flex-row lg:space-x-8 lg:mt-0">
            <li>
              <a href="/home"
                 class="block py-2 pr-4 pl-3 text-white rounded bg-primary-700 lg:bg-transparent lg:text-primary-700 lg:p-0 dark:text-white"
                 aria-current="page">Home</a>
            </li>
          </ul>
        </div>
      </div>
    </nav>
  </header>
</template>
