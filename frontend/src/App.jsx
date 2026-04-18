import { lazy, Suspense, useEffect } from 'react'
import { Routes, Route, useNavigate } from 'react-router-dom'
import ProtectedRoute from './components/frontend/ProtectedRoute/ProtectedRoute'

const Home = lazy(() => import('./pages/Home/Home'))
const About = lazy(() => import('./pages/About/About'))
const Login = lazy(() => import('./pages/Login/Login'))

function TokenHandler() {
    const navigate = useNavigate()

    useEffect(() => {
        const params = new URLSearchParams(window.location.search)
        const token = params.get('token')
        if (token) {
            localStorage.setItem('token', token)
            navigate('/', { replace: true })
        }
    }, [navigate])

    return null
}

function App() {
    return (
        <Suspense fallback={null}>
            <TokenHandler />
            <Routes>
                <Route path="/login" element={<Login />} />
                <Route path="/" element={<ProtectedRoute><Home /></ProtectedRoute>} />
                <Route path="/watch/:id" element={<ProtectedRoute><Home /></ProtectedRoute>} />
                <Route path="/about" element={<About />} />
            </Routes>
        </Suspense>
    )
}

export default App
