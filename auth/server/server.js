const express = require('express')
const { auth } = require('express-oauth2-jwt-bearer')
const cors = require('cors')

const port = process.env.PORT || 8081

const jwtCheck = auth({
  audience: 'pathfinder-backend',
  issuerBaseURL: 'https://pathfinder-1.us.auth0.com/',
  tokenSigningAlg: 'RS256',
})

const app = express()

app.use(cors({
  origin: 'http://localhost:8080',
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

app.listen(port)

console.log('Running on port ', port);

