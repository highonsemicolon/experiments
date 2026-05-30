import { useAuth0 } from '@auth0/auth0-react'

const AuthButtons = () => {
    const { loginWithRedirect, logout, isAuthenticated } = useAuth0()

    return (
        <div>
            {isAuthenticated ? (
                <button onClick={() => logout({ returnTo: window.location.origin })}>
                    Logout
                </button>
            ) : (
                <button onClick={() => loginWithRedirect()}>Login</button>
            )}
        </div>
    )
}

export default AuthButtons
