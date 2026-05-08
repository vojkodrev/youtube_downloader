import { Navigate, useLocation } from 'react-router-dom'

function isPublicHost() {
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

export default function ProtectedRoute({ children }) {
    const location = useLocation()
    if (!isPublicHost() && !localStorage.getItem('token')) {
        const redirect = encodeURIComponent(location.pathname + location.search)
        return <Navigate to={`/login?redirect=${redirect}`} replace />
    }
    return children
}
