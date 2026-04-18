const API_URL = import.meta.env.VITE_API_URL

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
