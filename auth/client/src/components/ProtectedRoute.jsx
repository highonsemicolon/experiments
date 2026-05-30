import { useAuth0 } from '@auth0/auth0-react'
import { Navigate } from 'react-router-dom'
import PropTypes from 'prop-types'

const ProtectedRoute = ({ component: Component }) => {
    const { isAuthenticated } = useAuth0()

    return isAuthenticated ? <Component /> : <Navigate to="/" />
}

ProtectedRoute.propTypes = {
    component: PropTypes.elementType.isRequired,
}

export default ProtectedRoute
