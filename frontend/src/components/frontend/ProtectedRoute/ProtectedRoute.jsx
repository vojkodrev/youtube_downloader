import { Navigate, useLocation } from 'react-router-dom'
import { isPublicHost } from '../../../api'

export default function ProtectedRoute({ children }) {
    const location = useLocation()
    if (!isPublicHost() && !localStorage.getItem('token')) {
        const redirect = encodeURIComponent(location.pathname + location.search)
        return <Navigate to={`/login?redirect=${redirect}`} replace />
    }
    return children
}
