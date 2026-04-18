const API_URL = import.meta.env.VITE_API_URL

export default function Login() {
    return (
        <div className="flex h-screen items-center justify-center">
            <a
                href={`${API_URL}/auth/login`}
                className="rounded bg-blue-600 px-6 py-3 text-white hover:bg-blue-700"
            >
                Login with Google
            </a>
        </div>
    )
}
