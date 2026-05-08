export function isPublicHost() {
    const publicHosts = (import.meta.env.VITE_PUBLIC_HOSTS ?? '').split(',').map(h => h.trim()).filter(Boolean)
    const hostname = window.location.hostname
    return publicHosts.some(host => {
        if (host.includes('/')) {
            const [base, bits] = host.split('/')
            const mask = ~(0xffffffff >>> parseInt(bits))
            const toInt = ip => ip.split('.').reduce((acc, b) => (acc << 8) | parseInt(b), 0)
            return (toInt(hostname) & mask) === (toInt(base) & mask)
        }
        return hostname === host
    })
}

export const API_URL = isPublicHost() ? `http://${window.location.hostname}:8080` : import.meta.env.VITE_API_URL

function authHeaders(extra = {}) {
    return { Authorization: `Bearer ${localStorage.getItem('token')}`, ...extra }
}

function handleUnauthorized(res) {
    if (res.status === 401) {
        localStorage.removeItem('token')
        window.location.href = '/login'
    }
    return res
}

export function mediaUrl(path) {
    const token = localStorage.getItem('token')
    return `${API_URL}${path}?token=${token}`
}

export const api = {
    get(path) {
        return fetch(`${API_URL}${path}`, {
            headers: authHeaders(),
        }).then(handleUnauthorized)
    },

    post(path, body) {
        return fetch(`${API_URL}${path}`, {
            method: 'POST',
            headers: authHeaders({ 'Content-Type': 'application/json' }),
            body: JSON.stringify(body),
        }).then(handleUnauthorized)
    },
}
