const API_URL = import.meta.env.VITE_API_URL

export default function Login() {
    const redirect = new URLSearchParams(window.location.search).get('redirect') ?? '/'
    const loginUrl = `${API_URL}/auth/login?redirect=${encodeURIComponent(redirect)}`

    return (
        <div className="flex h-screen items-center justify-center">
            <a
                href={loginUrl}
                className="rounded bg-blue-600 px-6 py-3 text-white hover:bg-blue-700"
            >
                Login with Google
            </a>
        </div>
    )
}
