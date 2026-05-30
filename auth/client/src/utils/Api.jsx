import { useAuth0 } from '@auth0/auth0-react'

export const useApi = () => {
  const { getAccessTokenSilently } = useAuth0()

  const callApi = async (endpoint, options = {}) => {

    const token = await getAccessTokenSilently()
    console.log(token)
    const response = await fetch(endpoint, {
      ...options,
      headers: {
        ...options.headers,
        Authorization: `Bearer ${token}`,
      },
    })

    if (!response.ok) {
      throw new Error("Network response was not ok")
    }
    return response.json()
  }

  return { callApi }
}
