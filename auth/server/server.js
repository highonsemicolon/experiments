require('dotenv').config()
const express = require('express')
const { auth } = require('express-oauth2-jwt-bearer')
const cors = require('cors')

const {
  PORT = 8081,
  CORS_ORIGIN = 'http://localhost:8080',
  AUTH_AUDIENCE,
  AUTH_ISSUER_BASE_URL,
  AUTH_TOKEN_SIGNING_ALG = 'RS256',
} = process.env

const jwtCheck = auth({
  audience: AUTH_AUDIENCE,
  issuerBaseURL: AUTH_ISSUER_BASE_URL,
  tokenSigningAlg: AUTH_TOKEN_SIGNING_ALG,
})

const app = express()

app.use(cors({
  origin: CORS_ORIGIN,
  credentials: true,
  methods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS'],
  allowedHeaders: ['Content-Type', 'Authorization']
}))

app.get('/', function (req, res) {
  res.send('welcome')
})

app.get('/protected-endpoint', jwtCheck, (req, res) => {
  res.json([
    { id: 1, text: 'Hello' },
    { id: 2, text: 'Welcome' }
  ])
})

app.listen(PORT)

console.log(`Running on port ${PORT}`);

