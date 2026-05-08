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
