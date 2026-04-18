import { Navigate, useLocation } from 'react-router-dom'

export default function ProtectedRoute({ children }) {
    const location = useLocation()
    if (!localStorage.getItem('token')) {
        const redirect = encodeURIComponent(location.pathname + location.search)
        return <Navigate to={`/login?redirect=${redirect}`} replace />
    }
    return children
}
